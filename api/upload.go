package api

import (
	"arch-repo/pkg/desc"
	"encoding/base64"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type UploadRequest struct {
	Description desc.Description `json:"description"`
}

type UploadResponse struct {
	UploadURL    string              `json:"upload_url"`
	UploadMethod string              `json:"upload_method"`
	UploadHeader map[string][]string `json:"upload_header"`
}

var pattern = regexp.MustCompile("^[a-z0-9.\\-_:]+$")

func (a *API) Upload(c *gin.Context) {
	var request UploadRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}

	if !pattern.MatchString(request.Description.Name) {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid name"))
		return
	}

	if !pattern.MatchString(request.Description.Version) {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid version"))
		return
	}

	if !pattern.MatchString(request.Description.FileName) {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid filename"))
		return
	}

	if err := a.repo.Store(&request.Description); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	md5Hash := base64.StdEncoding.EncodeToString(request.Description.MD5Sum[:])

	operation, err := a.presign.PresignPutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket:     &a.bucket,
		Key:        &request.Description.FileName,
		ContentMD5: &md5Hash,
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid filename"))
		return
	}

	response := &UploadResponse{
		UploadURL:    operation.URL,
		UploadMethod: operation.Method,
		UploadHeader: operation.SignedHeader,
	}

	c.JSON(http.StatusOK, response)
}
