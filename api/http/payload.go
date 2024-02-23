package httpx

import "io"

type Payload interface {
	FromJSON(r io.Reader, checkRequired bool, requiredFields []string) error
	ToJSON(checkRequired bool, requiredFields []string) ([]byte, error)
}
