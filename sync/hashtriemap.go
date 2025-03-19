package sync

import (
	isync "github.com/metacubex/backport-std/internal/sync"
)

type HashTrieMap[K comparable, V any] struct {
	_ noCopy
	isync.HashTrieMap[K, V]
}
