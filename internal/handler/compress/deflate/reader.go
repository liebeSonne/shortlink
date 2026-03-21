package deflate

import (
	"compress/zlib"
	"io"
)

type Reader interface {
	io.ReadCloser
}

func NewDeflateReader(r io.ReadCloser) (Reader, error) {
	zr, err := zlib.NewReader(r)
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
	zr io.ReadCloser
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
