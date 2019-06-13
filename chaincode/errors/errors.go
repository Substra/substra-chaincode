package errors

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
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

// E return an error according to the args passed.
// The type of the arg determines its meaning.
// If more than one arg of a type is passed (string put aside)
// Only the last one is applied. In case of string, when encountering one
// all the other args are passed as parameters for the string formatter.
//
// The possible arg type are:
//	errors.Error
//		It will be copied
//	errors.Kind
//		The class of error, such as a key conflict
//	error
//		The underlying error
//	string
//		The string to add to the existing error message. As mention above
//		all the args following the first string will be handle as format
//		parameters.
func E(args ...interface{}) error {
	e := Error{}
	for i, arg := range args {
		switch arg := arg.(type) {
		case Error:
			e = arg
		case Kind:
			e.Kind = arg
		case error:
			if e.Err == nil {
				e.Err = arg
				continue
			}
			e.Err = fmt.Errorf("%s %s", arg.Error(), e.Error())
		case string:
			parameters := args[i+1:]
			msg := fmt.Sprintf(arg, parameters...)
			if e.Err == nil {
				e.Err = fmt.Errorf(msg)
			} else {
				e.Err = fmt.Errorf("%s %s", msg, e.Error())
			}
			return e
		}
	}
	return e
}

// Wrap converts an error interface to the internal Error type.
// Does nothing if an internal Error type is passed.
func Wrap(err error) Error {
	if e, ok := err.(Error); ok {
		return e
	}
	return Error{Err: err}
}

// NotFound returns an Error of a this specific type
func NotFound(args ...interface{}) error {
	args = append([]interface{}{notFound}, args...)
	return E(args...)
}

// Conflict returns an Error of a this specific type
func Conflict(args ...interface{}) error {
	args = append([]interface{}{conflict}, args...)
	return E(args...)
}

// BadRequest returns an Error of a this specific type
func BadRequest(args ...interface{}) error {
	args = append([]interface{}{badRequest}, args...)
	return E(args...)
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
	internal   Kind = iota // default error
	notFound               // Asset has not been found
	conflict               // Asset already exists
	badRequest             // Invalid request
)

// HTTPStatusCode returns for an error kind the associated http status
func (k Kind) HTTPStatusCode() int {
	switch k {
	case internal:
	case notFound:
		return http.StatusNotFound
	case conflict:
		return http.StatusConflict
	case badRequest:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
