package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleConfig() oauth2.Config {
	return oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
        ClientID:     os.Getenv("CLIENT_ID"),
        ClientSecret: os.Getenv("CLIENT_SECRET"),
        Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
        Endpoint:     google.Endpoint,
	}
}