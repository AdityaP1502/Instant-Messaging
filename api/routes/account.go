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
	notfound "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/not_found"
	toomanyrequest "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/too_many_request"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var querynator = &database.Querynator{}

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

	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to send mail. mail server return with status code %d", resp.StatusCode)
	}

	return nil
}

func registerHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	payload := r.Context().Value(middleware.PayloadKey).(*model.Account)

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

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

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

// func loginHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }

func resendOTPHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	confirmID := vars["otp_confirmation_id"]
	u := &model.UserOTP{OTPConfirmID: confirmID}

	err := querynator.FindOne(&model.UserOTP{OTPConfirmID: confirmID}, u, db, "user_otp",
		"otp_id",
		"last_resend",
	)

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	// Check last resend duration
	t, err := util.ParseTimestamp(u.LastResend)

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	duration := util.SecondsDifferenceFromNow(t)

	if duration < config.OTP.ResendDurationMinutes {
		return toomanyrequest.ResendIntervalNotReachedErr.Init()
	}

	// revoked user token
	claims := r.Context().Value(middleware.PayloadKey).(*util.Claims)
	token := r.Context().Value(middleware.TokenKey).(string)

	tId, err := querynator.Insert(&model.RevokedToken{Token: token, ExpiredAt: claims.ExpiresAt.UTC().Weekday().String(), TokenType: string(claims.AccessType)}, db, "revoked_token", "token_id")

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	// create a new user token
	newClaims := util.GenerateClaims(config, claims.Username, claims.Email, util.Basic, util.User)
	newToken, err := util.GenerateToken(newClaims, config.Session.SecretKeyRaw)

	if err != nil {
		querynator.Delete(&model.RevokedToken{TokenID: fmt.Sprintf("%d", tId)}, db, "user_otp")
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	otp, err := util.GenerateOTP()

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	// send the new otp
	sendMailHTTP(fmt.Sprintf("Your OTP is %d. Don't share with anyone.", otp),
		"User Verification",
		claims.Email,
		fmt.Sprintf("http://%s:%d/mail/send", config.MailAPI.Host, config.MailAPI.Port),
	)

	json, err := util.CreateJSONResponse(struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}{Status: "success", Message: "OTP has been re-send to your email.", Token: newToken})

	if err != nil {
		querynator.Delete(&model.RevokedToken{TokenID: fmt.Sprintf("%d", tId)}, db, "user_otp")
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	// Update the otp
	err = querynator.Update(&model.UserOTP{OTP: fmt.Sprintf("%d", otp)}, []string{"otp_id"}, []any{u.OTPID}, db, "user_otp")

	if err != nil {
		querynator.Delete(&model.RevokedToken{TokenID: fmt.Sprintf("%d", tId)}, db, "user_otp")
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	w.WriteHeader(200)
	w.Write(json)

	return nil
}

func verifyOTPHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
	var validOTP = &model.UserOTP{}

	payload := r.Context().Value(middleware.PayloadKey).(*model.UserOTP)
	claims := r.Context().Value(middleware.ClaimsKey).(*util.Claims)

	// Fill the username in the payload
	payload.Username = claims.Username

	err := querynator.FindOne(&model.UserOTP{Username: payload.Username, OTPConfirmID: payload.OTPConfirmID}, validOTP, db, "user_otp", "otp")

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	if validOTP.OTP == "" {
		// not found entry
		return notfound.NotFoundErr.Init("otp_confirmation_id", "OTP Confirmation ID")
	}

	if validOTP.OTP != payload.OTP {
		// otp is wrong
		return badrequest.InvalidOTPErr.Init()
	}

	json, err := util.CreateJSONResponse(&model.GenericResponse{Status: "success", Message: "your account has been activated successfully"})

	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	// otp is correct, update user to be an active user
	err = querynator.Update(&model.Account{IsActive: strconv.FormatBool(true)}, []string{"username"}, []any{claims.Username}, db, "account")
	if err != nil {
		return internalserviceerror.InternalServiceErr.Init(err.Error())
	}

	w.WriteHeader(200)
	w.Write(json)

	return nil
}

// func refreshTokenHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }

// func logOutHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }

// func patchUserInfoHandler(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }

// Register account subrouter
func SetAccountRoute(r *mux.Router, db *sql.DB, config *util.Config) {
	subrouter := r.PathPrefix("/account").Subrouter()

	// Create middleware here
	userPayloadMiddleware, err := middleware.PayloadCheckMiddleware(&model.Account{})

	if err != nil {
		log.Fatal(err)
	}

	otpPayloadMiddleware, err := middleware.PayloadCheckMiddleware(&model.UserOTP{})

	if err != nil {
		log.Fatal(err)
	}

	basicAccessAuthMiddleware, _ := middleware.AuthMiddleware(string(util.Basic))

	// REGISTER ROUTE //
	register := &middleware.Handler{
		DB:      db,
		Config:  config,
		Handler: registerHandler,
	}

	subrouter.Handle("/register", middleware.UseMiddleware(db, config, register, userPayloadMiddleware)).Methods("POST")

	// VERIFY OTP ROUTE //
	verifyOTP := &middleware.Handler{
		DB:      db,
		Config:  config,
		Handler: verifyOTPHandler,
	}

	subrouter.Handle("/otp/verify", middleware.UseMiddleware(db, config, verifyOTP, otpPayloadMiddleware, basicAccessAuthMiddleware)).Methods("POST")

	// RESEND OTP ROUTE //
	resendOTP := &middleware.Handler{
		DB:      db,
		Config:  config,
		Handler: resendOTPHandler,
	}

	subrouter.Handle("/otp/{otp_confimartion_id}/resend", middleware.UseMiddleware(db, config, resendOTP, otpPayloadMiddleware, basicAccessAuthMiddleware)).Methods("POST")

	// subrouter.Handle("/login", login).Methods("POST")

	// subrouter.HandleFunc("/logout", logOutHandler).Methods("POST")
	// subrouter.HandleFunc("/token/refresh", refreshTokenHandler).Methods("POST")
	// subrouter.HandleFunc("/{username}", patchUserInfoHandler).Methods("PATCH")
}
