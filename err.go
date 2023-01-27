package errors

import (
	"reflect"
	"runtime"

	"github.com/pkg/errors"
)

const (
	defaultStackDepth uint8 = 32
)

// Either this of a list of errors instead of frames

// TODO: possibly try using runtime.Frame instead

var _ Error = (*Err)(nil)

type Err struct {
	err error
	// TODO: might want to make the Cause a non-pointer so that if people do things like err.Cause().Error() it doesn't panic
	cause  *Cause
	frames Stack
	// stackDepth uint8
}

func New(err error) *Err {
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

func NewBecause(err error, causeFn interface{}, causeErr error) *Err {
	return New(err).
		Because(causeFn, causeErr)
}

func (e *Err) Because(fn interface{}, err error) *Err {
	if e == nil {
		return nil
	}

	var ptr = runtime.FuncForPC(reflect.ValueOf(fn).Pointer())

	e.frames = append([]errors.Frame{errors.Frame(ptr.Entry() + 1)}, e.frames...)
	e.cause = &Cause{
		name:  ptr.Name(),
		err:   err,
		value: ptr,
	}

	return e
}

// func (e *Err) Cut() *Err {
// 	if e == nil {
// 		return nil
// 	}

// 	e.frames = e.Stack().Cut()

// 	return e
// }

func (e *Err) Cut() *Err {
	if e == nil {
		return nil
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	for i := range e.frames {
		if uintptr(e.frames[i]) == pcs[0] {
			e.frames = e.frames[:i]
			break
		}
	}

	return e
}

func (e *Err) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return e.err.Error()
}

func (e *Err) IsNil() bool {
	return e == nil
}

func (e *Err) Is(err error) bool {
	if err == nil && e == nil {
		return true
	}

	var target, ok = err.(*Err)
	if !ok {
		return false
	}

	return target.err == e.err && target.cause == e.cause
}

func (e *Err) As(target any) bool {
	// TODO: expand on this later
	return errors.As(e.err, target)
}

func (e *Err) Err() error {
	if e.err == nil {
		return nil
	}

	return e.err
}

func (e *Err) Cause() *Cause {
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

func (e *Err) StackTrace() errors.StackTrace {
	if e == nil {
		return nil
	}

	return errors.StackTrace(e.frames)
}

func (e *Err) Trace() *Err {
	if e == nil {
		return nil
	}

	var (
		pcs [1]uintptr
		n   = runtime.Callers(3, pcs[:])
	)

	if n > 0 {
		e.frames = append(e.frames, errors.Frame(pcs[0]))
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
	var frames = make([]errors.Frame, len(pcs))

	for i := range pcs {
		frames[i] = errors.Frame(pcs[i])
	}

	return Stack(frames)
}

// func (g Err) Verbose() string {
// 	return fmt.Sprintf("%s%+v", e.err, e.StackTrace(true))
// }

// func (e *Err) StackTrace(short bool) errors.StackTrace {
// 	var frames []errors.Frame
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
