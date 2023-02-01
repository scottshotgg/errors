package errors

// Reason implements Cause
var _ Cause = (*Reason)(nil)

func NewReason(err error, name string) *Reason {
	return &Reason{
		name: name,
		err:  err,
		// value: fn,
	}
}

type Reason struct {
	name string
	err  error
	// value *runtime.Func
}

func (c *Reason) Error() string {
	if c == nil || c.err == nil {
		return ""
	}

	return c.err.Error()
}

func (c *Reason) Name() string {
	if c == nil {
		return ""
	}

	return c.name
}
