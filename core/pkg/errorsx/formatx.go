package errorsx

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"runtime"
)

func JSONTrace(err error) string {
	var te *TraceError
	if errors.As(err, &te) {
		b, _ := json.MarshalIndent(te, "", "  ")
		return string(b)
	}

	var he *HumanError
	if errors.As(err, &he) {
		return JSONTrace(he.Unwrap())
	}

	_, file, line, _ := runtime.Caller(0)
	te = &TraceError{
		Root: Frame{
			File: filepath.Base(file),
			Line: line,
			Msg:  err.Error(),
		},
	}

	return JSONTrace(te)
}

func EnrichTrace(trace string, code int, message string) string {
	var m map[string]interface{}
	_ = json.Unmarshal([]byte(trace), &m)

	m["code"] = code
	m["message"] = message

	enrichmentTrace, _ := json.MarshalIndent(m, "", "  ")

	return string(enrichmentTrace)
}
