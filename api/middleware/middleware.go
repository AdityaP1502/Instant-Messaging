package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/AdityaP1502/Instant-Messaging/api/api/model"
	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
	badrequest "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/bad_request"
	"github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/unauthorized"
	mapset "github.com/deckarep/golang-set/v2"
)

type ContextKey string

const (
	PayloadKey ContextKey = "payload"
	ClaimsKey  ContextKey = "Claims"
)

type middleware func(next http.Handler, db *sql.DB, config *util.Config) http.Handler

func UseMiddleware(db *sql.DB, config *util.Config, handler http.Handler, middlewares ...middleware) http.Handler {
	chained := handler

	for i := len(middlewares) - 1; i > -1; i-- {
		chained = middlewares[i](chained, db, config)
	}

	return chained
}

func AuthMiddleware(next http.Handler, db *sql.DB, config *util.Config) http.Handler {
	fn := func(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
		// Check for authorization header

		var token string

		auth := r.Header.Get("Authorization")

		if auth == "" {
			return unauthorized.EmptyAuthHeaderErr.Init()
		}

		if authType, authValue, _ := strings.Cut(auth, " "); authType != "Bearer" {
			return unauthorized.InvalidAuthHeaderErr.Init(authType)
		} else {
			token = authValue
		}

		claims, err := util.VerifyToken(token, config.Session.SecretKeyRaw)

		if err != nil {
			return err
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))

		return nil
	}

	return &Handler{
		DB:      db,
		Config:  config,
		Handler: HandlerLogic(fn),
	}
}

func IdentitiyAccessManagementMiddleware(allowAccessRole ...string) (middleware, error) {
	roleAccess := mapset.NewSet[string](allowAccessRole...)

	return func(next http.Handler, db *sql.DB, config *util.Config) http.Handler {
		fn := func(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
			// Do access checking here
			claims := r.Context().Value(ClaimsKey).(*util.Claims)

			if !roleAccess.Contains(string(claims.Roles)) {
				// TODO: Change into 404 Not Found Error
				return fmt.Errorf("Path not found")
			}

			next.ServeHTTP(w, r)

			return nil
		}

		return &Handler{
			DB:      db,
			Config:  config,
			Handler: HandlerLogic(fn),
		}
	}, nil
}

func PayloadCheckMiddleware(template model.Model) (middleware, error) {
	var payload model.Model

	if reflect.ValueOf(template).Kind() != reflect.Ptr {
		err := fmt.Errorf("Cannot create middleware. template isn't a pointer")
		return nil, err
	}

	return func(next http.Handler, db *sql.DB, config *util.Config) http.Handler {
		fn := func(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {

			if r.Header.Get("Content-Type") != "application/json" {
				return badrequest.HeaderMismatchErr.Init("Content-Type")
			}

			payload = reflect.New(reflect.ValueOf(template).Elem().Type()).Interface().(model.Model)

			if r.Body == nil {
				return badrequest.InvalidPayloadErr.Init()
			}

			err := payload.FromJSON(r.Body, true)

			if err != nil {
				return err
			}

			ctx := context.WithValue(r.Context(), PayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))

			return nil
		}

		return &Handler{
			DB:      db,
			Config:  config,
			Handler: HandlerLogic(fn),
		}
	}, nil
}
