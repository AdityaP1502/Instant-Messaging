package badrequest

import (
	"fmt"

	requesterror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error"
)

// MissingParametersError
//
// Name: "MissingParameter"
//
// Code: 400
//
// Message: Required field %s is empty, where %s is the field in which the field is empty
type MissingParameterError struct {
	Name    string
	Code    int
	Message string
}

func (e *MissingParameterError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

// Init the value of empty MissingParameterError
//
// args: f (string): the missing field
func (e *MissingParameterError) Init(f string) error {
	return &MissingParameterError{
		Name:    string(MissingParameter),
		Message: fmt.Sprintf("Required field %s is empty", f),
		Code:    400,
	}
}

func (e *MissingParameterError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

// MissingParametersError
//
// Name: "MissingParameter"
//
// Code: 400
//
// Message: Required field %s is empty, where %s is the field in which the field is empty
type HeaderMismatchError struct {
	Name    string
	Code    int
	Message string
}

func (e *HeaderMismatchError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

// Init the value of HeaderMismatchError
//
// args:
//
// h (string): header name where mismatch occured
func (e *HeaderMismatchError) Init(h string) error {
	return &HeaderMismatchError{
		Name:    string(HeaderValueMistmatch),
		Message: fmt.Sprintf("Mismatch value in header %s.", h),
		Code:    400,
	}
}

func (e *HeaderMismatchError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

type ValueNotUniqueError struct {
	Name    string
	Code    int
	Message string
}

func (e *ValueNotUniqueError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

// Create a new ValueNotUniqueError
//
// args:
//
// v (string): value that not unique
//
// t (errorType): either UsernameExits or EmailExits
func (e *ValueNotUniqueError) Init(t errorType, v string) error {
	return &ValueNotUniqueError{
		Name:    string(t),
		Message: fmt.Sprintf("%s already taken", v),
		Code:    400,
	}
}

func (e *ValueNotUniqueError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

type WeakPasswordError struct {
	Name    string
	Code    int
	Message string
}

func (e *WeakPasswordError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *WeakPasswordError) Init() error {
	return &WeakPasswordError{
		Name:    string(PasswordWeak),
		Message: fmt.Sprintf("Password are too weak. Password need to be at minumum of 8 character with combination with letter and symbol"),
		Code:    400,
	}
}

func (e *WeakPasswordError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

type InvalidEmailError struct {
	Name    string
	Code    int
	Message string
}

func (e *InvalidEmailError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *InvalidEmailError) Init() error {
	return &InvalidEmailError{
		Name:    string(EmailInvalid),
		Message: fmt.Sprintf("Email are invalid"),
		Code:    400,
	}
}

func (e *InvalidEmailError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

type UsernameInvalidError struct {
	Name    string
	Code    int
	Message string
}

func (e *UsernameInvalidError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *UsernameInvalidError) Init() error {
	return &UsernameInvalidError{
		Name:    string(PasswordWeak),
		Message: "Username invalid. Username need to be at max 64 characters and don't contain Uppercase characters and invalid characters",
		Code:    400,
	}
}

func (e *UsernameInvalidError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

type InvalidPayloadError struct {
	Name    string
	Code    int
	Message string
}

func (e *InvalidPayloadError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *InvalidPayloadError) Init() error {
	return &InvalidPayloadError{
		Name:    string(PayloadInvalid),
		Message: "Payload seems to be invalid and cannot be processed.",
		Code:    400,
	}
}

func (e *InvalidPayloadError) Get() *requesterror.RequestError {
	return &requesterror.RequestError{
		Code:    e.Code,
		Message: e.Message,
		Name:    e.Name,
	}
}

// Global variables to init error type

var MissingParameterErr = &MissingParameterError{}

var HeaderMismatchErr = &HeaderMismatchError{}

var ValueNotUniqueErr = &ValueNotUniqueError{}

var WeakPasswordErr = &WeakPasswordError{}

var InvalidEmaiErr = &InvalidEmailError{}

var UsernameInvalidErr = &UsernameInvalidError{}

var InvalidPayloadErr = &InvalidPayloadError{}
