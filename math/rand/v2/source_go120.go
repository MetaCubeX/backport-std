//go:build go1.20 && !go1.22

package v2

import (
	_ "unsafe" // for go:linkname
)

//go:linkname runtimefastrand64 runtime.fastrand64
func runtimefastrand64() uint64

func (runtimeSource) Uint64() uint64 {
	return runtimefastrand64()
}
