package audit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/logger"
)

func TestNewURLObserver_Update(t *testing.T) {
	time1 := time.Now().Truncate(time.Second)
	time2 := time.Now().Truncate(time.Second).Add(-20 * time.Second)
	url1 := "https://example1.com"
	url2 := "https://example2.com"
	userID1 := "user1"

	type on struct {
		events []Event
	}
	testCases := []struct {
		name string
		on   on
	}{
		{
			"empty events",
			on{events: []Event{}},
		},
		{
			"not empty events",
			on{
				events: []Event{
					{
						Time:   time1,
						URL:    url1,
						UserID: nil,
						Action: ActionShorted,
					},
					{
						Time:   time2,
						URL:    url2,
						UserID: &userID1,
						Action: ActionFollow,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			countEvents := 0
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodPost, r.Method)
				countEvents++
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
			}))
			defer mockServer.Close()

			l := logger.NewMockLogger(t)
			observer := NewURLObserver(mockServer.URL, l)

			for _, event := range tc.on.events {
				observer.Update(event)
			}

			require.Equal(t, countEvents, len(tc.on.events))
		})
	}
}
