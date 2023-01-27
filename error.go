package errors

import (
	"reflect"
	"runtime"

	"github.com/pkg/errors"
)

const (
	defaultStackDepth = 32
)

// Either this of a list of errors instead of frames

// TODO: possibly try using runtime.Frame instead

type Error struct {
	err    error
	cause  *Cause
	frames Stack
}

func New(err error) *Error {
	return &Error{
		err:    err,
		frames: callers(),
	}
}

func NewBecause(err error, causeFn interface{}, causeErr error) *Error {
	return New(err).
		Because(causeFn, causeErr)
}

func (g *Error) Because(fn interface{}, err error) *Error {
	if g == nil {
		return nil
	}

	var ptr = runtime.FuncForPC(reflect.ValueOf(fn).Pointer())

	g.frames = append([]errors.Frame{errors.Frame(ptr.Entry() + 1)}, g.frames...)
	g.cause = &Cause{
		name:  ptr.Name(),
		err:   err,
		value: ptr,
	}

	return g
}

// func (e *Error) Cut() *Error {
// 	if e == nil {
// 		return nil
// 	}

// 	e.frames = e.Stack().Cut()

// 	return e
// }

func (g *Error) Cut() *Error {
	if g == nil {
		return nil
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	for i := range g.frames {
		if uintptr(g.frames[i]) == pcs[0] {
			g.frames = g.frames[:i]
			break
		}
	}

	return g
}

func (g *Error) Error() string {
	if g == nil || g.err == nil {
		return ""
	}

	return g.err.Error()
}

func (g *Error) IsNil() bool {
	return g == nil
}

func (g *Error) Is(target error) bool {
	if target == nil && g == nil {
		return true
	}

	var tgt, ok = target.(*Error)
	if !ok {
		return false
	}

	return tgt.err == g.err && tgt.cause == g.cause
}

func (g *Error) Err() error {
	if g.err == nil {
		return nil
	}

	return g.err
}

func (g *Error) Cause() *Cause {
	if g == nil {
		return nil
	}

	return g.cause
}

func (g *Error) Stack() Stack {
	if g == nil {
		return nil
	}

	return g.frames
}

func (g *Error) StackTrace() errors.StackTrace {
	if g == nil {
		return nil
	}

	return errors.StackTrace(g.frames)
}

func (g *Error) Trace() *Error {
	if g == nil {
		return nil
	}

	var (
		pcs [1]uintptr
		n   = runtime.Callers(3, pcs[:])
	)

	if n > 0 {
		g.frames = append(g.frames, errors.Frame(pcs[0]))
	}

	return g
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

// func (g Error) Verbose() string {
// 	return fmt.Sprintf("%s%+v", g.err, g.StackTrace(true))
// }

// func (g *Error) StackTrace(short bool) errors.StackTrace {
// 	var frames []errors.Frame
// 	if short {
// 		for _, frame := range g.frames {
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

// 	return g.frames
// }

// func (g *Error) StackString() string {
// 	var lines = make([]string, len(g.frames))

// 	for i, frame := range g.frames {
// 		var rf = runtime.FuncForPC(uintptr(frame))

// 		file, _ := rf.FileLine(uintptr(frame))

// 		var nameSplit = strings.Split(rf.Name(), "/")

// 		lines[i] = fmt.Sprintf("%s:%s", filepath.Base(file), nameSplit[len(nameSplit)-1])
// 	}

// 	return fmt.Sprintf("[ %s ]", strings.Join(lines, ", "))
// }
