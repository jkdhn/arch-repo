package repository

import (
	"arch-repo/pkg/desc"
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

type Encoder struct {
	out io.Writer
}

func NewEncoder(out io.Writer) *Encoder {
	return &Encoder{
		out: out,
	}
}

func (e *Encoder) Encode(r *Repository) error {
	writer := gzip.NewWriter(e.out)
	tarWriter := tar.NewWriter(writer)

	defer func() {
		_ = tarWriter.Close()
		_ = writer.Close()
	}()

	if err := r.storage.ForEach(func(description *desc.Description) error {
		folder := fmt.Sprintf("%v-%v/", description.Name, description.Version)

		if err := tarWriter.WriteHeader(&tar.Header{
			Typeflag: tar.TypeDir,
			Name:     folder,
			Mode:     0755,
		}); err != nil {
			return err
		}

		data := new(bytes.Buffer)
		if err := desc.NewEncoder(data).Encode(description); err != nil {
			return err
		}

		if err := tarWriter.WriteHeader(&tar.Header{
			Typeflag: tar.TypeReg,
			Name:     folder + "desc",
			Size:     int64(data.Len()),
			Mode:     0644,
		}); err != nil {
			return err
		}

		if _, err := io.Copy(tarWriter, data); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if err := tarWriter.Close(); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}
