package responseerror

import (
	"fmt"
)

const (
	UserMarkedInActive errorType = "user_marked_inactive"
	InvalidCredentials errorType = "invalid_credentials"
)

type InactiveUserError struct {
	Name    string
	Code    int
	Message string
}

func (e *InactiveUserError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

// Init the value of empty InactiveUserError
//
// args: f (string): the missing field
func (e *InactiveUserError) Init() error {
	return &InactiveUserError{
		Name:    string(UserMarkedInActive),
		Message: "The user isn't registered successfully and currently marked as inavtive",
		Code:    401,
	}
}

func (e *InactiveUserError) Get() *ResponseError {
	return &ResponseError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

type InvalidCredentialsError struct {
	Name    string
	Code    int
	Message string
}

func (e *InvalidCredentialsError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

// Init the value of empty InvalidCredentialsError
//
// args: f (string): the missing field
func (e *InvalidCredentialsError) Init() error {
	return &InvalidCredentialsError{
		Name:    string(InvalidCredentials),
		Message: "Email or password are invalid",
		Code:    401,
	}
}

func (e *InvalidCredentialsError) Get() *ResponseError {
	return &ResponseError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

var InactiveUserErr = &InactiveUserError{}

var InvalidCredentialsErr = &InvalidCredentialsError{}
