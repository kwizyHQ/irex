package utils

import (
	"io"
	"os"
	"path/filepath"
)

// WriteToFile writes data to the specified file path, creating directories if needed.
func WriteToFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// ReadFromFile reads the contents of the specified file and returns it as a byte slice.
func ReadFromFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// FileExists checks if the given file exists and is not a directory.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// EnsureDir ensures that a directory exists at the given path, creating it if necessary.
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
