package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/liebeSonne/shortlink/internal/logger"
)

func NewURLObserver(
	host string,
	logger logger.Logger,
) *URLObserver {
	return &URLObserver{
		host:   host,
		logger: logger,
	}
}

type URLObserver struct {
	host   string
	logger logger.Logger
}

func (o *URLObserver) Update(event Event) {
	err := o.sendEvent(event)
	if err != nil {
		o.logger.Errorf("failed to send event: %v", err)
	}
}

func (o *URLObserver) sendEvent(event Event) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = http.Post(o.host, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

	return nil
}
