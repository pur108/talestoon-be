package exception

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized action")
	ErrForbidden    = errors.New("forbidden action")
	ErrNotFound     = errors.New("resource not found")
)
