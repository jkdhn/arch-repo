package storage

import "arch-repo/pkg/desc"

type Storage interface {
	Put(description *desc.Description) error
	ForEach(func(*desc.Description) error) error
}
