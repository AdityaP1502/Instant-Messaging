package routes

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/AdityaP1502/Instant-Messanging/api/database"
	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/middleware"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/config"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/jwtutil"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/payload"
	"github.com/gorilla/mux"
)

type RevokedToken struct {
	TokenID   string `db:"token_id"`
	Username  string `db:"username"`
	Token     string `db:"token"`
	TokenType string `db:"type"`
	ExpiredAt string `db:"expired_at"`
}

var querynator = &database.Querynator{}

func IssueTokenHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) error {
	var roles jwtutil.Roles
	var ok bool

	cf := conf.(*config.Config)

	body := r.Context().Value(middleware.PayloadKey).(*payload.Credentials)

	if roles, ok = jwtutil.ParseRoles(body.Roles); !ok {
		return responseerror.InternalServiceErr.Init("Cannot parse given roles")
	}

	// Create a new token
	token := &payload.Token{}
	err := token.GenerateTokenPair(cf, body.Username, body.Email, roles)

	if err != nil {
		return responseerror.InternalServiceErr.Init(err.Error())
	}

	resp := &struct {
		Status string `json:"status"`
		*payload.Token
	}{
		Status: "success",
		Token:  token,
	}

	json, err := jsonutil.EncodeToJson(resp)

	if err != nil {
		return responseerror.InternalServiceErr.Init(err.Error())
	}

	w.WriteHeader(200)
	w.Write(json)

	return nil
}

func RefreshTokenHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) error {
	cf := conf.(*config.Config)

	body := r.Context().Value(middleware.PayloadKey).(*payload.Token)

	// Create a new token
	claims, err := body.CheckRefreshEligibility(cf)

	if err != nil {
		return err
	}

	// check if the token is revoked
	isExist, err := querynator.IsExists(&RevokedToken{Token: body.AccessToken}, db, "revoked_token")

	if err != nil {
		return responseerror.InternalServiceErr.Init(err.Error())
	}

	if isExist {
		return responseerror.InvalidTokenErr.Init("trying to refresh a revoked token")
	}

	newToken := &payload.Token{}

	err = newToken.GenerateTokenPair(cf, claims.Username, claims.Email, jwtutil.Roles(claims.Roles))
	if err != nil {
		return err
	}

	resp := &struct {
		Status string `json:"status"`
		*payload.Token
	}{
		Status: "success",
		Token:  newToken,
	}

	json, err := jsonutil.EncodeToJson(resp)

	if err != nil {
		return responseerror.InternalServiceErr.Init(err.Error())
	}

	w.WriteHeader(200)
	w.Write(json)

	return nil
}

func VerifyTokenHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) error {
	cf := conf.(*config.Config)
	body := r.Context().Value(middleware.PayloadKey).(*payload.Access)

	// verify token
	_, err := jwtutil.VerifyToken(body.AccessToken, cf.Session.SecretKeyRaw)

	if err != nil {
		return err
	}

	// TODO: check if current roles can access the endpoint
	w.WriteHeader(200)
	return nil
}

func RevokeTokenHandler(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) error {
	cf := conf.(*config.Config)
	body := r.Context().Value(middleware.PayloadKey).(*payload.Token)

	claims, err := jwtutil.VerifyToken(body.AccessToken, cf.Session.SecretKeyRaw)

	if err != nil {
		return err
	}

	_, err = querynator.Insert(&RevokedToken{
		Token:     body.AccessToken,
		TokenType: string(jwtutil.Auth),
		ExpiredAt: claims.ExpiresAt.Format(time.RFC3339),
	}, db, "revoked_token", "")

	if err != nil {
		return responseerror.InternalServiceErr.Init(err.Error())
	}

	w.WriteHeader(200)

	return nil
}

func SetAuthRoute(r *mux.Router, db *sql.DB, conf *config.Config) {
	subrouter := r.PathPrefix("/auth").Subrouter()

	credentialsPayloadMiddleware, err := middleware.PayloadCheckMiddleware(&payload.Credentials{}, "Username", "Roles", "Email")

	if err != nil {
		log.Fatal(err)
	}

	refreshpayloadMiddleware, err := middleware.PayloadCheckMiddleware(&payload.Token{}, "RefreshToken", "AccessToken")

	if err != nil {
		log.Fatal(err)
	}

	revokepayloadMiddleware, err := middleware.PayloadCheckMiddleware(&payload.Token{}, "AccessToken")

	if err != nil {
		log.Fatal(err)
	}

	accesspayloadMiddleware, err := middleware.PayloadCheckMiddleware(&payload.Access{}, "Endpoint", "AccessToken")

	if err != nil {
		log.Fatal(err)
	}

	issueToken := &httpx.Handler{
		DB:      db,
		Config:  conf,
		Handler: httpx.HandlerLogic(IssueTokenHandler),
	}

	refreshToken := &httpx.Handler{
		DB:      db,
		Config:  conf,
		Handler: httpx.HandlerLogic(RefreshTokenHandler),
	}

	verifyToken := &httpx.Handler{
		DB:      db,
		Config:  conf,
		Handler: httpx.HandlerLogic(VerifyTokenHandler),
	}

	revokeToken := &httpx.Handler{
		DB:      db,
		Config:  conf,
		Handler: httpx.HandlerLogic(RevokeTokenHandler),
	}

	subrouter.Handle("/token/issue", middleware.UseMiddleware(db, conf, issueToken, credentialsPayloadMiddleware))
	subrouter.Handle("/token/refresh", middleware.UseMiddleware(db, conf, refreshToken, refreshpayloadMiddleware))
	subrouter.Handle("/token/verify", middleware.UseMiddleware(db, conf, verifyToken, accesspayloadMiddleware))
	subrouter.Handle("/token/revoke", middleware.UseMiddleware(db, conf, revokeToken, revokepayloadMiddleware))

}