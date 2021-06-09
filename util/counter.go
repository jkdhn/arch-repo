package util

import "io"

type Counter struct {
	reader io.Reader
	v      uint64
}

func NewCounter(reader io.Reader) *Counter {
	return &Counter{
		reader: reader,
	}
}

func (c *Counter) Read(p []byte) (n int, err error) {
	n, err = c.reader.Read(p)
	if err == nil {
		c.v += uint64(n)
	}
	return
}

func (c *Counter) Length() uint64 {
	return c.v
}

var _ io.Reader = (*Counter)(nil)
