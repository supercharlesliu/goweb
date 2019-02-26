package goweb

import (
	"errors"
)

var (
	// ErrUnknowContentType on header "Content-Type" present.
	ErrUnknowContentType = errors.New("Content-Type header not found")

	ErrMethodNotFound = errors.New("HTTP")
)
