package routes

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AdityaP1502/Instant-Messaging/api/api/database"
	"github.com/AdityaP1502/Instant-Messaging/api/api/middleware"
	"github.com/AdityaP1502/Instant-Messaging/api/api/model"
	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
	badrequest "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/bad_request"
	internalserviceerror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/internal_service_error"
	"github.com/google/uuid"
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

func sendMailHTTP(message string, subject string, to string, url string) error {
	//TODO: Send http request to node js server

	var client = &http.Client{}

	var mail struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	mail.To = to
	mail.Subject = subject
	mail.Message = message

	json, err := util.CreateJSONResponse(&mail)

	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST", url, bytes.NewReader(json))
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to send mail. Mail server return with status code %d", resp.StatusCode)
	}

	return nil
}

func registerHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	payload := r.Context().Value(middleware.PayloadKey).(*model.Account)
	querynator := database.Querynator{}

	// Check username and email exist or not
	exist, err := querynator.IsExists(&model.Account{Email: payload.Email}, db, "account")

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	if exist {
		return badrequest.ValueNotUniqueErr.Init(badrequest.EmailExists, "email")
	}

	exist, err = querynator.IsExists(&model.Account{Username: payload.Username}, db, "account")
	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	if exist {
		return badrequest.ValueNotUniqueErr.Init(badrequest.UsernameExists, "username")
	}

	hash, salt, err := util.HashPassword(payload.Password, config.Hash.SecretKeyRaw)

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	// set the hash password and salt to the model for storage
	payload.Password = hash
	payload.Salt = salt
	payload.IsActive = strconv.FormatBool(false)

	id, err := querynator.Insert(payload, db, "account", "account_id")

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	otp, err := util.GenerateOTP()

	otpUUID := uuid.NewString()

	otpData := &model.UserOTP{
		Username:          payload.Username,
		OTP:               fmt.Sprintf("%d", otp),
		LastResend:        util.GenerateTimestamp(),
		OTPConfirmID:      otpUUID,
		MarkedForDeletion: strconv.FormatBool(false),
	}

	otpID, err := querynator.Insert(otpData, db, "user_otp", "otp_id")

	if err != nil {
		querynator.Delete(&model.Account{AccountID: fmt.Sprintf("%d", id)}, db, "account")
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	claims := util.GenerateClaims(config, payload.Username, payload.Email, util.Basic, util.User)
	token, err := util.GenerateToken(claims, config.Session.SecretKeyRaw)

	if err != nil {
		querynator.Delete(&model.UserOTP{OTPID: fmt.Sprintf("%d", otpID)}, db, "user_otp")
		querynator.Delete(&model.Account{AccountID: fmt.Sprintf("%d", id)}, db, "account")
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	// send mail
	err = sendMailHTTP(
		fmt.Sprintf("Your OTP is %s. Don't share with anyone.", otpData.OTP),
		"User Verification",
		payload.Email,
		fmt.Sprintf("http://%s:%d/mail/send", config.MailAPI.Host, config.MailAPI.Port),
	)

	if err != nil {
		querynator.Delete(&model.UserOTP{OTPID: fmt.Sprintf("%d", otpID)}, db, "user_otp")
		querynator.Delete(&model.Account{AccountID: fmt.Sprintf("%d", id)}, db, "account")
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	resp := &RegisterResponse{
		Status:       "success",
		OTPConfirmID: otpUUID,
		AccessToken:  token,
	}

	jsonResponse, err := util.CreateJSONResponse(resp)

	if err != nil {
		querynator.Delete(&model.UserOTP{OTPID: fmt.Sprintf("%d", otpID)}, db, "user_otp")
		querynator.Delete(&model.Account{AccountID: fmt.Sprintf("%d", id)}, db, "account")
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	w.WriteHeader(200)
	w.Write(jsonResponse)

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
	subrouter := r.PathPrefix("/account").Subrouter()

	subrouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}).Methods("GET")
	// register account path and its handler here

	register := &middleware.Handler{
		DB:      db,
		Config:  config,
		Handler: registerHandler,
	}

	// create a payload middleware
	userPayloadMiddleware, err := middleware.PayloadCheckMiddleware(&model.Account{})

	if err != nil {
		log.Fatal(err)
	}

	subrouter.Handle("/register", middleware.UseMiddleware(db, config, register, userPayloadMiddleware)).Methods("POST")

	// login := &util.Handler{
	// 	DB:      db,
	// 	Config:  config,
	// 	Handler: loginHandler,
	// }

	// subrouter.Handle("/login", login).Methods("POST")

	// subrouter.HandleFunc("/logout", logOutHandler).Methods("POST")
	// subrouter.HandleFunc("/otp/verify", loginHandler).Methods("POST")
	// subrouter.HandleFunc("/otp/{otp_confimartion_id}/resend", resendOTPHandler).Methods("POST")
	// subrouter.HandleFunc("/token/refresh", refreshTokenHandler).Methods("POST")
	// subrouter.HandleFunc("/{username}", patchUserInfoHandler).Methods("PATCH")
}
