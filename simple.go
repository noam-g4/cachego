package cachego

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type simple[K comparable, V any] struct {
	size int32
	used int32
	ttl  int16 // in seconds
	data map[K]V
	mx   *sync.Mutex
	file File
}

type Opts struct {
	Size int32
	TTL  int16
	File File
}

// NewCache creates a new thread-safe instance of a cache with the specified size and ttl.
// If the size is less than or equal to zero, a default size of 100 will be used.
// If the ttl is less than or equal to zero, the cache will not expire.
func NewCache[K comparable, V any](opts Opts) Cache[K, V] {
	s := int32(defaultSize)
	if opts.Size > 0 {
		s = opts.Size
	}

	var used int32
	data := make(map[K]V, s)

	if opts.File != nil {

		if bytes, err := opts.File.Load(); err == nil {

			if err = json.Unmarshal(bytes, &data); err != nil {
				log.Printf("error unmarshalling cache data: %v", err)
			} else {

				l := int32(len(data))
				if l > s {
					log.Printf("cache data size %v is larger than cache size %v", l, s)
					data = make(map[K]V, s)
				} else {
					used = l
				}

			}

		} else {
			log.Printf("loading cache data failed: %v", err)
		}

	}

	return &simple[K, V]{
		size: s,
		used: used,
		data: data,
		mx:   &sync.Mutex{},
		ttl:  opts.TTL,
		file: opts.File,
	}
}

// Set stores the provided value under the given key in the cache.
// If the cache is full (reached its capacity), it returns an error "cache is full".
// If the key already exists in the cache, the associated value will be updated.
// This method is thread-safe.
func (c *simple[K, V]) Set(key K, value V) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if c.used >= c.size {
		return fmt.Errorf("cache is full")
	}

	if _, ok := c.data[key]; !ok {
		c.used++
	}

	c.data[key] = value

	if c.ttl > 0 {
		ctx, _ := c.setDeadline(key)
		go c.destroy(ctx, key)
	}

	return nil
}

// Get retrieves the value associated with the given key from the cache.
// If the key is found in the cache, the corresponding value and nil error will be returned.
// If the key is not found, the zero value of the value type and an error will be returned.
// This method is thread-safe.
func (c *simple[K, V]) Get(key K) (V, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if v, ok := c.data[key]; ok {
		return v, nil
	}

	var empty V
	return empty, fmt.Errorf("key %v not found", key)
}

// Delete removes the key-value pair associated with the given key from the cache.
// If the key is found in the cache, it will be deleted, and a nil error will be returned.
// If the key is not found, an error will be returned.
// This method is thread-safe.
func (c *simple[K, V]) Delete(key K) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.data[key]; !ok {
		return fmt.Errorf("key %v not found", key)
	}

	delete(c.data, key)
	c.used--
	return nil
}

// Clear clears the entire cache, removing all key-value pairs.
// After this operation, the cache will be empty, and a nil error will be returned.
// This method is thread-safe.
func (c *simple[K, V]) Clear() error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if c.file != nil {
		bytes, _ := json.Marshal(c.data)
		if err := c.file.Dump(bytes); err != nil {
			return err
		}
	}

	c.data = make(map[K]V, c.size)
	c.used = 0
	return nil
}

func (c *simple[K, V]) setDeadline(key K) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(time.Duration(c.ttl)*time.Second))
}

func (c *simple[K, V]) destroy(ctx context.Context, key K) {
	<-ctx.Done()
	c.Delete(key)
}
