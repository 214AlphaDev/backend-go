package config

import (
	"encoding/base64"
	vo "github.com/214alphadev/community-bl/value_objects"
)

type Configuration struct {
	databaseUrl           string
	accessTokenSigningKey vo.AccessTokenSigningKey
	port                  string
	sendGridApiKey        string
	fromMail              string
	applicationMode       string
}

func (c Configuration) DatabaseURL() string {
	return c.databaseUrl
}

func (c Configuration) AccessTokenSigningKey() vo.AccessTokenSigningKey {
	return c.accessTokenSigningKey
}

func (c Configuration) Port() string {
	return c.port
}

func (c Configuration) SendGridApiKey() string {
	return c.sendGridApiKey
}

func (c Configuration) FromMail() (vo.EmailAddress, error) {
	return vo.NewEmailAddress(c.fromMail)
}

func (c Configuration) Development() bool {
	return c.applicationMode == "Development"
}

func NewConfiguration(getEnv func(key string) string) (Configuration, error) {

	accessTokenSigningKeyBytes, err := base64.URLEncoding.DecodeString(getEnv("ACCESS_TOKEN_SIGNING_KEY"))
	if err != nil {
		return Configuration{}, err
	}
	accessTokenSigningKey, err := vo.NewAccessTokenSigningKey(accessTokenSigningKeyBytes)
	if err != nil {
		return Configuration{}, err
	}

	return Configuration{
		databaseUrl:           getEnv("DATABASE_URL"),
		accessTokenSigningKey: accessTokenSigningKey,
		port:                  getEnv("PORT"),
		sendGridApiKey:        getEnv("SENDGRID_API_KEY"),
		fromMail:              getEnv("SEND_EMAIL_FROM"),
		applicationMode:       getEnv("MODE"),
	}, nil
}
