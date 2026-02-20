package handler

import (
	"errors"
	"net/url"
)

var ErrEmptyID = errors.New("empty ID")
var ErrEmptyURL = errors.New("empty URL")
var ErrInvalidURL = errors.New("invalid URL")

func validateShortLinkID(shortLinkID string) error {
	if len(shortLinkID) == 0 {
		return ErrEmptyID
	}
	return nil
}

func validateLink(str string) error {
	if len(str) == 0 {
		return ErrEmptyURL
	}

	u, err := url.ParseRequestURI(str)
	if err != nil {
		return err
	}
	if u.Scheme == "" || u.Host == "" {
		return ErrInvalidURL
	}
	return nil
}
