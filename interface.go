package cachego

const defaultSize = 100

// Cache is an interface that represents a generic key-value cache.
// Implementations of this interface are expected to provide mechanisms
// for storing and retrieving data in an efficient manner based on the
// specified key and value types.
type Cache[K comparable, V any] interface {
	// Set stores the provided value under the given key in the cache.
	// If the key already exists in the cache, its associated value will be updated.
	// An error may be returned if there's a problem while setting the key-value pair.
	Set(key K, value V) error

	// Get retrieves the value associated with the given key from the cache.
	// If the key is found in the cache, the corresponding value and nil error will be returned.
	// If the key is not found, the zero value of the value type and an error will be returned.
	Get(key K) (V, error)

	// Delete removes the key-value pair associated with the given key from the cache.
	// If the key is found in the cache, it will be deleted, and a nil error will be returned.
	// If the key is not found, an error will be returned.
	Delete(key K) error

	// Clear clears the entire cache, removing all key-value pairs.
	// After this operation, the cache will be empty, and a nil error will be returned.
	// Implementations may choose to release any resources associated with the cache
	// during this operation.
	Clear() error
}

// File represents an interface for loading from and dumping data to a file.
type File interface {
	// Load reads the contents of the file and returns the data read from the file as a byte slice.
	// If the operation is successful, it returns the read data and a nil error.
	// If an error occurs during the load operation, it returns a non-nil error.
	Load() ([]byte, error)

	// Dump writes the given data as a byte slice to the file.
	// If the operation is successful, it returns a nil error.
	// If an error occurs during the dump operation, it returns a non-nil error.
	Dump(data []byte) error
}
