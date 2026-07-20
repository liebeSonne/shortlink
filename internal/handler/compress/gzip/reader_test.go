package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGzipReader(t *testing.T) {
	originalData := []byte("test data for compression")
	compressedData := compressData(t, originalData)

	type on struct {
		body io.ReadCloser
	}
	type want struct {
		err error
	}
	testCases := []struct {
		name string
		on   on
		want want
	}{
		{
			"success with valid gzip data",
			on{body: io.NopCloser(bytes.NewReader(compressedData))},
			want{err: nil},
		},
		{
			"error with invalid gzip data",
			on{body: io.NopCloser(bytes.NewReader([]byte("invalid data")))},
			want{err: gzip.ErrHeader},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader, err := NewGzipReader(tc.on.body)

			if tc.want.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.want.err)
				require.Nil(t, reader)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, reader)
			require.NoError(t, reader.Close())
		})
	}
}

func TestGzipReader_Read(t *testing.T) {
	originalData := []byte("test data for compression read")
	compressedData := compressData(t, originalData)

	type on struct {
		bufferSize int
	}
	type want struct {
		data []byte
		err  error
	}
	testCases := []struct {
		name string
		on   on
		want want
	}{
		{
			"read all data at once",
			on{bufferSize: len(originalData) + 10},
			want{data: originalData, err: nil},
		},
		{
			"read data in chunks",
			on{bufferSize: 5},
			want{data: nil, err: nil},
		},
		{
			"read with small buffer",
			on{bufferSize: 1},
			want{data: nil, err: nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader, err := NewGzipReader(io.NopCloser(bytes.NewReader(compressedData)))
			require.NoError(t, err)
			defer reader.Close()

			var result []byte
			buf := make([]byte, tc.on.bufferSize)

			for {
				n, readErr := reader.Read(buf)
				if n > 0 {
					result = append(result, buf[:n]...)
				}
				if readErr == io.EOF {
					break
				}
				require.NoError(t, readErr)
			}

			if tc.want.data != nil {
				assert.Equal(t, tc.want.data, result)
			} else {
				assert.Equal(t, originalData, result)
			}
		})
	}
}

func TestGzipReader_Close(t *testing.T) {
	originalData := []byte("test data for compression close")
	compressedData := compressData(t, originalData)

	type want struct {
		err error
	}
	testCases := []struct {
		name string
		want want
	}{
		{
			"close success",
			want{err: nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader, err := NewGzipReader(io.NopCloser(bytes.NewReader(compressedData)))
			require.NoError(t, err)

			err = reader.Close()
			if tc.want.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.want.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
