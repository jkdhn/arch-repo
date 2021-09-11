package repository

import (
	"arch-repo/pkg/desc"
	"arch-repo/storage"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"strings"
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

func (r *Repository) Cleanup() error {
	files := make(map[string]bool)

	input := &s3.ListObjectsV2Input{
		Bucket: &r.bucket,
	}
	for {
		output, err := r.client.ListObjectsV2(context.Background(), input)
		if err != nil {
			return err
		}

		for _, obj := range output.Contents {
			key := *obj.Key
			if !strings.HasSuffix(key, ".db") {
				files[*obj.Key] = true
			}
		}

		input.ContinuationToken = output.NextContinuationToken
		if !output.IsTruncated {
			break
		}
	}

	if err := r.storage.ForEach(func(description *desc.Description) error {
		delete(files, description.FileName)
		return nil
	}); err != nil {
		return err
	}

	for key := range files {
		fmt.Println("Deleting", key)
		_, err := r.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: &r.bucket,
			Key:    &key,
		})
		if err != nil {
			return err
		}
	}
	fmt.Println("Deleted", len(files), "unreferenced packages")

	return nil
}
