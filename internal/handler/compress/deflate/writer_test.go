package deflate

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewDeflateWriter(t *testing.T) {
	type on struct {
		w http.ResponseWriter
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
			"mock writer", on{new(mockResponseWriter)}, want{nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			writer, err := NewDeflateWriter(tc.on.w)

			if tc.want.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, writer)
		})
	}
}

func TestZlibWriter_Header(t *testing.T) {
	type when struct {
		header http.Header
	}
	type want struct {
		header http.Header
	}
	testCases := []struct {
		name string
		when when
		want want
	}{
		{"empty http header", when{http.Header{}}, want{http.Header{}}},
		{"not empty http header", when{http.Header{
			"Accept-Encoding":  []string{"deflate", "111"},
			"Content-Encoding": []string{"deflate", "222"},
		}}, want{http.Header{
			"Accept-Encoding":  []string{"deflate", "111"},
			"Content-Encoding": []string{"deflate", "222"},
		}}},
		{"empty map header", when{map[string][]string{}}, want{map[string][]string{}}},
		{"not map header", when{map[string][]string{
			"Accept-Encoding":  []string{"deflate", "111"},
			"Content-Encoding": []string{"deflate", "222"},
		}}, want{map[string][]string{
			"Accept-Encoding":  []string{"deflate", "111"},
			"Content-Encoding": []string{"deflate", "222"},
		}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := new(mockResponseWriter)
			w.On("Header").Return(tc.when.header)

			writer, err := NewDeflateWriter(w)
			require.NoError(t, err)
			require.NotNil(t, writer)

			h := writer.Header()
			assert.Equal(t, tc.want.header, h)
		})
	}
}

func TestDeflateWriter_WriteHeader(t *testing.T) {
	type on struct {
		statusCode int
	}
	type want struct {
		header map[string]string
	}
	testCases := []struct {
		name string
		on   on
		want want
	}{
		{"status ok", on{http.StatusOK}, want{map[string]string{"Content-Encoding": "deflate"}}},
		{"status crated", on{http.StatusCreated}, want{map[string]string{"Content-Encoding": "deflate"}}},
		{"status found", on{http.StatusFound}, want{map[string]string{}}},
		{"status bad request", on{http.StatusBadRequest}, want{map[string]string{}}},
		{"status server error", on{http.StatusInternalServerError}, want{map[string]string{}}},
		{"zero", on{0}, want{map[string]string{"Content-Encoding": "deflate"}}},
		{"negative", on{-1}, want{map[string]string{"Content-Encoding": "deflate"}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			header := http.Header{}
			w := new(mockResponseWriter)
			w.On("WriteHeader", tc.on.statusCode).Once()
			w.On("Header").Return(header)

			writer, err := NewDeflateWriter(w)
			require.NoError(t, err)
			require.NotNil(t, writer)

			writer.WriteHeader(tc.on.statusCode)
			for k, v := range tc.want.header {
				assert.Equal(t, v, writer.Header().Get(k))
			}
		})
	}
}

func TestDeflateWriter_Write(t *testing.T) {
	err1 := errors.New("error 1")

	type on struct {
		bytes []byte
	}
	type when struct {
		len int
		err error
	}
	type want struct {
		err error
	}
	testCases := []struct {
		name string
		on   on
		when when
		want want
	}{
		{"empty", on{[]byte("")}, when{0, nil}, want{nil}},
		{"not empty", on{[]byte("not empty")}, when{10, nil}, want{nil}},
		{"error", on{[]byte("123")}, when{10, err1}, want{err1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := new(mockResponseWriter)
			w.On("Write", mock.Anything).Return(tc.when.len, tc.when.err)

			writer, err := NewDeflateWriter(w)
			require.NoError(t, err)
			require.NotNil(t, writer)

			count, err := writer.Write(tc.on.bytes)
			if tc.want.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, count, 0)
		})
	}
}

func TestDeflateWriter_Close(t *testing.T) {
	err1 := errors.New("error 1")

	type want struct {
		err error
	}
	type when struct {
		len int
		err error
	}
	testCases := []struct {
		name string
		when when
		want want
	}{
		{"empty", when{0, nil}, want{nil}},
		{"not empty", when{123, nil}, want{nil}},
		{"error", when{123, err1}, want{err1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := new(mockResponseWriter)
			w.On("Write", mock.Anything).Return(tc.when.len, tc.when.err)

			writer, err := NewDeflateWriter(w)
			require.NoError(t, err)
			require.NotNil(t, writer)

			err = writer.Close()
			if tc.want.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
		})
	}
}
