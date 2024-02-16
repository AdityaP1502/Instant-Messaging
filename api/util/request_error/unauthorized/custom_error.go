package unauthorized

import (
	"fmt"

	requesterror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error"
)

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

func (e *InvalidAuthHeaderError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Name:    e.Name,
		Code:    e.Code,
		Message: e.Message,
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
		Message: "Required authorization header in request header",
		Code:    403,
	}
}

func (e *EmptyAuthHeaderError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Name:    e.Name,
		Code:    e.Code,
		Message: e.Message,
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
		Name:    string(InvalidToken),
		Message: fmt.Sprintf("Invalid token.%s", f),
		Code:    403,
	}
}

func (e *InvalidTokenError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Name:    e.Name,
		Code:    e.Code,
		Message: e.Message,
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
		Name:    string(TokenExpired),
		Message: "Your token has expired.",
		Code:    403,
	}
}

func (e *TokenExpiredError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Name:    e.Name,
		Code:    e.Code,
		Message: e.Message,
	}
}

var TokenExpiredErr = &TokenExpiredError{}

type RefreshDeniedError struct {
	Name    string
	Code    int
	Message string
}

func (e *RefreshDeniedError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *RefreshDeniedError) Init() error {
	return &RefreshDeniedError{
		Name:    string(RefreshDenied),
		Message: "Cannot get new access token when the previous one still active",
		Code:    403,
	}
}

func (e *RefreshDeniedError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Name:    e.Name,
		Code:    e.Code,
		Message: e.Message,
	}
}

var RefreshDeniedErr = &RefreshDeniedError{}

type ClaimsMismatchError struct {
	Name    string
	Code    int
	Message string
}

func (e *ClaimsMismatchError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *ClaimsMismatchError) Init() error {
	return &ClaimsMismatchError{
		Name:    string(ClaimsMismatch),
		Message: "Refresh claims and username claims don't share the same credentials",
		Code:    403,
	}
}

func (e *ClaimsMismatchError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Name:    e.Name,
		Code:    e.Code,
		Message: e.Message,
	}
}

var ClaimsMismatchErr = &ClaimsMismatchError{}
