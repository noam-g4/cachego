package cachego

import (
	"fmt"
	"sync"
)

type lru[K comparable, V any] struct {
	size  int32
	used  int32
	head  *node[K, V]
	tail  *node[K, V]
	cache map[K]*node[K, V]
	mx    *sync.Mutex
}

type node[K comparable, T any] struct {
	value T
	key   K
	next  *node[K, T]
	prev  *node[K, T]
}

// NewLRUCache creates a new thread-safe instance of an LRU cache with the given size.
// It returns a Cache[K, V] interface that can be used to interact with the cache.
func NewLRUCache[K comparable, V any](size int32) Cache[K, V] {
	return &lru[K, V]{
		size:  size,
		cache: make(map[K]*node[K, V]),
		mx:    &sync.Mutex{},
	}
}

// Set adds or updates a key-value pair in the LRU cache.
// If the key already exists in the cache, it updates its value and moves the item to the front of the cache (MRU position).
// If the key is new and the cache is already at its maximum size, it removes the least recently used item from the cache before adding the new item.
// Thread-safe.
func (l *lru[K, V]) Set(key K, value V) error {
	l.mx.Lock()
	defer l.mx.Unlock()

	if n, ok := l.cache[key]; ok {
		n.value = value
		l.unshift(n)
		return nil
	}

	n := &node[K, V]{key: key, value: value}
	l.unshift(n)
	l.cache[key] = n
	l.used++

	if l.used > l.size {
		l.pop()
		l.used--
	}

	return nil
}

// Get retrieves the value associated with the given key from the LRU cache.
// If the key is found in the cache, it moves the corresponding item to the front (MRU position) and returns its value.
// If the key is not found in the cache, it returns an error indicating that the key was not found.
// Thread-safe.
func (l *lru[K, V]) Get(key K) (V, error) {
	l.mx.Lock()
	defer l.mx.Unlock()

	var empty V
	err := fmt.Errorf("key %v not found", key)

	if n, ok := l.cache[key]; ok {
		l.pull(n)
		l.unshift(n)
		return n.value, nil
	}

	return empty, err
}

// Delete removes the key-value pair associated with the given key from the LRU cache.
// If the key is found in the cache, it removes the corresponding item from the cache and updates the cache size accordingly.
// If the key is not found in the cache, it returns an error indicating that the key was not found.
// Thread-safe.
func (l *lru[K, V]) Delete(key K) error {
	l.mx.Lock()
	defer l.mx.Unlock()

	if n, ok := l.cache[key]; ok {
		l.pull(n)
		delete(l.cache, key)
		l.used--
		return nil
	}

	return fmt.Errorf("key %v not found", key)
}

// Clear removes all items from the LRU cache, making it empty.
// Thread-safe.
func (l *lru[K, V]) Clear() error {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.head = nil
	l.tail = nil
	l.cache = make(map[K]*node[K, V])
	l.used = 0
	return nil
}

func (l *lru[K, V]) unshift(n *node[K, V]) {
	if l.head == nil {
		l.head = n
		l.tail = n
		return
	}

	n.next = l.head
	l.head.prev = n
	l.head = n
}

func (l *lru[K, V]) pop() {
	delete(l.cache, l.tail.key)
	l.tail.prev.next = nil
	l.tail = l.tail.prev
}

func (l *lru[K, V]) pull(n *node[K, V]) {
	if n.prev != nil {
		n.prev.next = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	}
}
