package errorsx

import (
	"backend/core/types"
	"errors"
)

func WrapHuman(err error, humanCode types.HumanErrorCode, status int) error {
	if err == nil {
		return nil
	}

	var e *Error
	errors.As(Wrap(err, string(humanCode)), &e)
	return e.WithHuman(humanCode, status)
}

func WrapJSON(err error, message string) string {
	if err == nil {
		return ""
	}

	return Extract(Wrap(err, message)).(*Error).ToJSON()
}

func WrapfJSON(err error, format string, args ...interface{}) string {
	if err == nil {
		return ""
	}

	return Extract(Wrapf(err, format, args...)).(*Error).ToJSON()
}
