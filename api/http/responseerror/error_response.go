package responseerror

type ErrorResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorType string `json:"error_type"`
}
