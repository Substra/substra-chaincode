package main

import (
	"fmt"
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
func NewError(args ...interface{}) error {
	e := Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Error:
			e = arg
		case Kind:
			e.Kind = arg
		case string:
			if e.Err == nil {
				e.Err = fmt.Errorf(arg)
				continue
			}
			e.Err = fmt.Errorf("%s %s", arg, e.Error())
		case error:
			if e.Err == nil {
				e.Err = arg
				continue
			}
			e.Err = fmt.Errorf("%s %s", arg.Error(), e.Error())
		}
	}
	return e
}

// HTTPStatusCode wrap the HTTPStatusCode methods of the Kind parameter
func (e Error) HTTPStatusCode() int {
	return e.Kind.HTTPStatusCode()
}

// Kind is the type use to discriminate between errors type
// It's not intended for user print but to handle errors correctly
type Kind uint8

// Possible errors kinds. Beware, this declaration is order sensitive.
const (
	Default    Kind = iota // default unrecognized error
	NotFound               // Asset has not been found
	Conflict               // Asset already exists
	BadRequest             // Invalid request
)

// HTTPStatusCode returns for an error kind the associated http status
func (k Kind) HTTPStatusCode() int {
	switch k {
	case Default:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case Conflict:
		return http.StatusConflict
	case BadRequest:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
