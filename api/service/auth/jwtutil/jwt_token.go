package jwtutil

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/service/auth/config"
	"github.com/golang-jwt/jwt/v4"
)

type AccessType string
type Roles string

var JWT_SIGNING_METHOD jwt.SigningMethod = jwt.SigningMethodHS256

const (
	Refresh AccessType = "Refresh"
	Auth    AccessType = "Auth"
)

const (
	User  Roles = "user"
	Admin Roles = "admin"
)

var (
	rolesMap = map[string]Roles{
		"user":  User,
		"admin": Admin,
	}
)

func ParseRoles(str string) (Roles, bool) {
	r, ok := rolesMap[strings.ToLower(str)]

	return r, ok
}

type Claims struct {
	Username   string
	Email      string
	AccessType AccessType
	Roles      string
	jwt.RegisteredClaims
}

func GenerateClaims(config *config.Config, username string, email string, role Roles) *Claims {
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

func GenerateRefreshClaims(config *config.Config, username string, email string, role Roles) *Claims {
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

func VerifyToken(tokenString string, key []byte) (*Claims, responseerror.HTTPCustomError) {
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
			return &claims, responseerror.CreateUnauthorizedError(
				responseerror.TokenExpired,
				responseerror.TokenExpiredMessage,
				nil,
			)
		}

		return nil, responseerror.CreateUnauthorizedError(
			responseerror.InvalidToken,
			responseerror.InvalidTokenMessage,
			map[string]string{
				"description": err.Error(),
			},
		)
	}

	if !token.Valid {
		return nil, responseerror.CreateUnauthorizedError(
			responseerror.InvalidToken,
			responseerror.InvalidTokenMessage,
			map[string]string{
				"description": "",
			},
		)
	}

	// claims, ok := token.Claims.(*Claims)

	// if !ok {
	// 	return nil, requesterror.InvalidTokenErr.Init("Unrecognized claims")
	// }

	// }

	return &claims, nil
}
