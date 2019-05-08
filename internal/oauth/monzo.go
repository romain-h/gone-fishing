package oauth

import (
	"encoding/json"

	"github.com/romain-h/gone-fishing/internal/cache"
	"github.com/romain-h/gone-fishing/internal/config"
	"golang.org/x/oauth2"
)

var EndpointMonzo = oauth2.Endpoint{
	AuthURL:  "https://auth.monzo.com/",
	TokenURL: "https://api.monzo.com/oauth2/token",
}

const endpointMonzoMe = "https://api.monzo.com/ping/whoami"

func NewMonzo(cfg config.Config, cache cache.CacheManager) AuthProvider {
	cacheKey := "monzo_tk"
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.Monzo.AuthProvider.ClientID,
		ClientSecret: cfg.Monzo.AuthProvider.ClientSecret,
		RedirectURL:  cfg.AppURL + "monzo/callback",
		Endpoint:     EndpointMonzo,
	}
	config := Config{
		Config:   oauthCfg,
		cache:    cache,
		cacheKey: cacheKey,
	}
	var tk oauth2.Token

	strToken, _ := cache.Get(cacheKey)
	json.Unmarshal([]byte(strToken), &tk)

	client := config.Client(oauth2.NoContext, &tk)

	return &provider{
		config:     config,
		client:     client,
		endpointMe: endpointMonzoMe,
	}
}
