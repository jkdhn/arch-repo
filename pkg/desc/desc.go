package desc

import (
	"arch-repo/pkg/info"
	"crypto/md5"
	"crypto/sha256"
)

type Description struct {
	info.Info
	FileName       string            `pkgdesc:"FILENAME" json:"file_name"`
	CompressedSize uint64            `pkgdesc:"CSIZE" json:"compressed_size"`
	MD5Sum         [md5.Size]byte    `pkgdesc:"MD5SUM" json:"md5_sum"`
	SHA256Sum      [sha256.Size]byte `pkgdesc:"SHA256SUM" json:"sha256_sum"`
	PGPSignature   string            `pkgdesc:"PGPSIG" json:"pgp_signature"`
}
