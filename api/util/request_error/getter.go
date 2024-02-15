package requesterror

type RequestError struct {
	Code    int
	Message string
	Name    string
}

type CustomErrorGetter interface {
	Get() *RequestError
}
