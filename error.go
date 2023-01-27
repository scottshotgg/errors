package errors

import (
	"github.com/pkg/errors"
)

type Error interface {
	// Error implements error
	error

	Because(fn interface{}, err error) Error
	Cut() Error
	IsNil() bool
	Cause() Cause
	Stack() Stack
	StackTrace() errors.StackTrace
	Trace() Error

	Is(err error) bool
	As(target any) bool

	// TODO: think about a FromStatus(status.Status) Error function
}

// Only used for the type
var e Error

// TODO: make sure this function works
func FromError(err error) Error {
	if err == nil {
		return nil
	}

	if !errors.Is(err, &Err{}) {
		return New(err)
	}

	return err.(*Err)
}

// TODO: think about branches
// func Chain(err error) Error {
// 	return FromError(err).Append()
// }
