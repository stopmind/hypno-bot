package core

import (
	"os"
	"path/filepath"
)

type Storage struct {
	path string
}

func (s *Storage) OpenFile(path string) (*os.File, error) {
	return openWriteFile(filepath.Join(s.path, path))
}

func (s *Storage) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.path, path))
}

func newStorage(path string) *Storage {
	storage := new(Storage)
	storage.path = path

	return storage
}
