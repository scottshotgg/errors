package errors

import (
	"errors"
	"reflect"
	"runtime"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

const (
	// TODO: fix this default
	defaultStackDepth uint8 = 16
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
	cause Cause
	// causes []Cause
	stack Stack
	// stackDepth uint8

	details map[string]any
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
		err:     err,
		stack:   callers(),
		details: map[string]any{},
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

	var rv = reflect.ValueOf(fn)
	if rv.Kind() != reflect.Func {
		return &Err{
			err: errors.New("causeFn must be 'reflect.Func'"),
		}
	}

	var ptr = runtime.FuncForPC(rv.Pointer())
	// var entry = ptr.Entry()
	// _, line := ptr.FileLine(entry)

	e.stack = append([]Frame{
		{
			pc: ptr.Entry(),
			// pc: entry,
			// e: &Entry{
			// 	Name: ptr.Name(),
			// 	Line: line,
			// },
		},
	}, e.stack...)

	var split = strings.Split(ptr.Name(), "/")

	e.cause = NewReason(err, split[len(split)-1])
	// e.causes = []Cause{
	// NewReason(err, ptr),
	// }

	return e
}

// func (e *Err) Cut() *Err {
// 	if e == nil {
// 		return nil
// 	}

// 	e.stack = e.Stack().Cut()

// 	return e
// }

func (e *Err) Cut(dir CutDirection) Error {
	e.stack = e.stack.CutAt(3, dir)

	return e
}

func (e *Err) Details() map[string]any {
	return e.details
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
	if err == nil && e == nil {
		return true
	}

	var _, ok = err.(*Err)
	if !ok {
		return false
	}

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

// func (e *Err) Causes() []Cause {
// 	if e == nil {
// 		return nil
// 	}

// 	return e.causes
// }

func (e *Err) Stack() Stack {
	if e == nil {
		return nil
	}

	return e.stack
}

// func (e *Err) StackTrace() pkgerrors.StackTrace {
// 	if e == nil {
// 		return nil
// 	}

// 	return pkgerrors.StackTrace(e.stack)
// }

// func (e *Err) Trace() Error {
// 	if e == nil {
// 		return nil
// 	}

// 	var (
// 		pc, _, _, ok = runtime.Caller(2)
// 	)

// 	if !ok {
// 		return e
// 	}

// 	e.stack = append(e.stack, pkgerrors.Frame(pc))

// 	return e
// }

func callers() Stack {
	var (
		pcs [defaultStackDepth]uintptr
		i   int

		// TODO: not sure if this should actually be n, need to look into it
		n      = runtime.Callers(3, pcs[:])
		frames = make([]Frame, n-1)
		r      = runtime.CallersFrames(pcs[:n])
	)

	for {
		var frame, more = r.Next()
		if !more {
			break
		}

		// if frame.Func == nil {
		// 	frames = append(frames, Frame{
		// 		pc: frame.PC,
		// 	})

		// 	continue
		// }

		frames[i] = Frame{
			pc: frame.PC,
			// TODO: unfortunately - this Function is too long and Func is nil for inlined functions ...
			// Figure this out latersS
			// e: &Entry{
			// 	Name: frame.Function,
			// 	Line: frame.Line,
			// },
		}

		i++
	}

	return frames
}

// func (g Err) Verbose() string {
// 	return fmt.Sprintf("%s%+v", e.err, e.StackTrace(true))
// }

// func (e *Err) StackTrace(short bool) pkgerrors.StackTrace {
// 	var frames []pkgerrors.Frame
// 	if short {
// 		for _, frame := range e.stack {
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

// 	return e.stack
// }

// func (e *Err) StackString() string {
// 	var lines = make([]string, len(e.stack))

// 	for i, frame := range e.stack {
// 		var rf = runtime.FuncForPC(uintptr(frame))

// 		file, _ := rf.FileLine(uintptr(frame))

// 		var nameSplit = strings.Split(rf.Name(), "/")

// 		lines[i] = fmt.Sprintf("%s:%s", filepath.Base(file), nameSplit[len(nameSplit)-1])
// 	}

// 	return fmt.Sprintf("[ %s ]", strings.Join(lines, ", "))
// }
