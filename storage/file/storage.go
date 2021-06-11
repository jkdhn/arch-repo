package file

import (
	"arch-repo/pkg/desc"
	"arch-repo/storage"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type Storage struct {
	mutex sync.Mutex
	path  string
}

var _ storage.Storage = (*Storage)(nil)

func NewStorage(path string) *Storage {
	_ = os.MkdirAll(path, 0755)
	return &Storage{
		path: path,
	}
}

func (s *Storage) Put(description *desc.Description) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	name := filepath.Join(s.path, fmt.Sprintf("%v.json", description.Name))
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	if err := json.NewEncoder(file).Encode(description); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) ForEach(f func(*desc.Description) error) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return filepath.WalkDir(s.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer func() {
			_ = file.Close()
		}()

		var description desc.Description
		if err := json.NewDecoder(file).Decode(&description); err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}

		return f(&description)
	})
}
