package oauth

import (
	"net/http"

	"golang.org/x/oauth2"
)

type provider struct {
	config     Config
	client     *http.Client
	endpointMe string
}

func (p *provider) GetClient() *http.Client {
	return p.client
}

func (p *provider) GetRedirectURL() string {
	return p.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (p *provider) Callback(code string) error {
	tk, _ := p.config.Exchange(oauth2.NoContext, code)

	p.client = p.config.Client(oauth2.NoContext, tk)
	_, err := p.client.Get(p.endpointMe)
	return err
}
