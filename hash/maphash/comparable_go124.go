//go:build go1.24

package maphash

import "hash/maphash"

func Comparable[T comparable](s Seed, v T) uint64 {
	return maphash.Comparable(s, v)
}

func WriteComparable[T comparable](h *Hash, x T) {
	maphash.WriteComparable(h, x)
}
