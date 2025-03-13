package main

import (
	"net/http"

	"github.com/EBOOST-LTD/eboost-cms-partner-BE/app/interface/api/middleware"
)

const (
	herokuEnv = "heroku"
	devEnv    = "development"
)

func (s *service) initCookieAuthenticator(env string) {
	sameSite := http.SameSiteLaxMode
	if env == herokuEnv {
		// This configuration is for heroku review app
		// Cookie is not sent with cross domain unless setting None
		sameSite = http.SameSiteNoneMode
	}

	var secure bool
	if env != devEnv {
		secure = true
	}

	middleware.InitCookieOptions(
		sameSite,
		secure,
	)
}
