package routes

import (
	"database/sql"
	"net/http"

	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
	"github.com/gorilla/mux"
)

type LoginResponse struct {
	Status string `json:"status"`
	Token  struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"token"`
}

type RegisterResponse struct {
	Status       string `json:"status"`
	OTPConfirmID string `json:"otp_confirmation_id"`
	AccessToken  string `json:"access_token"`
}

func registerHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func loginHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func resendOTPHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func verifyOTPHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func refreshTokenHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func logOutHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func patchUserInfoHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Register account subrouter
func SetAccountRoute(r *mux.Router, db *sql.DB, config *util.Config) {
	subrouter := r.PathPrefix("/account/").Subrouter()

	// register account path and its handler here

	register := &util.Handler{
		DB:      db,
		Config:  config,
		Handler: registerHandler,
	}

	login := &util.Handler{
		DB:      db,
		Config:  config,
		Handler: loginHandler,
	}

	subrouter.Handle("/register", register).Methods("POST")
	subrouter.Handle("/login", login).Methods("POST")

	// subrouter.HandleFunc("/logout", logOutHandler).Methods("POST")
	// subrouter.HandleFunc("/otp/verify", loginHandler).Methods("POST")
	// subrouter.HandleFunc("/otp/{otp_confimartion_id}/resend", resendOTPHandler).Methods("POST")
	// subrouter.HandleFunc("/token/refresh", refreshTokenHandler).Methods("POST")
	// subrouter.HandleFunc("/{username}", patchUserInfoHandler).Methods("PATCH")
}
