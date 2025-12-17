package errorsx_test

import (
	"errors"
	"strings"
	"testing"

	"backend/core/pkg/errorsx"
)

func TestError_Basic(t *testing.T) {
	origErr := errors.New("original failure")
	err := errorsx.Error(origErr)

	if err == nil {
		t.Fatal("expected non-nil error")
	}

	var te *errorsx.TraceError
	ok := errors.As(err, &te)
	if !ok {
		t.Fatalf("expected TraceError, got %T", err)
	}

	if te.Root.Msg != "original failure" {
		t.Errorf("expected message 'original failure', got '%s'", te.Root.Msg)
	}

	if len(te.Trace) != 1 {
		t.Errorf("expected 1 trace frame, got %d", len(te.Trace))
	}
}

func TestError_Wrapping(t *testing.T) {
	base := errors.New("root cause")
	err1 := errorsx.Error(base)
	err2 := errorsx.Error(err1)

	var te *errorsx.TraceError
	ok := errors.As(err2, &te)
	if !ok {
		t.Fatalf("expected TraceError, got %T", err2)
	}

	if len(te.Trace) != 2 {
		t.Errorf("expected 2 frames after wrapping, got %d", len(te.Trace))
	}

	if te.Root.Msg != "root cause" {
		t.Errorf("expected root msg 'root cause', got '%s'", te.Root.Msg)
	}
}

func TestErrorf(t *testing.T) {
	err := errorsx.Errorf("failed to load %s #%d", "user", 42)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var te *errorsx.TraceError
	ok := errors.As(err, &te)
	if !ok {
		t.Fatalf("expected TraceError, got %T", err)
	}

	if !strings.Contains(te.Root.Msg, "failed to load user #42") {
		t.Errorf("unexpected message: %s", te.Root.Msg)
	}

	if len(te.Trace) != 1 {
		t.Errorf("expected 1 frame, got %d", len(te.Trace))
	}
}

func TestPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		traceErr := errorsx.Panic(r)

		if traceErr == nil {
			t.Fatal("expected non-nil TraceError from Panic()")
		}
		if traceErr.Root.Msg != "simulated panic" {
			t.Errorf("expected message 'simulated panic', got '%s'", traceErr.Root.Msg)
		}
		if len(traceErr.Trace) == 0 {
			t.Error("expected non-empty trace from Panic()")
		}

		found := false
		for _, frame := range traceErr.Trace {
			if strings.HasSuffix(frame.File, "_test.go") {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected at least one frame from test file in stack trace")
		}
	}()

	panic("simulated panic")
}
