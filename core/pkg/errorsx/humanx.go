package errorsx

import "backend/core/types"

type HumanError struct {
	Message types.HumanErrorCode
	Status  int
	Err     error
}

func (h *HumanError) Error() string {
	return string(h.Message)
}

func (h *HumanError) Unwrap() error {
	return h.Err
}

func Humanf(err error, message types.HumanErrorCode, status int) error {
	return &HumanError{
		Message: message,
		Status:  status,
		Err:     Error(err),
	}
}
