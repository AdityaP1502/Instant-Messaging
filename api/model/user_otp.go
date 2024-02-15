package model

import (
	"io"

	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
)

type UserOTP struct {
	OTPID             string `json:"-" db:"otp_id"`
	Username          string `json:"-" db:"username"`
	OTPConfirmID      string `json:"otp_confirmation_id" db:"otp_confirmation_id"`
	OTP               string `json:"otp" db:"otp"`
	LastResend        string `json:"-" db:"last_resend"`
	MarkedForDeletion string `json:"-" db:"marked_for_deletion"`
}

func (o *UserOTP) FromJSON(r io.Reader, checkRequired bool) error {
	err := util.DecodeJSONBody(r, o)

	if err != nil {
		return err
	}

	if checkRequired {
		return util.CheckParametersUnity(o)
	}

	return nil
}

func (o *UserOTP) ToJSON(checkRequired bool) ([]byte, error) {
	var err error

	if checkRequired {
		if err = util.CheckParametersUnity(o); err != nil {
			return nil, err
		}
	}

	return util.CreateJSONResponse(o)
}
