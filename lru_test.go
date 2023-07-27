package cachego

import (
	"fmt"
	"sync"
	"testing"
)

// nolint:errcheck
func TestLRUCache(t *testing.T) {
	// Create an LRU cache with size 3
	cache := NewLRUCache[int, string](3)
	arr := []string{"zero", "one", "two"}
	// Test Set and Get operations
	for key, val := range arr {
		cache.Set(key, val)
	}

	// Retrieve elements from the cache
	for key, val := range arr {
		value, err := cache.Get(key)
		if err != nil {
			t.Errorf("Expected key %v to be found in cache, but it was not found", key)
		}
		if value != val {
			t.Errorf("Expected value '%v' for key %v, but got '%v'", val, key, value)
		}
	}

	// Test eviction of the least recently used element
	cache.Set(3, "three")
	_, err := cache.Get(0)
	if err == nil {
		t.Errorf("Expected key %v to be deleted from the cache, but it was found", 0)
	}

	// Test setting a value for an entry that is already in cache
	cache.Set(1, "new")
	value, err := cache.Get(1)
	if err != nil {
		t.Errorf("Expected key %v to be found in cache, but it was not found", 1)
	}
	if value != "new" {
		t.Errorf("Expected value 'new' for key %v, but got '%v'", 1, value)
	}

	// Test Get for a non-existent key
	_, err = cache.Get(4)
	expectedError := fmt.Sprintf("key %v not found", 4)
	if err == nil {
		t.Errorf("Expected error: %v, but got nil", expectedError)
	}
	if err.Error() != expectedError {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err.Error())
	}

	// Test Delete operation
	cache.Delete(3)
	_, err = cache.Get(3)
	if err == nil {
		t.Errorf("Expected key %v to be deleted from the cache, but it was found", 3)
	}

	// Test Delete for a non-existent key
	err = cache.Delete(3)
	if err == nil {
		t.Errorf("Expected error: %v, but got nil", expectedError)
	}

	// Test Clear operation
	cache.Clear()
	_, err = cache.Get(1)
	if err == nil {
		t.Errorf("Expected cache to be cleared, but key %v was found in the cache", 1)
	}

	// Test setting new values after clearing
	cache.Set(5, "five")
	value, err = cache.Get(5)
	if err != nil {
		t.Errorf("Expected key %v to be found in cache, but it was not found", 5)
	}
	if value != "five" {
		t.Errorf("Expected value 'five' for key %v, but got '%v'", 5, value)
	}
}

// nolint:errcheck
func TestLRUCacheConcurrency(t *testing.T) {
	// Create an LRU cache with size 3
	cache := NewLRUCache[int, string](3)

	// Number of concurrent operations
	numOps := 100

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numOps * 2)

	// Concurrently set values
	for i := 1; i <= numOps; i++ {
		key := i
		value := fmt.Sprintf("value%d", i)

		go func() {
			cache.Set(key, value)
			wg.Done()
		}()
	}

	// Concurrently get values
	for i := 1; i <= numOps; i++ {
		key := i

		go func() {
			cache.Get(key)
			wg.Done()
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	var count int
	for i := 1; i <= numOps; i++ {
		_, err := cache.Get(i)
		if err == nil {
			count++
		}
	}

	if count != 3 {
		t.Errorf("Expected cache size to be %v, but got %v", 3, count)
	}
}
