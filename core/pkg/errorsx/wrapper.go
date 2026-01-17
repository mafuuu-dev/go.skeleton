package errorsx

import (
	"backend/core/types"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
)

func New(message string) error {
	_, file, line, _ := runtime.Caller(1)

	return &Error{
		stack: []StackFrame{{
			File:    filepath.Base(file),
			Line:    line,
			Message: message,
		}},
		httpStatus: http.StatusInternalServerError,
	}
}

func Newf(format string, args ...interface{}) error {
	_, file, line, _ := runtime.Caller(1)

	return &Error{
		stack: []StackFrame{{
			File:    filepath.Base(file),
			Line:    line,
			Message: fmt.Sprintf(format, args...),
		}},
		httpStatus: http.StatusInternalServerError,
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	_, file, line, _ := runtime.Caller(1)

	var e *Error
	if errors.As(err, &e) {
		e.stack = append(e.stack, StackFrame{
			File:    filepath.Base(file),
			Line:    line,
			Message: message,
		})
		return e
	}

	return &Error{
		stack: []StackFrame{
			{
				File:    filepath.Base(file),
				Line:    line,
				Message: err.Error(),
			},
			{
				File:    filepath.Base(file),
				Line:    line,
				Message: message,
			},
		},
		httpStatus: http.StatusInternalServerError,
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	return Wrap(err, fmt.Sprintf(format, args...))
}

func Human(code types.HumanErrorCode, httpStatus int) error {
	_, file, line, _ := runtime.Caller(1)

	return &Error{
		stack: []StackFrame{{
			File:    filepath.Base(file),
			Line:    line,
			Message: string(code),
			IsHuman: true,
		}},
		humanCode:  string(code),
		httpStatus: httpStatus,
	}
}

func Extract(err error) error {
	if err == nil {
		return nil
	}

	var e *Error
	if errors.As(err, &e) {
		return e
	}

	_, file, line, _ := runtime.Caller(1)
	return &Error{
		stack: []StackFrame{{
			File:    filepath.Base(file),
			Line:    line,
			Message: err.Error(),
		}},
		httpStatus: http.StatusInternalServerError,
	}
}

func Recover(r interface{}) error {
	if r == nil {
		return nil
	}

	err, ok := r.(error)
	if !ok {
		return FromPanic(fmt.Errorf("%v", r))
	}

	var e *Error
	if errors.As(err, &e) {
		return e
	}

	return FromPanic(err)
}

func FromPanic(err error) error {
	if err == nil {
		return nil
	}

	const maxStackSize = 32
	pcs := make([]uintptr, maxStackSize)
	n := runtime.Callers(3, pcs)

	frames := runtime.CallersFrames(pcs[:n])
	stack := make([]StackFrame, 0, n)

	firstFrame, more := frames.Next()
	stack = append(stack, StackFrame{
		File:    filepath.Base(firstFrame.File),
		Line:    firstFrame.Line,
		Message: err.Error(),
	})

	for more {
		frame, hasMore := frames.Next()
		more = hasMore

		stack = append(stack, StackFrame{
			File:    filepath.Base(frame.File),
			Line:    frame.Line,
			Message: fmt.Sprintf("at %s", frame.Function),
		})
	}

	return &Error{
		stack:      stack,
		httpStatus: http.StatusInternalServerError,
	}
}
