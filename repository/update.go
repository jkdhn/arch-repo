package repository

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (r *Repository) Update() error {
	r.updateMutex.Lock()
	defer r.updateMutex.Unlock()

	data := new(bytes.Buffer)
	if err := NewEncoder(data).Encode(r); err != nil {
		return err
	}

	md5Sum := md5.Sum(data.Bytes())
	md5Hash := base64.StdEncoding.EncodeToString(md5Sum[:])

	reader := bytes.NewReader(data.Bytes())

	if _, err := r.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        &r.bucket,
		Key:           aws.String(fmt.Sprintf("%v.db", r.name)),
		ContentMD5:    &md5Hash,
		Body:          reader,
	}); err != nil {
		return err
	}

	return nil
}
