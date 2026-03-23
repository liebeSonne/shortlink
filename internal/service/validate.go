package service

import (
	"errors"
	"log"
	"net/url"
)

var ErrEmptyURL = errors.New("empty URL")
var ErrInvalidURL = errors.New("invalid URL")

func validateLink(str string) error {
	if len(str) == 0 {
		return ErrEmptyURL
	}

	u, err := url.ParseRequestURI(str)
	if err != nil {
		log.Printf("error on parse url: %v", err)
		return ErrInvalidURL
	}
	if u.Scheme == "" || u.Host == "" {
		return ErrInvalidURL
	}
	return nil
}
