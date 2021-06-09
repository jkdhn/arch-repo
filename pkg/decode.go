package pkg

import (
	"arch-repo/pkg/desc"
	"arch-repo/pkg/info"
	"arch-repo/util"
	"archive/tar"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/klauspost/compress/zstd"
	"io"
)

type Decoder struct {
	in       io.Reader
	filename string
}

func NewDecoder(in io.Reader, filename string) *Decoder {
	return &Decoder{
		in:       in,
		filename: filename,
	}
}

func (d *Decoder) Decode() (*Package, error) {
	input := util.NewCounter(d.in)
	md5Hash := md5.New()
	sha256Hash := sha256.New()
	tarIn, tarOut := io.Pipe()
	writer := io.MultiWriter(md5Hash, sha256Hash, tarOut)

	defer func() {
		_ = tarIn.Close()
	}()

	go func() {
		_, err := io.Copy(writer, input)
		_ = tarOut.CloseWithError(err)
	}()

	decoder, err := zstd.NewReader(tarIn)
	if err != nil {
		return nil, err
	}

	defer decoder.Close()

	var description *desc.Description
	tarReader := tar.NewReader(decoder)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if header.Name == ".PKGINFO" {
			i, err := info.NewDecoder(tarReader).Decode()
			if err != nil {
				return nil, err
			}

			description = &desc.Description{
				Info: *i,
			}
		}
	}

	if description == nil {
		return nil, fmt.Errorf(".PKGINFO not found")
	}

	description.FileName = d.filename
	description.CompressedSize = input.Length()
	copy(description.MD5Sum[:], md5Hash.Sum(nil))
	copy(description.SHA256Sum[:], sha256Hash.Sum(nil))

	p := &Package{
		description: description,
	}

	return p, nil
}
