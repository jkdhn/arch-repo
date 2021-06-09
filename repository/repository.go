package repository

import (
	"arch-repo/pkg/desc"
	"arch-repo/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"sync"
)

type Repository struct {
	updateMutex sync.Mutex
	storage     storage.Storage
	client      *s3.Client
	bucket      string
	name        string
}

func NewRepository(s storage.Storage, client *s3.Client, bucket string, name string) *Repository {
	return &Repository{
		storage: s,
		client:  client,
		bucket:  bucket,
		name:    name,
	}
}

func (r *Repository) Store(description *desc.Description) error {
	if err := r.storage.Put(description); err != nil {
		return err
	}

	if err := r.Update(); err != nil {
		return err
	}

	return nil
}
