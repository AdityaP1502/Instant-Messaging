package httpx

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

type HTTPRequest struct {
	Request            http.Request
	Payload            []byte
	SuccessStatusCode  int
	ReturnedStatusCode int
	Status             int
	IsSuccess          bool
}

func (h *HTTPRequest) CreateRequest(host string, port int, endpoint string, method string, successStatus int, payload interface{}) (*HTTPRequest, error) {
	url := fmt.Sprintf("http://%s:%d/%s", host, port, endpoint)

	json, err := jsonutil.EncodeToJson(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(json))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return &HTTPRequest{
		Request:           *req,
		Payload:           nil,
		SuccessStatusCode: successStatus,
	}, nil
}

func (h *HTTPRequest) Send(dest interface{}) error {
	var client = &http.Client{}

	resp, err := client.Do(&h.Request)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != h.SuccessStatusCode {
		// if not provided a destination or that the status code don't match
		// the expected return code
		// store the payload in the payload field

		respBytes, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(respBytes))

		errorResponse := &responseerror.ErrorResponse{}
		err = jsonutil.DecodeJSON(resp.Body, errorResponse)

		if err != nil {
			return err
		}

		return &responseerror.ResponseError{
			Code:    resp.StatusCode,
			Message: errorResponse.Message,
			Name:    errorResponse.ErrorType,
		}
	}

	if dest == nil {
		return nil
	}

	return jsonutil.DecodeJSON(resp.Body, dest)
}
