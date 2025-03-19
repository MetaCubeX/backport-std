// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tls13 implements the TLS 1.3 Key Schedule as specified in RFC 8446,
// Section 7.1 and allowed by FIPS 140-3 IG 2.4.B Resolution 7.
package tls13

import (
	"github.com/metacubex/backport-std/crypto/internal/hkdf"
	"github.com/metacubex/backport-std/internal/byteorder"
	"hash"
)

// We don't set the service indicator in this package but we delegate that to
// the underlying functions because the TLS 1.3 KDF does not have a standard of
// its own.

// ExpandLabel implements HKDF-Expand-Label from RFC 8446, Section 7.1.
func ExpandLabel(hash func() hash.Hash, secret []byte, label string, context []byte, length int) []byte {
	if len("tls13 ")+len(label) > 255 || len(context) > 255 {
		// It should be impossible for this to panic: labels are fixed strings,
		// and context is either a fixed-length computed hash, or parsed from a
		// field which has the same length limitation.
		//
		// Another reasonable approach might be to return a randomized slice if
		// we encounter an error, which would break the connection, but avoid
		// panicking. This would perhaps be safer but significantly more
		// confusing to users.
		panic("tls13: label or context too long")
	}
	hkdfLabel := make([]byte, 0, 2+1+len("tls13 ")+len(label)+1+len(context))
	hkdfLabel = byteorder.BEAppendUint16(hkdfLabel, uint16(length))
	hkdfLabel = append(hkdfLabel, byte(len("tls13 ")+len(label)))
	hkdfLabel = append(hkdfLabel, "tls13 "...)
	hkdfLabel = append(hkdfLabel, label...)
	hkdfLabel = append(hkdfLabel, byte(len(context)))
	hkdfLabel = append(hkdfLabel, context...)
	return hkdf.Expand(hash, secret, string(hkdfLabel), length)
}

func extract(hash func() hash.Hash, newSecret, currentSecret []byte) []byte {
	if newSecret == nil {
		newSecret = make([]byte, hash().Size())
	}
	return hkdf.Extract(hash, newSecret, currentSecret)
}

func deriveSecret(hash func() hash.Hash, secret []byte, label string, transcript hash.Hash) []byte {
	if transcript == nil {
		transcript = hash()
	}
	return ExpandLabel(hash, secret, label, transcript.Sum(nil), transcript.Size())
}

const (
	resumptionBinderLabel         = "res binder"
	clientEarlyTrafficLabel       = "c e traffic"
	clientHandshakeTrafficLabel   = "c hs traffic"
	serverHandshakeTrafficLabel   = "s hs traffic"
	clientApplicationTrafficLabel = "c ap traffic"
	serverApplicationTrafficLabel = "s ap traffic"
	earlyExporterLabel            = "e exp master"
	exporterLabel                 = "exp master"
	resumptionLabel               = "res master"
)

type EarlySecret struct {
	secret []byte
	hash   func() hash.Hash
}

func NewEarlySecret(hashFn func() hash.Hash, psk []byte) *EarlySecret {
	return &EarlySecret{
		secret: extract(hashFn, psk, nil),
		hash:   func() hash.Hash { return hashFn() },
	}
}

func (s *EarlySecret) ResumptionBinderKey() []byte {
	return deriveSecret(s.hash, s.secret, resumptionBinderLabel, nil)
}

// ClientEarlyTrafficSecret derives the client_early_traffic_secret from the
// early secret and the transcript up to the ClientHello.
func (s *EarlySecret) ClientEarlyTrafficSecret(transcript hash.Hash) []byte {
	return deriveSecret(s.hash, s.secret, clientEarlyTrafficLabel, transcript)
}

type HandshakeSecret struct {
	secret []byte
	hash   func() hash.Hash
}

func (s *EarlySecret) HandshakeSecret(sharedSecret []byte) *HandshakeSecret {
	derived := deriveSecret(s.hash, s.secret, "derived", nil)
	return &HandshakeSecret{
		secret: extract(s.hash, sharedSecret, derived),
		hash:   s.hash,
	}
}

// ClientHandshakeTrafficSecret derives the client_handshake_traffic_secret from
// the handshake secret and the transcript up to the ServerHello.
func (s *HandshakeSecret) ClientHandshakeTrafficSecret(transcript hash.Hash) []byte {
	return deriveSecret(s.hash, s.secret, clientHandshakeTrafficLabel, transcript)
}

// ServerHandshakeTrafficSecret derives the server_handshake_traffic_secret from
// the handshake secret and the transcript up to the ServerHello.
func (s *HandshakeSecret) ServerHandshakeTrafficSecret(transcript hash.Hash) []byte {
	return deriveSecret(s.hash, s.secret, serverHandshakeTrafficLabel, transcript)
}

type MasterSecret struct {
	secret []byte
	hash   func() hash.Hash
}

func (s *HandshakeSecret) MasterSecret() *MasterSecret {
	derived := deriveSecret(s.hash, s.secret, "derived", nil)
	return &MasterSecret{
		secret: extract(s.hash, nil, derived),
		hash:   s.hash,
	}
}

// ClientApplicationTrafficSecret derives the client_application_traffic_secret_0
// from the master secret and the transcript up to the server Finished.
func (s *MasterSecret) ClientApplicationTrafficSecret(transcript hash.Hash) []byte {
	return deriveSecret(s.hash, s.secret, clientApplicationTrafficLabel, transcript)
}

// ServerApplicationTrafficSecret derives the server_application_traffic_secret_0
// from the master secret and the transcript up to the server Finished.
func (s *MasterSecret) ServerApplicationTrafficSecret(transcript hash.Hash) []byte {
	return deriveSecret(s.hash, s.secret, serverApplicationTrafficLabel, transcript)
}

// ResumptionMasterSecret derives the resumption_master_secret from the master secret
// and the transcript up to the client Finished.
func (s *MasterSecret) ResumptionMasterSecret(transcript hash.Hash) []byte {
	return deriveSecret(s.hash, s.secret, resumptionLabel, transcript)
}

type ExporterMasterSecret struct {
	secret []byte
	hash   func() hash.Hash
}

// ExporterMasterSecret derives the exporter_master_secret from the master secret
// and the transcript up to the server Finished.
func (s *MasterSecret) ExporterMasterSecret(transcript hash.Hash) *ExporterMasterSecret {
	return &ExporterMasterSecret{
		secret: deriveSecret(s.hash, s.secret, exporterLabel, transcript),
		hash:   s.hash,
	}
}

// EarlyExporterMasterSecret derives the exporter_master_secret from the early secret
// and the transcript up to the ClientHello.
func (s *EarlySecret) EarlyExporterMasterSecret(transcript hash.Hash) *ExporterMasterSecret {
	return &ExporterMasterSecret{
		secret: deriveSecret(s.hash, s.secret, earlyExporterLabel, transcript),
		hash:   s.hash,
	}
}

func (s *ExporterMasterSecret) Exporter(label string, context []byte, length int) []byte {
	secret := deriveSecret(s.hash, s.secret, label, nil)
	h := s.hash()
	h.Write(context)
	return ExpandLabel(s.hash, secret, "exporter", h.Sum(nil), length)
}

func TestingOnlyExporterSecret(s *ExporterMasterSecret) []byte {
	return s.secret
}
