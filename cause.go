package errors

import "runtime"

type Cause struct {
	name  string
	err   error
	value *runtime.Func
}

func (c *Cause) Error() string {
	return c.err.Error()
}

func (c *Cause) Name() string {
	return c.name
}
