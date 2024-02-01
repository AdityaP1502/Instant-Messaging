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
		ExpireTime      int    `json:"expireTime,string"`
		SecretKeyBase64 string `json:"secretKey"`
		secretKeyRaw    []byte `json:"-"`
	} `json:"Session"`

	SMTPConfig struct {
		Host       string `json:"host"`
		Port       int    `json:"port,string"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		SenderName string `json:"sender_name"`
	}
}

func ReadJSONConfiguration(path string) (*Config, error) {
	var config Config

	configFile, err := os.Open(path)
	defer configFile.Close()

	if err != nil {
		return nil, err
	}

	err = DecodeJSONBody(configFile, &config)

	if err != nil {
		return nil, err
	}

	// Convert secretkey base64 to raw
	key, err := base64.StdEncoding.DecodeString(config.Session.SecretKeyBase64)

	if err != nil {
		return nil, err
	}

	config.Session.secretKeyRaw = key

	return &config, nil
}
