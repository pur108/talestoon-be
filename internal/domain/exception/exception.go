package exception

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized action")
	ErrNotFound     = errors.New("resource not found")
)
