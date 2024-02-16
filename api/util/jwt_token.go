package util

import (
	"errors"
	"fmt"
	"time"

	requesterror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/unauthorized"
	"github.com/golang-jwt/jwt/v4"
)

type AccessType string
type Roles string

var JWT_SIGNING_METHOD jwt.SigningMethod = jwt.SigningMethodHS256

const (
	Refresh AccessType = "Refresh"
	Auth    AccessType = "Auth"
	Basic   AccessType = "Basic"
)

const (
	User  Roles = "User"
	Admin Roles = "Admin"
)

type Claims struct {
	Username   string
	Email      string
	AccessType AccessType
	Roles      string
	jwt.RegisteredClaims
}

func GenerateBasicClaims(config *Config, username string, email string, role Roles) *Claims {
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.OTP.OTPDurationMinutes) * time.Minute)),
		},
		Username:   username,
		Email:      email,
		AccessType: Auth,
		Roles:      string(role),
	}
}
func GenerateClaims(config *Config, username string, email string, role Roles) *Claims {
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Session.ExpireTime) * time.Minute)),
		},
		Username:   username,
		Email:      email,
		AccessType: Auth,
		Roles:      string(role),
	}
}

func GenerateRefreshClaims(config *Config, username string, email string, role Roles) *Claims {
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Session.ExpireTime) * time.Minute)),
		},
		Username:   username,
		Email:      email,
		AccessType: Refresh,
		Roles:      string(role),
	}
}
func GenerateToken(claim *Claims, key []byte) (string, error) {
	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claim)
	signedToken, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(tokenString string, key []byte) (*Claims, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method invalid")
			} else if method != JWT_SIGNING_METHOD {
				return nil, fmt.Errorf("signing method invalid")
			}

			return key, nil
		})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return &claims, requesterror.TokenExpiredErr.Init()
		}

		return nil, requesterror.InvalidTokenErr.Init(err.Error())
	}

	if !token.Valid {
		return nil, requesterror.InvalidTokenErr.Init("")
	}

	// claims, ok := token.Claims.(*Claims)

	// if !ok {
	// 	return nil, requesterror.InvalidTokenErr.Init("Unrecognized claims")
	// }

	// }

	return &claims, nil
}
