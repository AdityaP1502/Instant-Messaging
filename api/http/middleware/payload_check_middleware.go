package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
)

var PayloadKey ContextKey = "payload"

func PayloadCheckMiddleware(template httpx.Payload, requiredFields ...string) (Middleware, error) {
	var payload httpx.Payload

	p := reflect.ValueOf(template)
	if p.Kind() != reflect.Ptr {
		err := fmt.Errorf("cannot create middleware. template isn't a pointer")
		return nil, err
	}

	// Struct reflect
	s := p.Elem()
	sType := s.Type()

	// Check if the requiredFields is valid
	for _, field := range requiredFields {
		f := s.FieldByName(field)
		if !f.IsValid() {
			return nil, fmt.Errorf("struct of %s don't have field named %s", s.Type(), field)
		}
	}

	return func(next http.Handler, db *sql.DB, config interface{}) http.Handler {
		fn := func(db *sql.DB, config interface{}, w http.ResponseWriter, r *http.Request) responseerror.HTTPCustomError {

			if r.Header.Get("Content-Type") != "application/json" {
				return responseerror.CreateBadRequestError(
					responseerror.HeaderValueMistmatch,
					responseerror.HeaderValueMistmatchMessage,
					map[string]string{
						"name": "Content-Type",
					},
				)
			}

			// create a new struct with the same type as template
			payload = reflect.New(sType).Interface().(httpx.Payload)

			if r.Body == nil {
				return responseerror.CreateBadRequestError(
					responseerror.PayloadInvalid,
					responseerror.PayloadInvalidMessage,
					nil,
				)
			}

			err := payload.FromJSON(r.Body, true, requiredFields)

			if err != nil {
				return responseerror.CreateInternalServiceError(err)
			}

			ctx := context.WithValue(r.Context(), PayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))

			return nil
		}

		return &httpx.Handler{
			DB:      db,
			Config:  config,
			Handler: httpx.HandlerLogic(fn),
		}
	}, nil
}
