package audit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/liebeSonne/shortlink/internal/logger"
)

func TestPublisher(t *testing.T) {
	stringUserID := "u1"

	type on struct {
		events []Event
	}
	type want struct {
		countObserverEvents int
	}
	testCases := []struct {
		name string
		on   on
		want want
	}{
		{
			"no events",
			on{[]Event{}},
			want{0},
		},
		{
			"one events",
			on{[]Event{
				{Time: time.Now(), Action: ActionShorted, UserID: &stringUserID, URL: "url1"},
			}},
			want{1},
		},
		{
			"many events",
			on{[]Event{
				{Time: time.Now(), Action: ActionShorted, UserID: &stringUserID, URL: "url1"},
				{Time: time.Now(), Action: ActionFollow, UserID: nil, URL: "url2"},
			}},
			want{2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewMockObserver(t)
			countObserverEvents := 0
			o.EXPECT().Update(mock.Anything).Run(func(_ Event) {
				countObserverEvents++
			}).Maybe()
			l := logger.NewMockLogger(t)
			l.EXPECT().Debugf(mock.Anything, mock.Anything).Maybe()

			p := NewPublisher(t.Context(), 1, 1, l)
			p.Subscribe(o)

			for _, event := range tc.on.events {
				p.Notify(event)
			}

			time.Sleep(200 * time.Millisecond)

			assert.Equal(t, tc.want.countObserverEvents, countObserverEvents)
		})
	}
}
