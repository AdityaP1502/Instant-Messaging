package routes

type ErrorResponse struct {
	Status      string `json:"status"`
	ErrorType   string `json:"error_type"`
	Description string `json:"description"`
}

type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
