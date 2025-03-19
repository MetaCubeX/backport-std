//go:build !go1.24

package maphash

import (
	"github.com/metacubex/backport-std/internal/abi"
	"unsafe"
)

// ptrSize is the size of a pointer in bytes - unsafe.Sizeof(uintptr(0)) but as an ideal constant.
// It is also the size of the machine's native word size (that is, 4 on 32-bit systems, 8 on 64-bit).
const ptrSize = 4 << (^uintptr(0) >> 63)

func Comparable[T comparable](s Seed, v T) uint64 {
	return comparableHash(*(*seedTyp)(unsafe.Pointer(&s)), v)
}

func comparableHash[T comparable](seed seedTyp, v T) uint64 {
	s := seed.s
	var m map[T]struct{}
	mTyp := abi.TypeOf(m)
	var hasher func(unsafe.Pointer, uintptr) uintptr
	if abi.SwissMap {
		hasher = (*abi.SwissMapType)(unsafe.Pointer(mTyp)).Hasher
	} else {
		hasher = (*abi.OldMapType)(unsafe.Pointer(mTyp)).Hasher
	}

	p := abi.Escape(unsafe.Pointer(&v))

	if ptrSize == 8 {
		return uint64(hasher(p, uintptr(s)))
	}
	lo := hasher(p, uintptr(s))
	hi := hasher(p, uintptr(s>>32))
	return uint64(hi)<<32 | uint64(lo)
}

// WriteComparable adds x to the data hashed by h.
func WriteComparable[T comparable](h *Hash, x T) {
	// writeComparable (not in purego mode) directly operates on h.state
	// without using h.buf. Mix in the buffer length so it won't
	// commute with a buffered write, which either changes h.n or changes
	// h.state.
	hash := (*hashTyp)(unsafe.Pointer(h))
	if hash.n != 0 {
		hash.state.s = comparableHash(hash.state, hash.n)
	}
	hash.state.s = comparableHash(hash.state, x)
}

// go/src/hash/maphash/maphash.go
type hashTyp struct {
	_     [0]func() // not comparable
	seed  seedTyp   // initial seed used for this hash
	state seedTyp   // current hash of all flushed bytes
	buf   [128]byte // unflushed byte buffer
	n     int       // number of unflushed bytes
}

type seedTyp struct {
	s uint64
}
