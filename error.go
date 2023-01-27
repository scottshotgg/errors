package errors

import (
	"github.com/pkg/errors"
)

type Error interface {
	error

	Because(fn interface{}, err error) *Err
	Cut() *Err
	IsNil() bool
	Cause() *Cause
	Stack() Stack
	StackTrace() errors.StackTrace
	Trace() *Err

	Is(err error) bool
	As(target any) bool
}
