package main

import (
	"arch-repo/api"
	"arch-repo/auth"
	"arch-repo/repository"
	"arch-repo/storage/file"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var issuer string
var jwksUrl string
var skipAuth bool
var bucket string
var name string

func init() {
	serverCmd.Flags().StringVarP(&issuer, "issuer", "i", "", "JWT issuer")
	serverCmd.Flags().StringVarP(&jwksUrl, "jwks", "k", "", "JWT key store URL")
	serverCmd.Flags().BoolVar(&skipAuth, "skip-auth", false, "Disable authentication")
	serverCmd.Flags().StringVarP(&bucket, "bucket", "b", "", "AWS bucket")
	serverCmd.Flags().StringVarP(&name, "name", "n", "", "Repository name")
	_ = serverCmd.MarkFlagRequired("bucket")
	_ = serverCmd.MarkFlagRequired("name")
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			logrus.WithError(err).Fatal("Loading AWS config failed")
		}

		client := s3.NewFromConfig(cfg)
		presign := s3.NewPresignClient(client)

		s := file.NewStorage(".")
		repo := repository.NewRepository(s, client, bucket, name)
		a := api.New(bucket, repo, presign)

		r := gin.Default()

		if issuer != "" {
			var verifier *oidc.IDTokenVerifier
			if jwksUrl != "" {
				keySet := oidc.NewRemoteKeySet(context.Background(), jwksUrl)
				verifier = oidc.NewVerifier(issuer, keySet, &oidc.Config{
					SkipClientIDCheck: true,
				})
			} else {
				provider, err := oidc.NewProvider(context.Background(), issuer)
				if err != nil {
					return err
				}
				verifier = provider.Verifier(&oidc.Config{
					SkipClientIDCheck: true,
				})
			}
			r.Use(auth.JWTMiddleware(verifier))
		} else if !skipAuth {
			return errors.New("authentication not configured")
		}

		r.POST("/upload", a.Upload)

		return r.Run()
	},
}
