package unauthorized

import "fmt"

type InvalidAuthHeaderError struct {
	Name    string
	Code    int
	Message string
}

func (e *InvalidAuthHeaderError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *InvalidAuthHeaderError) Init(f string) error {
	return &InvalidAuthHeaderError{
		Name:    string(InvalidAuthHeader),
		Message: fmt.Sprintf("Not accepted authorization of type '%s'", f),
		Code:    403,
	}
}

var InvalidAuthHeaderErr = &InvalidAuthHeaderError{}

type EmptyAuthHeaderError struct {
	Name    string
	Code    int
	Message string
}

func (e *EmptyAuthHeaderError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *EmptyAuthHeaderError) Init() error {
	return &EmptyAuthHeaderError{
		Name:    string(EmptyAuthHeader),
		Message: fmt.Sprintf("Required authorization header in request header"),
		Code:    403,
	}
}

var EmptyAuthHeaderErr = &EmptyAuthHeaderError{}

type InvalidTokenError struct {
	Name    string
	Code    int
	Message string
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *InvalidTokenError) Init(f string) error {
	return &InvalidTokenError{
		Name:    string(InvalidAuthHeader),
		Message: fmt.Sprintf("Invalid token.%s", f),
		Code:    403,
	}
}

var InvalidTokenErr = &InvalidTokenError{}

type TokenExpiredError struct {
	Name    string
	Code    int
	Message string
}

func (e *TokenExpiredError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *TokenExpiredError) Init() error {
	return &TokenExpiredError{
		Name:    string(InvalidAuthHeader),
		Message: fmt.Sprintf("Your token has expired."),
		Code:    403,
	}
}

var TokenExpiredErr = &TokenExpiredError{}
