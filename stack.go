package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	// errors_pb "github.com/scottshotgg/errors/proto"
)

// var wd string

// func init() {
// 	fmt.Println("runtime.GOROOT():", runtime.GOROOT())

// 	var workDir, err = os.Getwd()
// 	if err != nil {
// 		panic(fmt.Sprintln("could not get working dir:", err))
// 	}

// 	fmt.Println("workDir:", workDir)

// 	wd = workDir
// }

type (
	Stack []Frame

	Frame struct {
		pc uintptr
		e  *Entry
	}

	Entry struct {
		Name string
		Line int
	}
)

func (e Entry) String() string {
	return e.Name + ":" + strconv.Itoa(e.Line)
	// // return fmt.Sprintf("%s:%s:%d",
	// return fmt.Sprintf("%s:%d",
	// 	// e.File,
	// 	// strings.Split(e.Name, ".")[1],
	// 	e.Name,
	// 	e.Line,
	// )
}

func (s Stack) String() string {
	return strings.Join(s.Strings(), ", ")
}

func (s Stack) Strings() []string {
	var lines = make([]string, len(s))

	for i := range s {
		lines[i] = s[i].String()
	}

	return lines
}

func (f Frame) String() string {
	if f.e == nil {
		f.Resolve()
	}

	return f.e.String()
}

func (s Stack) Frames() {
	var pcs = make([]uintptr, len(s))
	for i := range s {
		pcs[i] = s[i].pc
	}

	var r = runtime.CallersFrames(pcs[:])

	var i int
	for {
		var frame, more = r.Next()
		if !more {
			break
		}

		fmt.Printf("frame: %#v\n", frame)
		i++

		if i > 10 {
			panic("wtf")
		}
	}
}

func (f *Frame) Resolve() {
	var rf = runtime.FuncForPC(f.pc)
	if rf == nil {
		f.e = &Entry{
			Name: "runtime.FuncForPC returned nil",
			Line: -1,
		}

		return
	}

	var _, line = rf.FileLine(f.pc)

	f.e = &Entry{
		// File: strings.Replace(file, wd+"/", "", 1),
		Line: line,
		Name: filepath.Base(rf.Name()),
	}
}

// // TODO: use runtime.Frame here instead
// func (s Stack) Entries() Entries {
// 	if s == nil {
// 		return nil
// 	}

// 	var (
// 		frames  = []errors.Frame(s)
// 		entries = make([]Entry, len(frames))
// 		wd, err = os.Getwd()
// 	)

// 	if err != nil {
// 		return nil
// 	}

// 	for i := range frames {
// 		var (
// 			rf         = runtime.FuncForPC(frames[i])
// 			file, line = rf.FileLine(frames[i])
// 		)

// 		entries[i] = Entry{
// 			File: strings.Replace(file, wd+"/", "", 1),
// 			Line: line,
// 			Name: filepath.Base(rf.Name()),
// 		}
// 	}

// 	return entries
// }

func (s Stack) Cut(dir CutDirection) Stack {
	return s.CutAt(3, dir)
}

func (s Stack) CutAt(skip int, dir CutDirection) Stack {
	if s == nil || dir > Down {
		return nil
	}

	var pc, _, _, ok = runtime.Caller(skip)
	if !ok {
		return s
	}

	for i := range s {
		if s[i].pc == pc {
			switch dir {
			case Up:
				return s[i-1:]

			case Down:
				return s[:i]
			}

			break
		}
	}

	return s
}

// func StackFromPB(frs []*errors_pb.Frame) Stack {
// 	var frames = make([]Frame, len(frs))

// 	for i := range frs {
// 		frames[i] = FrameFromPB(frs[i])
// 	}

// 	return frames
// }

// func FrameFromPB(fr *errors_pb.Frame) Frame {
// 	return Frame{
// 		e: &Entry{
// 			Name: fr.Name,
// 			Line: int(fr.Line),
// 		},
// 	}
// }

// func (s Stack) ToPB() []*errors_pb.Frame {
// 	var frames = make([]*errors_pb.Frame, len(s))

// 	for i := range s {
// 		frames[i] = s[i].ToPB()
// 	}

// 	return frames
// }

// func (f Frame) ToPB() *errors_pb.Frame {
// 	if f.e == nil {
// 		f.Resolve()
// 	}

// 	return &errors_pb.Frame{
// 		Name: f.e.Name,
// 		Line: int32(f.e.Line),
// 	}
// }

// TODO: these should probably be built at some point be we don't
// need them right now

// func StackFromStrings(s []string) Stack {
// 	return nil
// }

// func FrameFromString(s string) Frame {
// 	return Frame{
// 		e: &Entry{},
// 	}
// }
