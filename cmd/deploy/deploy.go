package main

import (
	"arch-repo/api"
	"arch-repo/pkg"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var apiUrl string
var authToken string

func init() {
	deployCmd.Flags().StringVarP(&apiUrl, "api-url", "a", "", "repository API URL")
	deployCmd.Flags().StringVarP(&authToken, "auth-token", "t", "", "JWT token")
	_ = deployCmd.MarkFlagRequired("api-url")
}

var deployCmd = &cobra.Command{
	Use:  "deploy --api-url url filename",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}

		p, err := pkg.NewDecoder(file, filepath.Base(file.Name())).Decode()
		if err != nil {
			return err
		}

		payload := new(bytes.Buffer)

		request := &api.UploadRequest{
			Description: *p.Description(),
		}

		if err := json.NewEncoder(payload).Encode(request); err != nil {
			return err
		}

		fmt.Println("Adding to database")

		databaseRequest, err := http.NewRequest(http.MethodPost, apiUrl, payload)
		if err != nil {
			return err
		}

		databaseRequest.Header.Set("Content-Type", "application/json")
		if authToken != "" {
			databaseRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %v", authToken))
		}

		databaseResponse, err := http.DefaultClient.Do(databaseRequest)
		if err != nil {
			return err
		}

		if databaseResponse.StatusCode != http.StatusOK {
			data, _ := io.ReadAll(databaseResponse.Body)
			_ = databaseResponse.Body.Close()
			return fmt.Errorf("%v: %v", databaseResponse.Status, string(data))
		}

		var response api.UploadResponse
		if err := json.NewDecoder(databaseResponse.Body).Decode(&response); err != nil {
			_ = databaseResponse.Body.Close()
			return err
		}

		if err := databaseResponse.Body.Close(); err != nil {
			return err
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return err
		}

		fmt.Println("Uploading to repository")

		uploadRequest, err := http.NewRequest(response.UploadMethod, response.UploadURL, file)
		if err != nil {
			return err
		}

		uploadRequest.ContentLength = int64(p.Description().CompressedSize)
		uploadRequest.Header = response.UploadHeader

		uploadResponse, err := http.DefaultClient.Do(uploadRequest)
		if err != nil {
			return err
		}

		if uploadResponse.StatusCode != http.StatusOK {
			data, _ := io.ReadAll(uploadResponse.Body)
			_ = uploadResponse.Body.Close()
			return fmt.Errorf("%v: %v", uploadResponse.Status, string(data))
		}

		if err := uploadResponse.Body.Close(); err != nil {
			return err
		}

		return nil
	},
}
