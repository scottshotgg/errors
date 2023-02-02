package errors

import (
	stderrors "errors"
	"fmt"

	errors_pb "github.com/scottshotgg/errors/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Opt int

const (
	_ Opt = iota

	WithCause
	WithStack
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

	Err() error
	Cause() Cause
	Stack() Stack
	Details() map[string]any

	Is(err error) bool
	As(target any) bool

	// StackTrace() errors.StackTrace
	// Trace() Error
	// Causes() []Cause

	// TODO: think about a FromStatus(status.Status) Error function
}

// TODO: make sure this function works
func FromError(err error) Error {
	if err == nil {
		return nil
	}

	if !stderrors.Is(err, &Err{}) {
		return New(err)
	}

	return err.(*Err)
}

// TODO: think about branches
// func Chain(err error) Error {
// 	return FromError(err).Append()
// }

func Cut(e Error, dir CutDirection) Error {
	if e == nil || dir > Down {
		return nil
	}

	return &Err{
		err:   e.Err(),
		cause: e.Cause(),
		stack: e.Stack().CutAt(3, dir),
	}
}

// TODO: make a clone function; probably on implementation
func Clone(e Error, dir CutDirection) Error {
	if e == nil || dir > Down {
		return nil
	}

	return &Err{
		err:   e.Err(),
		cause: e.Cause(),
		stack: e.Stack().CutAt(3, dir),
	}
}

func ToStatus(e Error, opts ...Opt) *status.Status {
	if e == nil {
		return status.New(codes.OK, codes.OK.String())
	}

	// TODO: need to expand on this, possibly nil check
	var (
		st = status.New(codes.Unknown, e.Err().Error())

		pbErr *errors_pb.Error
	)

	for _, opt := range opts {
		switch opt {
		case WithCause:
			var c = e.Cause()
			if c != nil {
				if pbErr == nil {
					pbErr = &errors_pb.Error{}
				}

				pbErr.Cause = &errors_pb.Cause{
					Name:  c.Name(),
					Error: c.Error(),
				}
			}

		case WithStack:
			if pbErr == nil {
				pbErr = &errors_pb.Error{}
			}

			pbErr.Stack = e.Stack().ToPB()

		default:
			return status.New(codes.Internal, fmt.Sprintf("unknown opt type: %d", opt))
		}
	}

	if pbErr != nil {
		var err error

		st, err = st.WithDetails(pbErr)
		if err != nil {
			return nil
		}
	}

	return st
}

func FromStatus(s *status.Status) (Error, error) {
	if s == nil {
		return nil, nil
	}

	if s.Code() == codes.OK {
		return nil, nil
	}

	var e = &Err{
		err: stderrors.New(s.Message()),
	}

	var details = s.Details()
	if len(details) == 0 {
		return e, nil
	}

	var pbErr, ok = details[0].(*errors_pb.Error)
	if !ok {
		return nil, fmt.Errorf("details were not a valid Error: %T", details[0])
	}

	if pbErr.Cause != nil {
		e.cause = &Reason{
			name: pbErr.Cause.Name,
			err:  stderrors.New(pbErr.Cause.Error),
		}
	}

	if len(pbErr.Stack) > 0 {
		e.stack = StackFromPB(pbErr.Stack)
	}

	return e, nil
}
