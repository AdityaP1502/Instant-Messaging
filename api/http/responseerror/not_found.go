package responseerror

import (
	"fmt"
)

type NotFoundError struct {
	Name    string
	Code    int
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *NotFoundError) Init(id string, name string) error {
	return &NotFoundError{
		Name:    fmt.Sprintf("%s_not_found", id),
		Message: fmt.Sprintf("%s provided doesn't exist", name),
		Code:    404,
	}
}

func (e *NotFoundError) Get() ResponseError {
	return ResponseError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

var NotFoundErr *NotFoundError = &NotFoundError{}
