//go:build go1.21

package buildin

import (
	"cmp"
)

func Min[T cmp.Ordered](x, y T) T {
	return min(x, y)
}

func Max[T cmp.Ordered](x, y T) T {
	return max(x, y)
}

func ClearArray[T ~[]E, E any](t T) {
	clear(t)
}

func ClearMap[T ~map[K]V, K comparable, V any](t T) {
	clear(t)
}
