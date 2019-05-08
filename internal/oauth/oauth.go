package oauth

import (
	"net/http"
)

type AuthProvider interface {
	GetClient() *http.Client
	GetRedirectURL() string
	Callback(code string) error
}
