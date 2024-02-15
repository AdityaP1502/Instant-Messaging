package internalserviceerror

import (
	"fmt"

	requesterror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error"
)

type InternalServiceError struct {
	Name        string
	Code        int
	Message     string
	Description string
}

func (e *InternalServiceError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *InternalServiceError) Init(f string) error {
	return &InternalServiceError{
		Name:        "internal_service_error",
		Message:     "Sorry, it seems there are some problems on your request. Please try again!",
		Code:        500,
		Description: f,
	}
}

func (e *InternalServiceError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

var InternalServiceErr *InternalServiceError = &InternalServiceError{}
