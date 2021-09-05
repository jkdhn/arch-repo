package main

import (
	"arch-repo/api"
	"arch-repo/auth"
	"arch-repo/repository"
	"arch-repo/storage/file"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var issuer string
var jwksUrl string
var claimsJson string
var skipAuth bool
var bucket string
var name string
var endpoint string
var endpointRegion string

func init() {
	serverCmd.Flags().StringVarP(&issuer, "issuer", "i", "", "JWT issuer")
	serverCmd.Flags().StringVarP(&jwksUrl, "jwks", "k", "", "JWT key store URL")
	serverCmd.Flags().StringVarP(&claimsJson, "claims", "c", "", "required JWT token claims (JSON)")
	serverCmd.Flags().BoolVar(&skipAuth, "skip-auth", false, "disable authentication")
	serverCmd.Flags().StringVarP(&bucket, "bucket", "b", "", "AWS bucket")
	serverCmd.Flags().StringVarP(&name, "name", "n", "", "repository name")
	serverCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "S3 endpoint")
	serverCmd.Flags().StringVarP(&endpointRegion, "endpoint-region", "r", "", "S3 endpoint signing region")
	_ = serverCmd.MarkFlagRequired("bucket")
	_ = serverCmd.MarkFlagRequired("name")
}

var serverCmd = &cobra.Command{
	Use:  "server",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadDefaultConfig(
			context.Background(),
			func(options *config.LoadOptions) error {
				if endpoint != "" {
					options.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
						return aws.Endpoint{
							URL:           endpoint,
							SigningRegion: endpointRegion,
						}, nil
					})
				}
				return nil
			},
		)
		if err != nil {
			logrus.WithError(err).Fatal("Loading AWS config failed")
		}

		client := s3.NewFromConfig(cfg)
		presign := s3.NewPresignClient(client)

		s := file.NewStorage(".")
		repo := repository.NewRepository(s, client, bucket, name)
		a := api.New(bucket, repo, presign)

		r := gin.Default()

		if !skipAuth && issuer != "" {
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

			var claims map[string]interface{}
			if claimsJson != "" {
				if err := json.Unmarshal([]byte(claimsJson), &claims); err != nil {
					return err
				}
			}

			r.Use(auth.JWTMiddleware(verifier, claims))
		} else if !skipAuth {
			return errors.New("authentication not configured")
		}

		r.POST("/upload", a.Upload)

		return r.Run()
	},
}
