# arch-repo

Manages a [pacman](https://wiki.archlinux.org/title/Pacman) package repository in AWS S3.

## Server

Locally stores the package database.

Uploads it to S3 in the format understood by pacman.

### API

#### POST /upload

Adds the package metadata to the repository database.

Returns a presigned URL to upload the package archive to S3.

## Setup

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
