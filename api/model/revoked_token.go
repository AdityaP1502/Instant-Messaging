package model

import (
	"io"

	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
)

type RevokedToken struct {
	TokenID   string `json:"-" db:"token_id"`
	Username  string `json:"-" db:"username"`
	Token     string `json:"token" db:"token"`
	TokenType string `json:"-" db:"type"`
	ExpiredAt string `json:"-" db:"expired_at"`
}

func (t *RevokedToken) FromJSON(r io.Reader, checkRequired bool) error {
	err := util.DecodeJSONBody(r, t)

	if err != nil {
		return err
	}

	if checkRequired {
		return util.CheckParametersUnity(t)
	}

	return nil
}

func (t *RevokedToken) ToJSON(checkRequired bool) ([]byte, error) {
	var err error

	if checkRequired {
		if err = util.CheckParametersUnity(t); err != nil {
			return nil, err
		}
	}

	return util.CreateJSONResponse(t)
}
