package errors

import (
	"fmt"
	"reflect"
	"runtime"

	pkgerrors "github.com/pkg/errors"
)

const (
	// TODO: fix this default
	defaultStackDepth uint8 = 32
)

// Either this of a list of errors instead of frames

// TODO: possibly try using runtime.Frame instead

// Err implements Error
var _ Error = (*Err)(nil)

type Err struct {
	err error
	// TODO: might want to make the Cause a non-pointer so that if people do things like err.Cause().Error() it doesn't panic
	// TODO: might want to change the verbase here to Reasons, make it an array
	// that way you can just tag things along the way that you want to keep
	cause  Cause
	frames Stack
	// stackDepth uint8
}

func New(err error) *Err {
	if err == nil {
		return nil
	}

	// TODO: return what we already have?
	// var e, ok = err.(*Err)
	// if ok {
	// 	return e
	// }

	return &Err{
		err:    err,
		frames: callers(),
		// stackDepth: defaultStackDepth,
	}
}

// func (e *Err) WithStackDepth(sd uint8) *Err {
// 	e.stackDepth = sd

// 	return e
// }

func NewBecause(err error, causeFn interface{}, causeErr error) Error {
	return New(err).
		Because(causeFn, causeErr)
}

func (e *Err) Because(fn interface{}, err error) Error {
	if e == nil {
		return nil
	}

	var ptr = runtime.FuncForPC(reflect.ValueOf(fn).Pointer())

	e.frames = append([]pkgerrors.Frame{pkgerrors.Frame(ptr.Entry() + 1)}, e.frames...)
	e.cause = NewReason(ptr.Name(), err, ptr)

	return e
}

// func (e *Err) Cut() *Err {
// 	if e == nil {
// 		return nil
// 	}

// 	e.frames = e.Stack().Cut()

// 	return e
// }

func (e *Err) Cut(dir CutDirection) Error {
	if e == nil || dir > Down {
		return nil
	}

	var pc, _, _, ok = runtime.Caller(2)
	if !ok {
		return e
	}

	for i := range e.frames {
		if uintptr(e.frames[i]) == pc+1 {
			switch dir {
			case Up:
				e.frames = e.frames[i-1:]

			case Down:
				e.frames = e.frames[:i]
			}

			break
		}
	}

	return e
}

func (e *Err) String() string {
	// TODO: make this more verbose
	return "TODO: make this more verbose"
}

func (e *Err) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	// TODO: expand this, make it more verbose
	return e.err.Error()
}

func (e *Err) IsNil() bool {
	return e == nil
}

// TODO: think about this functionality
func (e *Err) Is(err error) bool {
	fmt.Println("running is func")
	if err == nil && e == nil {
		fmt.Println("both nil?")
		return true
	}

	var _, ok = err.(*Err)
	if !ok {
		fmt.Println("not ok?")
		return false
	}

	fmt.Println("compare?")

	return true
	// return target.err == e.err && target.cause == e.cause
}

func (e *Err) As(target any) bool {
	// TODO: expand on this later
	return pkgerrors.As(e.err, target)
}

func (e *Err) Err() error {
	if e.err == nil {
		return nil
	}

	return e.err
}

func (e *Err) Cause() Cause {
	if e == nil {
		return nil
	}

	return e.cause
}

func (e *Err) Stack() Stack {
	if e == nil {
		return nil
	}

	return e.frames
}

func (e *Err) StackTrace() pkgerrors.StackTrace {
	if e == nil {
		return nil
	}

	return pkgerrors.StackTrace(e.frames)
}

func (e *Err) Trace() Error {
	if e == nil {
		return nil
	}

	var (
		pcs [1]uintptr
		n   = runtime.Callers(3, pcs[:])
	)

	if n > 0 {
		e.frames = append(e.frames, pkgerrors.Frame(pcs[0]))
	}

	return e
}

func callers() Stack {
	var (
		pcs [defaultStackDepth]uintptr
		n   = runtime.Callers(3, pcs[:])
	)

	return toStack(pcs[:n])
}

func toStack(pcs []uintptr) Stack {
	var frames = make([]pkgerrors.Frame, len(pcs))

	for i := range pcs {
		frames[i] = pkgerrors.Frame(pcs[i])
	}

	return Stack(frames)
}

// func (g Err) Verbose() string {
// 	return fmt.Sprintf("%s%+v", e.err, e.StackTrace(true))
// }

// func (e *Err) StackTrace(short bool) pkgerrors.StackTrace {
// 	var frames []pkgerrors.Frame
// 	if short {
// 		for _, frame := range e.frames {
// 			var rf = runtime.FuncForPC(uintptr(frame))

// 			file, line := rf.FileLine(uintptr(frame))

// 			if strings.Contains(file, "paynearme") {
// 				frames = append(frames, frame)
// 				fmt.Printf("name: %+v\n", rf.Name())
// 				fmt.Printf("entry: %+v\n", rf.Entry())
// 				fmt.Printf("line: %+v\n", file)
// 				fmt.Printf("line: %+v\n", line)
// 				fmt.Println()
// 			}
// 		}
// 	}

// 	return e.frames
// }

// func (e *Err) StackString() string {
// 	var lines = make([]string, len(e.frames))

// 	for i, frame := range e.frames {
// 		var rf = runtime.FuncForPC(uintptr(frame))

// 		file, _ := rf.FileLine(uintptr(frame))

// 		var nameSplit = strings.Split(rf.Name(), "/")

// 		lines[i] = fmt.Sprintf("%s:%s", filepath.Base(file), nameSplit[len(nameSplit)-1])
// 	}

// 	return fmt.Sprintf("[ %s ]", strings.Join(lines, ", "))
// }
