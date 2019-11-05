// Copyright 2018 Owkin, inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	// Associated interface through errors methods
	context map[string]interface{}
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
//	errors.Kind
//		The class of error, such as a key conflict
//	error
//		The underlying error
//	string
//		The string to add to the existing error message. As mention above
//		all the args following the first string will be handle as format
//		parameters.
func E(args ...interface{}) Error {
	e := Error{}
	e.context = map[string]interface{}{}
	for i, arg := range args {
		switch arg := arg.(type) {
		case Kind:
			e.Kind = arg
		case Error:
			e.context = arg.context
			if e.Err == nil {
				e.Err = arg.Err
			} else {
				e.Err = fmt.Errorf("%s %s", arg.Error(), e.Error())
			}
		case error:
			if e.Err == nil {
				e.Err = arg
			} else {
				e.Err = fmt.Errorf("%s %s", arg.Error(), e.Error())
			}
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

// Internal returns an Error of a this specific type
func Internal(args ...interface{}) Error {
	args = append([]interface{}{internal}, args...)
	return E(args...)
}

// NotFound returns an Error of a this specific type
func NotFound(args ...interface{}) Error {
	args = append([]interface{}{notFound}, args...)
	return E(args...)
}

// Conflict returns an Error of a this specific type
func Conflict(args ...interface{}) Error {
	args = append([]interface{}{conflict}, args...)
	return E(args...)
}

// BadRequest returns an Error of a this specific type
func BadRequest(args ...interface{}) Error {
	args = append([]interface{}{badRequest}, args...)
	return E(args...)
}

// Forbidden returns an Error of a this specific type
func Forbidden(args ...interface{}) Error {
	args = append([]interface{}{forbidden}, args...)
	return E(args...)
}

// WithKey associate the given key to the error context
// It overwrites previous key if any.
func (e Error) WithKey(key string) Error {
	e.context["key"] = key
	return e
}

// WithKeys associate the given keys to the error context
// It overwrites previous keys' list if any.
func (e Error) WithKeys(keys []string) Error {
	e.context["keys"] = keys
	return e
}

// GetContext return the associated key if there is any
func (e Error) GetContext() map[string]interface{} {
	return e.context
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
	forbidden              // Forbidden request
)

// HTTPStatusCode returns for an error kind the associated http status
func (k Kind) HTTPStatusCode() int {
	switch k {
	case notFound:
		return http.StatusNotFound
	case conflict:
		return http.StatusConflict
	case badRequest:
		return http.StatusBadRequest
	case forbidden:
		return http.StatusForbidden
	}
	return http.StatusInternalServerError
}
