package memory

import (
	"arch-repo/pkg/desc"
	"arch-repo/storage"
)

type Storage struct {
	data []*desc.Description
}

var _ storage.Storage = (*Storage)(nil)

func NewStorage() *Storage {
	return new(Storage)
}

func (s *Storage) Put(description *desc.Description) error {
	s.data = append(s.data, description)
	return nil
}

func (s *Storage) ForEach(f func(*desc.Description) error) error {
	for _, d := range s.data {
		if err := f(d); err != nil {
			return err
		}
	}
	return nil
}
