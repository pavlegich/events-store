// Package errors contains variables with error explainations.
package errors

import "errors"

var (
	ErrEventNotFound      = errors.New("event not found")
	ErrEventAlreadyExists = errors.New("event already exists")
)
