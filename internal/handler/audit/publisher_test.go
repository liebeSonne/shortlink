package audit

import "testing"

func TestPublisher(t *testing.T) {
	type on struct {
		event Event
	}
	type want struct {
	}
	testCases := []struct {
		name string
		on   on
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := NewMockObserver(t)
			o.EXPECT().Update(tc.on.event)

			p := NewPublisher()
			p.Subscribe(o)
			p.Notify(tc.on.event)
		})
	}
}
