package errorsx_test

import (
	"errors"
	"net/http"
	"testing"

	"backend/core/pkg/errorsx"
)

func TestHumanf_Basic(t *testing.T) {
	origErr := errors.New("database unreachable")

	hErr := errorsx.Humanf(origErr, "failed to process request", http.StatusBadRequest)
	if hErr == nil {
		t.Fatal("expected non-nil HumanError")
	}

	var he *errorsx.HumanError
	ok := errors.As(hErr, &he)
	if !ok {
		t.Fatalf("expected *HumanError, got %T", hErr)
	}

	if he.Message != "failed to process request" {
		t.Errorf("expected message 'failed to process request', got '%s'", he.Message)
	}
	if he.Status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", he.Status)
	}
	if he.Err == nil {
		t.Errorf("expected internal error to be set")
	}

	if he.Error() != "failed to process request" {
		t.Errorf("Error() returned '%s', expected '%s'", he.Error(), he.Message)
	}
}

func TestHumanf_WrappedTrace(t *testing.T) {
	baseErr := errors.New("read timeout")
	hErr := errorsx.Humanf(baseErr, "request failed", http.StatusInternalServerError)

	var he *errorsx.HumanError
	ok := errors.As(hErr, &he)
	if !ok {
		t.Fatalf("expected *HumanError, got %T", hErr)
	}

	var traceError *errorsx.TraceError
	if !errors.As(he.Err, &traceError) {
		t.Errorf("expected wrapped error to be *TraceError, got %T", he.Err)
	}
}
