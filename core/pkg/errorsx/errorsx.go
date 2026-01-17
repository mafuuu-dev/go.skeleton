package errorsx

import (
	"backend/core/types"
	"encoding/json"
	"net/http"
	"path/filepath"
	"runtime"
)

type StackFrame struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
	IsHuman bool   `json:"is_human,omitempty"`
}

type Error struct {
	stack      []StackFrame
	humanCode  string
	httpStatus int
}

func (e *Error) Error() string {
	if e.humanCode != "" {
		return e.humanCode
	}
	if len(e.stack) > 0 {
		return e.stack[0].Message
	}
	return "unknown error"
}

func (e *Error) Unwrap() error {
	return nil
}

func (e *Error) WithHuman(code types.HumanErrorCode, httpStatus int) error {
	if e == nil {
		return nil
	}

	if e.humanCode == "" {
		e.humanCode = string(code)
		e.httpStatus = httpStatus
	}

	_, file, line, _ := runtime.Caller(1)
	e.stack = append(e.stack, StackFrame{
		File:    filepath.Base(file),
		Line:    line,
		Message: string(code),
		IsHuman: true,
	})

	return e
}

func (e *Error) GetHuman() (string, int) {
	if e == nil {
		return "", http.StatusInternalServerError
	}
	return e.humanCode, e.httpStatus
}

func (e *Error) GetFirstMessage() string {
	if e == nil || len(e.stack) == 0 {
		return ""
	}
	return e.stack[0].Message
}

func (e *Error) ToJSON() string {
	if e == nil {
		return "{}"
	}

	data := map[string]interface{}{
		"human_code":  e.humanCode,
		"http_status": e.httpStatus,
		"root_cause":  e.GetFirstMessage(),
		"stack":       e.stack,
	}

	b, _ := json.MarshalIndent(data, "", "  ")
	return string(b)
}
