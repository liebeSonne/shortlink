package audit

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/logger"
)

func TestFileObserver_Update(t *testing.T) {
	time1 := time.Now().Truncate(time.Second)
	time2 := time.Now().Truncate(time.Second).Add(-20 * time.Second)
	url1 := "https://example1.com"
	url2 := "https://example2.com"
	userID1 := "user1"

	type on struct {
		events []Event
	}
	tests := []struct {
		name string
		on   on
	}{
		{
			name: "not empty events",
			on: on{
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempFile, err := os.CreateTemp("", "test_observer_*.log")
			require.NoError(t, err)
			defer tempFile.Close()
			filePath := tempFile.Name()

			l := logger.NewMockLogger(t)

			observer, err := NewFileObserver(filePath, l)
			t.Cleanup(func() {
				if observer != nil {
					err = observer.Close()
					require.NoError(t, err)
				}
			})

			require.NoError(t, err)

			var wg sync.WaitGroup
			for i := 0; i < 3; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for _, event := range tc.on.events {
						observer.Update(event)
					}
				}()
			}
			wg.Wait()

			if len(tc.on.events) > 0 {
				fileContent, err := os.ReadFile(filePath)
				assert.NoError(t, err)
				assert.NotEmpty(t, fileContent)
			}
		})
	}
}
