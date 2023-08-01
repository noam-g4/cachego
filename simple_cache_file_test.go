package cachego

import (
	"os"
	"testing"
)

func TestSimpleCacheFile(t *testing.T) {
	nums := []string{"one", "two"}
	filename := "cache.json"

	// setting up a cache with a cache file
	file := NewSimpleCacheFile(filename)
	cache := NewCache[int, string](Opts{Size: 2, File: file})

	for i, num := range nums {
		cache.Set(i+1, num)
	}

	err := cache.Clear()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	// setting up a new cache with the same cache file
	cache2 := NewCache[int, string](Opts{Size: 2, File: file})

	for i, num := range nums {
		val, err := cache2.Get(i + 1)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if val != num {
			t.Errorf("expected %v, got %v", num, val)
		}
	}

	// test cache load with a file that is too large (file should be discarded)
	cache3 := NewCache[int, string](Opts{Size: 1, File: file})
	if _, err := cache3.Get(1); err == nil {
		t.Errorf("expected error, got nil")
	}

	os.Remove(filename)

	// test loading an invaid file
	os.WriteFile(filename, []byte("invalid"), 0644)
	cache4 := NewCache[int, string](Opts{Size: 1, File: file})
	if _, err := cache4.Get(1); err == nil {
		t.Errorf("expected error, got nil")
	}

	os.Remove(filename)
}
