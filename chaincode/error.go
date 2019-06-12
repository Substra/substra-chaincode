package main

import (
	"net/http"
)

// Error implements the error interface
// It Wrap the error with additional data, each in its own type
type Error struct {
	// Kind differentiate the different error type
	Kind Kind
	// The underlying error if any
	Err error
}

func (e Error) Error() string {
	return e.Err.Error()
}

// NewError return a error according to the args passed
func NewError(err error, kind Kind) error {
	return Error{Kind: kind, Err: err}
}

// Kind is the type use to discriminate between errors type
// It's not intended for user print but to handle errors correctly
type Kind uint8

// Possible errors kinds. Beware, this declaration is order sensitive.
const (
	Default  Kind = iota // default unrecognized error
	NotFound             // Assets has not been found
)

// StatusCode returns for an error kind the associated http status
func (k Kind) StatusCode() int {
	switch k {
	case Default:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
