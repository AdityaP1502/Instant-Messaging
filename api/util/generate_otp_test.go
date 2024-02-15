package util

import (
	"testing"
)

func TestGenerateOTP(t *testing.T) {
	otp, err := GenerateOTP()
	if err != nil {
		t.Error(err)
	}

	if otp > 999999 || otp < 100000 {
		t.Errorf("OTP %d isn't a 6 digit number", otp)
		return
	}

	t.Logf("Success. %d is a valid 6 digit number", otp)
}
