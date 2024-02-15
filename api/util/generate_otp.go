package util

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP() (int, error) {
	// Generate user OTP
	n, err := rand.Int(rand.Reader, big.NewInt(900000))

	if err != nil {
		return 0, nil
	}

	otp := n.Add(n, big.NewInt(100000)) // produce a 6 digit otp

	return int(otp.Int64()), nil
}
