package audit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go"

	"github.com/liebeSonne/shortlink/internal/logger"
)

var errUnexpectedStatusCode = errors.New("unexpected status code")

func NewURLObserver(
	host string,
	maxAttempts uint,
	delay time.Duration,
	logger logger.Logger,
) *URLObserver {
	return &URLObserver{
		host:        host,
		maxAttempts: maxAttempts,
		delay:       delay,
		logger:      logger,
	}
}

type URLObserver struct {
	host        string
	maxAttempts uint
	delay       time.Duration
	logger      logger.Logger
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

	return retry.Do(
		func() error {
			resp, err := http.Post(o.host, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				return fmt.Errorf("failed to send event: %w", err)
			}
			defer resp.Body.Close()

			return o.translateStatusCode(resp.StatusCode)
		},
		retry.Attempts(o.maxAttempts),
		retry.Delay(o.delay),
		retry.RetryIf(func(err error) bool {
			return errors.Is(err, errUnexpectedStatusCode)
		}),
	)
}

func (o *URLObserver) translateStatusCode(statusCode int) error {
	if statusCode > http.StatusBadRequest {
		return fmt.Errorf("unexpected status code: %d: %w", statusCode, errUnexpectedStatusCode)
	}
	return nil
}
