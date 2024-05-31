package core

import (
	"encoding/json"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
	"text/template"
)

type Storage struct {
	path string
}

func (s *Storage) OpenFile(path string) (*os.File, error) {
	return openFile(filepath.Join(s.path, path))
}

func (s *Storage) RemoveFile(path string) {
	_ = os.Remove(filepath.Join(s.path, path))
}

func (s *Storage) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.path, path))
}

func (s *Storage) WriteFile(path string, data []byte) error {
	return os.WriteFile(filepath.Join(s.path, path), data, 0777)
}

func (s *Storage) GetTemplate(name string) (*template.Template, error) {
	data, err := s.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return template.New(name).Funcs(templateFuncs).Parse(string(data))
}

func (s *Storage) ReadTOML(path string, value any) error {
	data, err := s.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, value)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ReadJson(path string, value any) error {
	data, err := s.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, value)
}

func (s *Storage) WriteJson(path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.WriteFile(path, data)
}

func newStorage(path string) *Storage {
	storage := new(Storage)
	storage.path = path

	return storage
}
