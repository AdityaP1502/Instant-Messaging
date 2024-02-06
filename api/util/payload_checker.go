package util

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	requesterror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/bad_request"
	mapset "github.com/deckarep/golang-set/v2"
)

// CheckParameterUnity check if all required json parameers are filled
//
// if a field with tag json and have a value empty, the function will
// return with a non nil value.
// This function will ignore field with "omitempty" tag
func CheckParametersUnity(v interface{}) error {
	// get interface field
	s := reflect.ValueOf(v)
	typeOfS := s.Type()

	for i := 0; i < typeOfS.NumField(); i++ {
		field := typeOfS.Field(i)
		jsonTag := field.Tag.Get("json")

		// Gatekeep conditional
		if jsonTag == "-" || jsonTag == "" {
			continue
		}

		if x := strings.SplitAfter(jsonTag, ","); len(x) > 1 {
			if x[1] == "omitempty" {
				continue
			}
		}

		a := s.Field(i).Interface()

		// Check if a field is empty/has value of "zero"
		if a == reflect.Zero(s.Field(i).Type()).Interface() {
			return requesterror.MissingParameterErr.Init(jsonTag)
		}
	}

	return nil
}

// check if a request header match with the expected value
func CheckHeader(h http.Header, headerName []string, expectedValue []mapset.Set[string]) error {
	if len(headerName) != len(expectedValue) {
		return errors.New("Name length isn't equal with value length")
	}

	for i := 0; i < len(headerName); i++ {
		if !expectedValue[i].Contains(h.Get(headerName[i])) {
			return requesterror.HeaderMismatchErr.Init(headerName[i])
		}
	}

	return nil
}
