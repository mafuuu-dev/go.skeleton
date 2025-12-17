package errorsx_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"backend/core/pkg/errorsx"
)

func TestJSONTrace_WithTraceError(t *testing.T) {
	orig := errors.New("db connection failed")
	err := errorsx.Error(orig)

	jsonStr := errorsx.JSONTrace(err)

	if !strings.Contains(jsonStr, `"msg": "db connection failed"`) {
		t.Errorf("expected message in JSONTrace output, got:\n%s", jsonStr)
	}

	if !strings.Contains(jsonStr, `"file":`) {
		t.Errorf("expected file field in JSON output")
	}
}

func TestJSONTrace_WithPlainError(t *testing.T) {
	err := errors.New("plain error")
	jsonStr := errorsx.JSONTrace(err)

	if !strings.Contains(jsonStr, `"msg": "plain error"`) {
		t.Errorf("expected message 'plain error' in JSONTrace output, got:\n%s", jsonStr)
	}

	if !strings.Contains(jsonStr, `"file": "formatx.go"`) {
		t.Errorf("expected file 'formatx.go' in JSONTrace output, got:\n%s", jsonStr)
	}

	if !strings.Contains(jsonStr, `"line":`) {
		t.Errorf("expected line number in JSONTrace output, got:\n%s", jsonStr)
	}

	if !strings.Contains(jsonStr, `"trace": null`) {
		t.Errorf("expected 'trace': null in JSONTrace output, got:\n%s", jsonStr)
	}
}

func TestEnrichTrace_AddsFields(t *testing.T) {
	traceErr := &errorsx.TraceError{
		Root: errorsx.Frame{
			File: "main.go",
			Line: 42,
			Msg:  "initial error",
		},
		Trace: []errorsx.Frame{{
			File: "db.go",
			Line: 77,
			Msg:  "query failed",
		}},
	}
	data, _ := json.Marshal(traceErr)
	trace := string(data)

	enriched := errorsx.EnrichTrace(trace, http.StatusBadRequest, "bad request")

	if !strings.Contains(enriched, `"code": 400`) {
		t.Errorf("expected 'code' in enriched trace, got:\n%s", enriched)
	}
	if !strings.Contains(enriched, `"message": "bad request"`) {
		t.Errorf("expected 'message' in enriched trace, got:\n%s", enriched)
	}
	if !strings.Contains(enriched, `"initial error"`) {
		t.Errorf("original trace lost:\n%s", enriched)
	}
}
