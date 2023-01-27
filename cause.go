package errors

type Cause interface {
	error

	Name() string

	// Could have multiple reasons here?
}
