package toomanyrequest

import (
	"fmt"

	requesterror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error"
)

type ResendIntervalNotReachedError struct {
	Name    string
	Code    int
	Message string
}

func (e *ResendIntervalNotReachedError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

// Init the value of empty ResendIntervalNotReachedError
//
// args: f (string): the missing field
func (e *ResendIntervalNotReachedError) Init() error {
	return &ResendIntervalNotReachedError{
		Name:    "otp_resend_interval_not_reached",
		Message: "Mail has already been sent",
		Code:    429,
	}
}

func (e *ResendIntervalNotReachedError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

var ResendIntervalNotReachedErr = &ResendIntervalNotReachedError{}
