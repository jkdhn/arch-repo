package api

import (
	"arch-repo/repository"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type API struct {
	repo    *repository.Repository
	bucket  string
	presign *s3.PresignClient
}

func New(bucket string, repo *repository.Repository, presign *s3.PresignClient) *API {
	return &API{
		repo:    repo,
		bucket:  bucket,
		presign: presign,
	}
}
