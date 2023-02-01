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

	var se = errors.FromError(err)
	fmt.Println("se:", se)
	// fmt.Println(se.Cut(errors.Up).Stack())
	// fmt.Println(se.Cut(errors.Down).Stack())

	// fmt.Println(err.Cut().StackTrace())
}

func someFunc() errors.Error {
	var err1 = someOtherFunc1()
	fmt.Println("up:", err1.Cut(errors.Up).Stack())
	var err2 = someOtherFunc1()
	fmt.Println("down:", err2.Cut(errors.Down).Stack())

	return someOtherFunc1()
}

func someOtherFunc1() errors.Error {
	return errors.New(stderrors.New("abc what the func"))
}

func someOtherFunc2() errors.Error {
	return errors.New(stderrors.New("123 what the func"))
}

/*
	Error could wrap a multierror of Errors

	Error.New().Cut():
		multierror:
			Error1.Cut()
			Error2.Cut()
*/