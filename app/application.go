package app

import (
	"github.com/214alphadev/community-bl"
	"net/http"
	"s33d-backend/config"
)

type IApplication interface {
	Start() error
}

type Application struct {
	config    config.Configuration
	Community community_bl.CommunityInterface
}

func (a *Application) Start() error {
	return http.ListenAndServe(":"+a.config.Port(), nil)
}
