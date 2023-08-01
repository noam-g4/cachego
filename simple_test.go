package cachego

import (
	"sync"
	"testing"
	"time"
)

func TestSimpleCache(t *testing.T) {
	c := NewCache[int, string](Opts{Size: 1})
	if c == nil {
		t.Error("NewSimple returned nil")
	}

	// Set
	if err := c.Set(1, "one"); err != nil {
		t.Errorf("Set returned error: %s", err)
	}

	// Set (full)
	if err := c.Set(2, "two"); err == nil {
		t.Errorf("Set returned nil error when cache is full")
	}

	// Get
	v, err := c.Get(1)
	if err != nil {
		t.Errorf("Get returned error: %s", err)
	}

	if v != "one" {
		t.Errorf("Get returned incorrect value: %s", v)
	}

	// Get (not found)
	if _, err := c.Get(2); err == nil {
		t.Errorf("Get returned nil error when key not found")
	}

	// Delete
	if err := c.Delete(1); err != nil {
		t.Errorf("Delete returned error: %s", err)
	}

	// delete (not found)
	if err := c.Delete(2); err == nil {
		t.Errorf("Delete returned nil error when key not found")
	}

	if _, err := c.Get(1); err == nil {
		t.Errorf("Get returned nil error after Delete")
	}

	// Clear
	if err := c.Set(1, "one"); err != nil {
		t.Errorf("Set returned error: %s", err)
	}

	c.Clear() // errcheck: ignore

	if _, err := c.Get(1); err == nil {
		t.Errorf("Get returned nil error after Clear")
	}
}

func TestSimpleCacheConcurrency(t *testing.T) {
	c := NewCache[int, string](Opts{Size: 100})
	wg := sync.WaitGroup{}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			if err := c.Set(i, "one"); err != nil {
				t.Errorf("Set returned error: %s", err)
			}
		}(i)
	}
	wg.Wait()

	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			if _, err := c.Get(i); err != nil {
				t.Errorf("Get returned error: %s", err)
			}
		}(i)
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			if err := c.Delete(i); err != nil {
				t.Errorf("Delete returned error: %s", err)
			}
		}(i)
	}
	wg.Wait()

	if _, err := c.Get(1); err == nil {
		t.Errorf("Get returned nil error after Clear")
	}
}

func TestCache(t *testing.T) {
	c := NewCache[int, string](Opts{Size: 1, TTL: 1})

	if err := c.Set(1, "one"); err != nil {
		t.Errorf("Set returned error: %s", err)
	}

	time.Sleep(2 * time.Second)

	if _, err := c.Get(1); err == nil {
		t.Errorf("Get returned nil error after TTL")
	}
}
