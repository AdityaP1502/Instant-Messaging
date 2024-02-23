package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/middleware"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/service/account/config"
)

var ClaimsKey middleware.ContextKey = "claims"
var TokenKey middleware.ContextKey = "token"

var AUTH_VERIFY_TOKEN_ENDPOINT = "v1/auth/token/verify"

func AuthMiddleware(next http.Handler, db *sql.DB, conf interface{}) http.Handler {
	fn := func(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) error {
		var token string
		var responseError *responseerror.ResponseError

		endpoint := r.Context().Value(middleware.EndpointKey).(string)

		cf := conf.(*config.Config)

		auth := r.Header.Get("Authorization")

		if auth == "" {
			return responseerror.EmptyAuthHeaderErr.Init()
		}

		if authType, authValue, _ := strings.Cut(auth, " "); authType != "Bearer" {
			return responseerror.InvalidAuthHeaderErr.Init(authType)

		} else {
			token = authValue
		}

		req := &httpx.HTTPRequest{}
		req, err := req.CreateRequest(
			cf.Services.Auth.Host,
			cf.Services.Auth.Port,
			AUTH_VERIFY_TOKEN_ENDPOINT,
			http.MethodPost,
			http.StatusOK,
			struct {
				Token    string `json:"token"`
				Endpoint string `json:"endpoint"`
			}{
				Token:    token,
				Endpoint: endpoint,
			},
		)

		if err != nil {
			return responseerror.InternalServiceErr.Init(err.Error())
		}

		err = req.Send(nil)

		if err != nil {
			if errors.As(err, &responseError) {
				return responseerror.TokenExpiredErr.Init()
			}

			return responseerror.InternalServiceErr.Init(err.Error())
		}

		return nil
	}

	return &httpx.Handler{
		DB:      db,
		Config:  conf,
		Handler: httpx.HandlerLogic(fn),
	}

}
