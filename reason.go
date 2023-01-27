package errors

import (
	"runtime"
)

// Reason implements Cause
var _ Cause = (*Reason)(nil)

func NewReason(name string, err error, value *runtime.Func) *Reason {
	return &Reason{
		name:  name,
		err:   err,
		value: value,
	}
}

type Reason struct {
	name  string
	err   error
	value *runtime.Func
}

func (c *Reason) Error() string {
	if c == nil || c.err == nil {
		return "<nil>"
	}

	return c.err.Error()
}

func (c *Reason) Name() string {
	if c == nil {
		return ""
	}

	return c.name
}
