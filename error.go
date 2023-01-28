package errors

import (
	"runtime"

	"github.com/pkg/errors"
)

type CutDirection uint

const (
	_ CutDirection = iota
	Up
	Down
)

type Error interface {
	// Error implements error
	error

	Because(fn interface{}, err error) Error
	Cut(dir CutDirection) Error
	IsNil() bool
	Cause() Cause
	Stack() Stack
	// StackTrace() errors.StackTrace
	// Trace() Error

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

func Cut(e *Err, dir CutDirection) Error {
	// var pcs [1]uintptr
	// runtime.Callers(3, pcs[:])
	var pc, _, _, ok = runtime.Caller(2)
	if !ok {
		return e
	}

	if e == nil {
		return nil
	}

	var frames Stack
	switch dir {
	case Up:
		for i := range e.frames {
			if uintptr(e.frames[i]) == pc {
				e.frames = e.frames[:i]
			}
		}

	case Down:
		for i := range e.frames {
			if uintptr(e.frames[i]) == pc {
				e.frames = e.frames[i:]
			}
		}
	}

	return &Err{
		err:    e.err,
		cause:  e.cause,
		frames: frames,
	}
}
