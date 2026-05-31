package audit

import (
	"time"
)

type Action int

const (
	ActionShorted Action = iota
	ActionFollow
)

type Event struct {
	Time   time.Time `json:"ts"`
	Action Action    `json:"action"`
	UserID *string   `json:"user_id,omitempty"`
	URL    string    `json:"url"`
}
