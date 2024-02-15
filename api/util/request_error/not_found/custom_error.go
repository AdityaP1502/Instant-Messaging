package notfound

import (
	"fmt"

	requesterror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error"
)

type NotFoundError struct {
	Name    string
	Code    int
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *NotFoundError) Init(name string) error {
	return &NotFoundError{
		Name:    fmt.Sprintf("%s_not_found", name),
		Message: fmt.Sprintf("%s provided doesn't exist", name),
		Code:    404,
	}
}

func (e *NotFoundError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

var InternalServiceErr *NotFoundError = &NotFoundError{}
