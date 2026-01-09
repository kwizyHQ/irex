package tempdir

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// TempDir manages a single temporary directory for the program lifecycle.
type TempDir struct {
	path string
}

var (
	once      sync.Once
	singleton *TempDir
)

// Get returns the singleton TempDir instance, creating it if necessary.
func Get() *TempDir {
	once.Do(func() {
		dir, err := os.MkdirTemp("", "irex_temp_*")
		if err != nil {
			panic(err)
		}
		singleton = &TempDir{path: dir}
	})
	return singleton
}

// Path returns the path to the temp directory.
func (t *TempDir) Path() string {
	return t.path
}

// Delete removes the temp directory and all its contents.
func (t *TempDir) Delete() error {
	return os.RemoveAll(t.path)
}

// Read a file from the temp directory.
func (t *TempDir) ReadFile(relPath string) (string, error) {
	data, err := os.ReadFile(filepath.Join(t.path, relPath))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Clean specfied directory inside the temp directory. Deletes all contents.
func (t *TempDir) Clean(relPath string) error {
	fullPath := filepath.Join(t.path, relPath)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		err := os.RemoveAll(filepath.Join(fullPath, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteFile removes a file in the temp directory.
func (t *TempDir) DeleteFile(relPath string) error {
	return os.Remove(filepath.Join(t.path, relPath))
}

// OverwriteFile writes data to a file in the temp directory, replacing if exists.
func (t *TempDir) OverwriteFile(relPath string, data []byte) error {
	return os.WriteFile(filepath.Join(t.path, relPath), data, 0644)
}

// CopyFile copies a file from src (can be fs.FS or os) to the temp directory.
func (t *TempDir) CopyFile(srcFS fs.FS, srcPath, destRelPath string) error {
	f, err := srcFS.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()
	out, err := os.Create(filepath.Join(t.path, destRelPath))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, f)
	return err
}

// CopyFolder copies all files from srcFS/srcDir to the temp directory under destRelPath.
func (t *TempDir) CopyFolder(srcFS fs.FS, srcDir, destRelPath string) error {
	return fs.WalkDir(srcFS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(srcDir, path)
		dest := filepath.Join(t.path, destRelPath, rel)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		return t.CopyFile(srcFS, path, filepath.Join(destRelPath, rel))
	})
}
