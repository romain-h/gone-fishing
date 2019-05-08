package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AuthProvider struct {
	ClientID     string
	ClientSecret string
}

type Splitwise struct {
	AuthProvider
}

type Monzo struct {
	AuthProvider
	AccountID string
}

type Config struct {
	AppURL    string
	RedisURL  string
	Splitwise Splitwise
	Monzo     Monzo
	StartDate string
}

func New() *Config {
	godotenv.Load()

	return &Config{
		AppURL:   os.Getenv("APP_URL"),
		RedisURL: os.Getenv("REDIS_URL"),
		Splitwise: Splitwise{
			AuthProvider: AuthProvider{
				ClientID:     os.Getenv("SPLITWISE_CLIENT_ID"),
				ClientSecret: os.Getenv("SPLITWISE_CLIENT_SECRET"),
			},
		},
		Monzo: Monzo{
			AuthProvider: AuthProvider{
				ClientID:     os.Getenv("MONZO_CLIENT_ID"),
				ClientSecret: os.Getenv("MONZO_CLIENT_SECRET"),
			},
			AccountID: os.Getenv("MONZO_ACCOUNT_ID"),
		},
		StartDate: os.Getenv("START_DATE"),
	}
}
