package util

import (
	"errors"
	"testing"
	"time"

	requesterror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/unauthorized"
)

var jwtKey = []byte("This is a super secret key")

func TestJWTGeneration(t *testing.T) {
	config := &Config{
		ApplicationName: "test-app",
		Session: struct {
			ExpireTime        int    "json:\"expireTimeMinutes,string\""
			RefreshExpireTime int    "json:\"refreshExpireTimeMinutes,string\""
			SecretKeyBase64   string "json:\"secretKey\""
			SecretKeyRaw      []byte "json:\"-\""
		}{
			ExpireTime:        5,
			RefreshExpireTime: 60,
			SecretKeyRaw:      jwtKey,
		},
	}

	claims := GenerateClaims(config, "aditya", "adityanotgeh@email.com", User)
	token, err := GenerateToken(claims, config.Session.SecretKeyRaw)

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(token)
	t.Log("Success")
}

func TestValidJWTToken(t *testing.T) {
	config := &Config{
		ApplicationName: "test-app",
		Session: struct {
			ExpireTime        int    "json:\"expireTimeMinutes,string\""
			RefreshExpireTime int    "json:\"refreshExpireTimeMinutes,string\""
			SecretKeyBase64   string "json:\"secretKey\""
			SecretKeyRaw      []byte "json:\"-\""
		}{
			ExpireTime:        5,
			RefreshExpireTime: 60,
			SecretKeyRaw:      jwtKey,
		},
	}

	claims := GenerateClaims(config, "aditya", "adityanotgeh@email.com", User)
	token, err := GenerateToken(claims, config.Session.SecretKeyRaw)

	if err != nil {
		t.Error(err)
		return
	}

	decodedClaims, err := VerifyToken(token, config.Session.SecretKeyRaw)

	if err != nil {
		t.Error(err)
		return
	}

	if decodedClaims.Username != claims.Username {
		t.Errorf("Claim username not match. expected %s received %s", claims.Username, decodedClaims.Username)
		return
	}

	t.Log("Success")
}

func TestInvalidToken(t *testing.T) {
	token := "xxxxxxxddddddd"
	config := &Config{
		ApplicationName: "test-app",
		Session: struct {
			ExpireTime        int    "json:\"expireTimeMinutes,string\""
			RefreshExpireTime int    "json:\"refreshExpireTimeMinutes,string\""
			SecretKeyBase64   string "json:\"secretKey\""
			SecretKeyRaw      []byte "json:\"-\""
		}{
			ExpireTime:        5,
			RefreshExpireTime: 60,
			SecretKeyRaw:      jwtKey,
		},
	}

	_, err := VerifyToken(token, config.Session.SecretKeyRaw)

	if !errors.As(err, &requesterror.InvalidTokenErr) {
		t.Errorf("Wrong error type found")
		t.Error(err)
		return
	}

	t.Log("Success")
}

func TestExpiredToken(t *testing.T) {
	config := &Config{
		ApplicationName: "test-app",
		Session: struct {
			ExpireTime        int    "json:\"expireTimeMinutes,string\""
			RefreshExpireTime int    "json:\"refreshExpireTimeMinutes,string\""
			SecretKeyBase64   string "json:\"secretKey\""
			SecretKeyRaw      []byte "json:\"-\""
		}{
			ExpireTime:        1,
			RefreshExpireTime: 60,
			SecretKeyRaw:      jwtKey,
		},
	}

	claims := GenerateClaims(config, "aditya", "adityanotgeh@email.com", User)
	token, err := GenerateToken(claims, config.Session.SecretKeyRaw)

	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Duration(1) * time.Minute)

	_, err = VerifyToken(token, config.Session.SecretKeyRaw)

	if !errors.As(err, &requesterror.TokenExpiredErr) {
		t.Errorf("Wrong error type")
		t.Error(err)
		return
	}

	t.Log("Success")
}
