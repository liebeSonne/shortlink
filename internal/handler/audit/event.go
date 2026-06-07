package audit

import (
	"time"
)

// Action - тип действия.
type Action int

// Значения типов действий.
const (
	ActionShorted Action = iota
	ActionFollow
)

// Event - данные события для аудита.
type Event struct {
	Time   time.Time `json:"ts"`
	Action Action    `json:"action"`
	UserID *string   `json:"user_id,omitempty"`
	URL    string    `json:"url"`
}
