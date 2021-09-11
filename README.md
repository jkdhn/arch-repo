# arch-repo

Manages a [pacman](https://wiki.archlinux.org/title/Pacman) package repository in AWS S3.

## Server

Locally stores the package database.

Uploads it to S3 in the format understood by pacman.

### API

#### POST /upload

Adds the package metadata to the repository database.

Returns a presigned URL to upload the package archive to S3.

### Usage

```
Usage:
  server [flags]

Flags:
  -b, --bucket string               AWS bucket
  -c, --claims string               required JWT token claims (JSON)
      --cleanup-interval duration   Cleanup interval (default 48h0m0s)
  -e, --endpoint string             S3 endpoint
  -r, --endpoint-region string      S3 endpoint signing region
  -h, --help                        help for server
  -i, --issuer string               JWT issuer
  -k, --jwks string                 JWT key store URL
  -n, --name string                 repository name
      --skip-auth                   disable authentication
```

### Example configuration

Environment variables:

```
GIN_MODE=release
PORT=80
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
AWS_DEFAULT_REGION=eu-central-1
```

Command line:

```
--bucket example-bucket --name example
```

Authentication with GitLab CI (`$CI_JOB_JWT`), restricted to jobs from a certain group:

```
--issuer gitlab.com --jwks https://gitlab.com/-/jwks --claims '{"namespace_id": "123"}'
```

`/etc/pacman.conf`:

```
[example]
SigLevel = Never
Server = https://example-bucket.s3.eu-central-1.amazonaws.com

```

## Deploy

Command line tool to deploy a package.

Generates the package metadata and uploads it to the API server.

Uses the returned presigned URL to upload the package archive to S3.

### Usage

```
Usage:
  deploy --api-url url filename [flags]

Flags:
  -a, --api-url string      repository API URL
  -t, --auth-token string   JWT token
  -h, --help                help for deploy
```

Example: `deploy -a https://example.com/upload example.pkg.tar.zst`
