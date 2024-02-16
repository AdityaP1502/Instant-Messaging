package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/AdityaP1502/Instant-Messaging/api/api/database"
	"github.com/AdityaP1502/Instant-Messaging/api/api/model"
	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
	badrequest "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/bad_request"
	internalserviceerror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/internal_service_error"
	notfound "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/not_found"
	"github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/unauthorized"
	mapset "github.com/deckarep/golang-set/v2"
)

type ContextKey string

const (
	PayloadKey ContextKey = "payload"
	ClaimsKey  ContextKey = "Claims"
	TokenKey   ContextKey = "token"
)

type middleware func(next http.Handler, db *sql.DB, config *util.Config) http.Handler

func UseMiddleware(db *sql.DB, config *util.Config, handler http.Handler, middlewares ...middleware) http.Handler {
	chained := handler

	for i := len(middlewares) - 1; i > -1; i-- {
		chained = middlewares[i](chained, db, config)
	}

	return chained
}

func AuthMiddleware(allowedTokenType ...util.AccessType) (middleware, error) {
	allowdTypes := mapset.NewSet[util.AccessType](allowedTokenType...)

	return func(next http.Handler, db *sql.DB, config *util.Config) http.Handler {
		fn := func(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
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

			if !allowdTypes.Contains(claims.AccessType) {
				return unauthorized.InvalidTokenErr.Init("Your token don't have the required access.")
			}

			// check if revoked
			querynator := database.Querynator{}

			isRevoked, err := querynator.IsExists(&model.RevokedToken{Token: token, TokenType: string(claims.AccessType)}, db, "revoked_token")

			if err != nil {
				return internalserviceerror.InternalServiceErr.Init(err.Error())
			}

			if isRevoked {
				return unauthorized.InvalidTokenErr.Init("Token is revoked")
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			ctx = context.WithValue(ctx, TokenKey, token)
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

func IdentitiyAccessManagementMiddleware(allowAccessRole ...string) (middleware, error) {
	roleAccess := mapset.NewSet[string](allowAccessRole...)

	return func(next http.Handler, db *sql.DB, config *util.Config) http.Handler {
		fn := func(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
			// Do access checking here
			claims := r.Context().Value(ClaimsKey).(*util.Claims)

			if !roleAccess.Contains(string(claims.Roles)) {
				return notfound.NotFoundErr.Init("path", "Path")
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

func PayloadCheckMiddleware(template model.Model, requiredFields ...string) (middleware, error) {
	var payload model.Model

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

	return func(next http.Handler, db *sql.DB, config *util.Config) http.Handler {
		fn := func(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {

			if r.Header.Get("Content-Type") != "application/json" {
				return badrequest.HeaderMismatchErr.Init("Content-Type")
			}

			// create a new struct with the same type as template
			payload = reflect.New(sType).Interface().(model.Model)

			if r.Body == nil {
				return badrequest.InvalidPayloadErr.Init()
			}

			err := payload.FromJSON(r.Body, true, requiredFields)

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
