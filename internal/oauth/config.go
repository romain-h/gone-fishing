package oauth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/romain-h/gone-fishing/internal/cache"
	"golang.org/x/oauth2"
)

type Config struct {
	*oauth2.Config
	cache    cache.CacheManager
	cacheKey string
}

func (c *Config) StoreToken(token *oauth2.Token) error {
	t, _ := json.Marshal(token)
	err := c.cache.Set(c.cacheKey, string(t))
	return err
}

func (c *Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := c.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	if err := c.StoreToken(token); err != nil {
		return nil, err
	}
	return token, nil
}

func (c *Config) Client(ctx context.Context, t *oauth2.Token) *http.Client {
	return oauth2.NewClient(ctx, c.TokenSource(ctx, t))
}

func (c *Config) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	rts := &CacheTokenSource{
		source: c.Config.TokenSource(ctx, t),
		config: c,
	}
	return oauth2.ReuseTokenSource(t, rts)
}

type CacheTokenSource struct {
	source oauth2.TokenSource
	config *Config
}

func (t *CacheTokenSource) Token() (*oauth2.Token, error) {
	token, err := t.source.Token()
	if err != nil {
		return nil, err
	}
	if err := t.config.StoreToken(token); err != nil {
		return nil, err
	}
	return token, nil
}
