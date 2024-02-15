package model

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/AdityaP1502/Instant-Messaging/api/api/database"
	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
	"github.com/google/uuid"
)

func TestInserUserOTP(t *testing.T) {
	db := connectToDB(t)

	querynator := database.Querynator{}
	_, err := querynator.Insert(newUser, db.DB, "account", "")

	if err != nil {
		t.Error(err)
		return
	}

	otp, err := util.GenerateOTP()

	otpData := &UserOTP{
		Username:          newUser.Username,
		OTPConfirmID:      uuid.NewString(),
		OTP:               fmt.Sprintf("%d", otp),
		LastResend:        util.GenerateTimestamp(),
		MarkedForDeletion: strconv.FormatBool(false),
	}

	otpID, err := querynator.Insert(otpData, db.DB, "user_otp", "otp_id")

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Success. The ID is %d", otpID)
}

func TestUpdateOTP(t *testing.T) {
	db := connectToDB(t)

	data := &Account{
		Username: "Dimwit",
		Email:    "email@domain.com",
		Name:     "MyGuyisGay",
		Salt:     "random_ah_salt",
		Password: "Test123",
		IsActive: strconv.FormatBool(false),
	}

	querynator := database.Querynator{}
	_, err := querynator.Insert(data, db.DB, "account", "")

	if err != nil {
		t.Error(err)
		return
	}

	otp, err := util.GenerateOTP()

	otpData := &UserOTP{
		Username:          data.Username,
		OTPConfirmID:      uuid.NewString(),
		OTP:               fmt.Sprintf("%d", otp),
		LastResend:        util.GenerateTimestamp(),
		MarkedForDeletion: strconv.FormatBool(false),
	}

	_, err = querynator.Insert(otpData, db.DB, "user_otp", "")

	if err != nil {
		t.Error(err)
		return
	}

	searchData := &Account{}
	err = querynator.FindOne(&Account{Username: data.Username}, searchData, db.DB, "account", "username")

	if err != nil {
		t.Error(err)
		return
	}

	otp, err = util.GenerateOTP()

	err = querynator.Update(&UserOTP{
		OTP:        fmt.Sprintf("%d", otp),
		LastResend: util.GenerateTimestamp(),
	}, []string{"username", "otp_confirmation_id"}, []any{searchData.Username, otpData.OTPConfirmID}, db.DB, "user_otp")

	searchOTP := &UserOTP{}

	err = querynator.FindOne(&UserOTP{OTPConfirmID: otpData.OTPConfirmID}, searchOTP, db.DB, "user_otp", "otp", "last_resend")

	if err != nil {
		t.Error(err)
		return
	}

	updateOTPd, err := strconv.Atoi(searchOTP.OTP)
	if err != nil {
		t.Error(err)
		return
	}

	if updateOTPd != otp {
		t.Errorf("Data hasn't been updated properly")
		return
	}

	t.Log("Success")
}
