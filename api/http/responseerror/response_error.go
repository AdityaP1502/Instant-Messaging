package responseerror

import "fmt"

type errorType string

type ResponseError struct {
	Code    int
	Message string
	Name    string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

type ResponseErrorGetter interface {
	Get() *ResponseError
}
