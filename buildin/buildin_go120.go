//go:build !go1.21

package buildin

import (
	"github.com/metacubex/backport-std/cmp"
)

// isNaN reports whether x is a NaN without requiring the math package.
// This will always return false if T is not floating-point.
func isNaN[T cmp.Ordered](x T) bool {
	return x != x
}

func Min[T cmp.Ordered](x, y T) T {
	if isNaN(x) {
		return x
	}
	if isNaN(y) {
		return y
	}
	if x < y {
		return x
	}
	return y
}

func Max[T cmp.Ordered](x, y T) T {
	if isNaN(x) {
		return x
	}
	if isNaN(y) {
		return y
	}
	if x < y {
		return y
	}
	return x
}

func ClearArray[T ~[]E, E any](t T) {
	var defaultValue E
	for i := range t {
		t[i] = defaultValue
	}
}

func ClearMap[T ~map[K]V, K comparable, V any](t T) {
	for k := range t {
		delete(t, k)
	}
}
