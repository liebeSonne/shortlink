package repository

import "fmt"

// ErrConflictURL - ошибка конфликта сохранения сокращенной ссылки.
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

// NewErrConflictURL - создание экземпляра ошибки конфликта сохранения сокращенной ссылки.
func NewErrConflictURL(url string, err error) error {
	return &ErrConflictURL{URL: url, Err: err}
}
