package util

import (
	"encoding/base64"
	"os"
)

type Config struct {
	ApplicationName string `json:"app_name"`
	Version         string `json:"app_ver"`
	Database        struct {
		Host     string `json:"host"`
		Port     int    `json:"port,string"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"database"`

	Server struct {
		Host   string `json:"host"`
		Port   int    `json:"port,string"`
		Secure string `json:"secure"`
	} `json:"server"`

	Certificate struct {
		CertFile string `json:"certFile"`
		KeyFile  string `json:"KeyFile"`
	} `json:"certificate"`

	Session struct {
		ExpireTime        int    `json:"expireTimeMinutes,string"`
		RefreshExpireTime int    `json:"refreshExpireTimeMinutes,string"`
		SecretKeyBase64   string `json:"secretKey"`
		SecretKeyRaw      []byte `json:"-"`
	} `json:"token"`

	MailAPI struct {
		Host string `json:"host"`
		Port int    `json:"port,string"`
	} `json:"mail"`

	Hash struct {
		SecretKeyBase64 string `json:"secretKey"`
		SecretKeyRaw    []byte `json:"-"`
	} `json:"prehash"`

	OTP struct {
		ResendDurationMinutes int `json:"resendDurationMinutes,string"`
	} `json:"otp"`
}

func ReadJSONConfiguration(path string) (*Config, error) {
	var config Config

	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	err = DecodeJSONBody(configFile, &config)

	if err != nil {
		return nil, err
	}

	// Convert secretkey base64 to raw
	key, err := base64.StdEncoding.DecodeString(config.Session.SecretKeyBase64)

	if err != nil {
		return nil, err
	}

	config.Session.SecretKeyRaw = key
	// Convert secretkey base64 to raw
	key, err = base64.StdEncoding.DecodeString(config.Hash.SecretKeyBase64)

	if err != nil {
		return nil, err
	}

	config.Hash.SecretKeyRaw = key

	return &config, nil
}
