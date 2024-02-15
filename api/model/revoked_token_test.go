package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/AdityaP1502/Instant-Messaging/api/api/database"
	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
)

func TestInsertRevokedToken(t *testing.T) {
	config := &util.Config{
		ApplicationName: "test-app",
		Session: struct {
			ExpireTime      int    "json:\"expireTime,string\""
			SecretKeyBase64 string "json:\"secretKey\""
			SecretKeyRaw    []byte "json:\"-\""
		}{
			ExpireTime:   5,
			SecretKeyRaw: []byte("Thisissecretkey"),
		},
	}

	db := connectToDB(t)

	claims := util.GenerateClaims(config, "lucas", "lucas12@gmail.com", util.Basic, util.User)
	token, err := util.GenerateToken(claims, config.Hash.SecretKeyRaw)

	if err != nil {
		t.Error(err.Error())
		return
	}

	querynator := database.Querynator{}

	_, err = querynator.Insert(&RevokedToken{Token: token, Username: "lucas", ExpiredAt: claims.ExpiresAt.Local().Format(time.RFC3339), TokenType: string(claims.AccessType)}, db.DB, "revoked_token", "token_id")

	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log("Success")
}

func TestSearchToken(t *testing.T) {
	config := &util.Config{
		ApplicationName: "test-app",
		Session: struct {
			ExpireTime      int    "json:\"expireTime,string\""
			SecretKeyBase64 string "json:\"secretKey\""
			SecretKeyRaw    []byte "json:\"-\""
		}{
			ExpireTime:   5,
			SecretKeyRaw: []byte("Thisissecretkey"),
		},
	}

	db := connectToDB(t)

	claims := util.GenerateClaims(config, "lucas14", "lucas14@gmail.com", util.Basic, util.User)
	token, err := util.GenerateToken(claims, config.Hash.SecretKeyRaw)

	if err != nil {
		t.Error(err.Error())
		return
	}

	querynator := database.Querynator{}

	id, err := querynator.Insert(&RevokedToken{Token: token, Username: "lucas14", ExpiredAt: claims.ExpiresAt.Local().Format(time.RFC3339), TokenType: string(claims.AccessType)}, db.DB, "revoked_token", "token_id")

	if err != nil {
		t.Error(err.Error())
		return
	}

	tokenData := &RevokedToken{}

	err = querynator.FindOne(&RevokedToken{TokenID: fmt.Sprintf("%d", id)}, tokenData, db.DB, "revoked_token", "expired_at", "token")

	if err != nil {
		t.Error(err.Error())
		return
	}

	expiredTime, err := util.ParseTimestamp(tokenData.ExpiredAt)

	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("Revoked token will expire in : %s", expiredTime.In(time.Local).Format(time.RFC3339))
	t.Log("Success")
}
