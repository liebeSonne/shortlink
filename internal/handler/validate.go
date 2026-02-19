package handler

import "errors"

var ErrEmptyID = errors.New("empty ID")
var ErrEmptyURL = errors.New("empty URL")

func validateShortLinkID(shortLinkID string) error {
	if len(shortLinkID) == 0 {
		return ErrEmptyID
	}
	return nil
}

func validateLink(url string) error {
	if len(url) == 0 {
		return ErrEmptyURL
	}
	return nil
}
