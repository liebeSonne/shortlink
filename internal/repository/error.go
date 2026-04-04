package repository

import "fmt"

type ErrConflictURL struct {
	URL string
	Err error
}

func (e *ErrConflictURL) Error() string {
	return fmt.Sprintf("conflict url '%s': %s", e.URL, e.Err)
}

func (e *ErrConflictURL) Unwrap() error {
	return e.Err
}

func NewErrConflictURL(url string, err error) error {
	return &ErrConflictURL{URL: url, Err: err}
}
