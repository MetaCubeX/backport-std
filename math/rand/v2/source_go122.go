//go:build go1.22

package rand

import "math/rand/v2"

func (runtimeSource) Uint64() uint64 {
	return rand.Uint64()
}
