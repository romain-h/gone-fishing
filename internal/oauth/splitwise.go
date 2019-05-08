package oauth

import (
	"encoding/json"

	"github.com/romain-h/gone-fishing/internal/cache"
	"github.com/romain-h/gone-fishing/internal/config"
	"golang.org/x/oauth2"
)

var EndpointSplitwise = oauth2.Endpoint{
	AuthURL:  "https://secure.splitwise.com/oauth/authorize",
	TokenURL: "https://secure.splitwise.com/oauth/token",
}

const endpointSplitwiseMe = "https://www.splitwise.com/api/v3.0/get_current_user"

func NewSplitwise(cfg config.Config, cache cache.CacheManager) AuthProvider {
	cacheKey := "splitwise_tk"
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.Splitwise.AuthProvider.ClientID,
		ClientSecret: cfg.Splitwise.AuthProvider.ClientSecret,
		RedirectURL:  cfg.AppURL + "splitwise/callback",
		Endpoint:     EndpointSplitwise,
	}
	config := Config{
		Config:   oauthCfg,
		cache:    cache,
		cacheKey: cacheKey,
	}
	var tk oauth2.Token

	// When initialising try to read token from cache
	str_token, _ := cache.Get(cacheKey)
	json.Unmarshal([]byte(str_token), &tk)

	client := config.Client(oauth2.NoContext, &tk)

	return &provider{
		config:     config,
		client:     client,
		endpointMe: endpointSplitwiseMe,
	}
}
