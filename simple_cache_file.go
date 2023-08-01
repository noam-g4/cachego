package cachego

import "os"

// simpleCacheFile is an implementation of the File interface.
// It represents a simple cache file that can be used for loading and dumping data.
type simpleCacheFile struct {
	path string
}

// NewSimpleCacheFile creates a new instance of the File interface backed by a simple cache file.
// The cache file is associated with the specified path.
// It returns a pointer to the File interface.
func NewSimpleCacheFile(path string) File {
	return &simpleCacheFile{path: path}
}

// Load reads the contents of the cache file and returns the data read from the file as a byte slice.
// If the operation is successful, it returns the read data and a nil error.
// If an error occurs during the load operation, it returns a non-nil error.
func (s *simpleCacheFile) Load() ([]byte, error) {
	return os.ReadFile(s.path)
}

// Dump writes the given data as a byte slice to the cache file.
// If the operation is successful, it returns a nil error.
// If an error occurs during the dump operation, it returns a non-nil error.
func (s *simpleCacheFile) Dump(data []byte) error {
	return os.WriteFile(s.path, data, 0644)
}
