package errorsx

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Frame struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Msg  string `json:"msg"`
}

type TraceError struct {
	Root  Frame   `json:"root"`
	Trace []Frame `json:"trace"`
}

func (e *TraceError) Error() string {
	return e.Root.Msg
}

func Error(err error) error {
	if err == nil {
		return nil
	}

	_, file, line, _ := runtime.Caller(1)

	var te *TraceError
	if errors.As(err, &te) {
		te.Trace = append(te.Trace, Frame{
			File: filepath.Base(file),
			Line: line,
			Msg:  te.Root.Msg,
		})

		return te
	}

	return &TraceError{
		Root: Frame{
			File: filepath.Base(file),
			Line: line,
			Msg:  err.Error(),
		},
		Trace: []Frame{{
			File: filepath.Base(file),
			Line: line,
			Msg:  err.Error(),
		}},
	}
}

func Errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	_, file, line, _ := runtime.Caller(1)

	return &TraceError{
		Root: Frame{
			File: filepath.Base(file),
			Line: line,
			Msg:  msg,
		},
		Trace: []Frame{{
			File: filepath.Base(file),
			Line: line,
			Msg:  msg,
		}},
	}
}

func Panic(r any) *TraceError {
	msg := fmt.Sprintf("%v", r)
	stack := debug.Stack()

	lines := bytes.Split(stack, []byte("\n"))
	trace := make([]Frame, 0, len(lines)/2)

	for i := 0; i < len(lines)-1; i++ {
		line := string(lines[i])
		if !strings.HasPrefix(line, "\t") {
			continue
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		file := filepath.Base(parts[0])
		lineNumStr := strings.Split(parts[1], " ")[0]
		lineNum, _ := strconv.Atoi(lineNumStr)

		trace = append(trace, Frame{
			File: file,
			Line: lineNum,
			Msg:  msg,
		})
	}

	root := Frame{}
	if len(trace) > 0 {
		root = trace[0]
	}

	return &TraceError{
		Root:  root,
		Trace: trace,
	}
}
