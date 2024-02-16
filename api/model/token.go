package model

import (
	"io"

	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
)

type Token struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (t *Token) FromJSON(r io.Reader, checkRequired bool, requiredFields []string) error {
	err := util.DecodeJSONBody(r, t)

	if err != nil {
		return err
	}

	if checkRequired {
		return util.CheckParametersUnity(t, requiredFields)
	}

	return nil
}

func (t *Token) ToJSON(checkRequired bool, requiredFields []string) ([]byte, error) {
	var err error

	if checkRequired {
		if err = util.CheckParametersUnity(t, requiredFields); err != nil {
			return nil, err
		}
	}

	return util.CreateJSONResponse(t)
}
