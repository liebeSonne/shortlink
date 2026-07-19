package subnet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/logger"
)

func TestNewTrustedSubnetMiddleware(t *testing.T) {
	type when struct {
		trustedSubnet string
		xRealIP       string
	}
	type want struct {
		code int
	}
	testCases := []struct {
		name string
		when when
		want want
	}{
		{
			name: "successful access with valid IP in subnet",
			when: when{"192.168.1.0/24", "192.168.1.100"},
			want: want{http.StatusOK},
		},
		{
			name: "successful access with valid IP in subnet (IPv6)",
			when: when{"2001:db8::/32", "2001:db8::1"},
			want: want{http.StatusOK},
		},
		{
			name: "forbidden when trusted subnet is empty",
			when: when{"", "192.168.1.100"},
			want: want{http.StatusForbidden},
		},
		{
			name: "forbidden when X-Real-IP header is missing",
			when: when{"192.168.1.0/24", ""},
			want: want{http.StatusForbidden},
		},
		{
			name: "forbidden when IP is not in trusted subnet",
			when: when{"192.168.1.0/24", "10.0.0.1"},
			want: want{http.StatusForbidden},
		},
		{
			name: "forbidden when invalid CIDR in trusted subnet",
			when: when{"invalid-cidr", "192.168.1.100"},
			want: want{http.StatusForbidden},
		},
		{
			name: "forbidden when invalid X-Real-IP format",
			when: when{"192.168.1.0/24", "invalid-ip"},
			want: want{http.StatusForbidden},
		},
		{
			name: "forbidden when X-Real-IP has extra spaces",
			when: when{"192.168.1.0/24", " 192.168.1.100 "},
			want: want{http.StatusOK},
		},
		{
			name: "forbidden when IP from different subnet",
			when: when{"10.0.0.0/8", "192.168.1.100"},
			want: want{http.StatusForbidden},
		},
		{
			name: "successful access with IP at subnet boundaries",
			when: when{"192.168.1.0/24", "192.168.1.0"},
			want: want{http.StatusOK},
		},
		{
			name: "successful access with IP at subnet broadcast",
			when: when{"192.168.1.0/24", "192.168.1.255"},
			want: want{http.StatusOK},
		},
		{
			name: "forbidden with IPv6 address in IPv4 subnet",
			when: when{"192.168.1.0/24", "2001:db8::1"},
			want: want{http.StatusForbidden},
		},
		{
			name: "forbidden with IPv4 address in IPv6 subnet",
			when: when{"2001:db8::/32", "192.168.1.100"},
			want: want{http.StatusForbidden},
		},
		{
			name: "successful access with private IPv6",
			when: when{"fd00::/8", "fd12:3456:789a::1"},
			want: want{http.StatusOK},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			l := logger.NewMockLogger(t)
			l.On("Errorf", mock.Anything, mock.Anything).Maybe()
			l.On("Warnf", mock.Anything, mock.Anything).Maybe()

			handler := NewTrustedSubnetMiddleware(nextHandler, tc.when.trustedSubnet, l)

			srv := httptest.NewServer(handler)
			defer srv.Close()

			client := resty.New()

			req := client.R()
			req.Method = http.MethodGet
			req.URL = srv.URL + "/test"
			if tc.when.xRealIP != "" {
				req.Header.Set("X-Real-IP", tc.when.xRealIP)
			}

			resp, err := req.Send()
			require.NoError(t, err)

			require.Equal(t, tc.want.code, resp.StatusCode(), fmt.Sprintf("expected status code %d but got %d with body: %s", tc.want.code, resp.StatusCode(), string(resp.Body())))
		})
	}
}
