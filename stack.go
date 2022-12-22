package errors

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type (
	Stack []errors.Frame

	Entry struct {
		File string
		Line int
		Name string
	}

	Entries []Entry
)

func (s Stack) Cut() Stack {
	if s == nil {
		return nil
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	for i := range s {
		if uintptr(s[i]) == pcs[0] {
			return s[:i]
		}
	}

	return s
}

func (e Entry) String() string {
	return fmt.Sprintf("%s:%s:%d",
		e.File,
		e.Name,
		e.Line,
	)
}

func (es Entries) String() string {
	var lines = make([]string, len(es))

	for i := range es {
		lines[i] = es[i].String()
	}

	return strings.Join(lines, ", ")
}

// TODO: use runtime.Frame here instead
func (s Stack) Entries() Entries {
	var (
		frames  = []errors.Frame(s)
		entries = make([]Entry, len(frames))
		wd, err = os.Getwd()
	)

	if err != nil {
		return nil
	}

	for i := range frames {
		var (
			rf         = runtime.FuncForPC(uintptr(frames[i]))
			file, line = rf.FileLine(uintptr(frames[i]))
		)

		entries[i] = Entry{
			File: strings.Replace(file, wd+"/", "", 1),
			Line: line,
			Name: filepath.Base(rf.Name()),
		}
	}

	return entries
}
