package gzip

import (
	"compress/gzip"
	"io"
)

type Reader interface {
	io.ReadCloser
}

func NewGzipReader(r io.ReadCloser) (Reader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:  r,
		zr: zr,
	}, nil
}

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (r *gzipReader) Read(p []byte) (n int, err error) {
	return r.zr.Read(p)
}

func (r *gzipReader) Close() error {
	if err := r.r.Close(); err != nil {
		return err
	}
	return r.zr.Close()
}
