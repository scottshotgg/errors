package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/scottshotgg/errors"
)

func TestErrors(t *testing.T) {
	var err error = someFunc()

	// var e errors.Error
	// if stderrors.Is(err, e) {
	// 	fmt.Println("it is that type")
	// }

	var se = errors.FromError(err).Cut(errors.Down)
	fmt.Println("se:", se)
	fmt.Println("se:", se.Stack())
	// fmt.Println(se.Cut(errors.Up).Stack())
	// fmt.Println(se.Cut(errors.Down).Stack())

	// se.Stack().Frames()

	var st = errors.ToStatus(se, errors.WithStack, errors.WithCause)
	fmt.Println("status:", st)
	fmt.Println("status.Details():", st.Details())

	e2, err := errors.FromStatus(st)
	if err != nil {
		panic(err)
	}

	fmt.Println("e2:", e2)
	fmt.Println("e2.cause:", e2.Cause())
	fmt.Println("e2.stack:", e2.Stack())
}

func someFunc() errors.Error {
	var err = someOtherFunc1()
	stack := err.Stack()
	fmt.Println("stack.up:", stack.Cut(errors.Up))
	fmt.Println("stack.down:", stack.Cut(errors.Down))

	fmt.Println("pkg.up:", errors.Cut(err, errors.Up).Stack())
	fmt.Println("pkg.down:", errors.Cut(err, errors.Down).Stack())

	// fmt.Println("err.up:", err.Cut(errors.Up).Stack())
	// fmt.Println("err.down:", err.Cut(errors.Down).Stack())

	return err
}

func someOtherFunc1() errors.Error {
	var err = someOtherFunc2()
	return errors.New(stderrors.New("something happened")).Because(someOtherFunc2, err)
}

var (
	randomLibError = stderrors.New("123 what the func")
)

func someOtherFunc2() error {
	return randomLibError
}

/*
	Error could wrap a multierror of Errors

	Error.New().Cut():
		multierror:
			Error1.Cut()
			Error2.Cut()
*/
