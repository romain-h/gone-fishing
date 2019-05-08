package cache

import (
	"net/url"

	"github.com/go-redis/redis"
	"github.com/romain-h/gone-fishing/internal/config"
)

type CacheManager interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	GetByte(key string) ([]byte, error)
	SetByte(key string, value []byte) error
}

type manager struct {
	client *redis.Client
}

func New(cfg config.Config) CacheManager {
	redisURL := cfg.RedisURL
	parsedURL, _ := url.Parse(redisURL)
	password, _ := parsedURL.User.Password()
	resolvedURL := parsedURL.Host
	client := redis.NewClient(&redis.Options{
		Addr:     resolvedURL,
		Password: password,
		DB:       0, // use default DB
	})

	return &manager{client: client}
}

func (mgr *manager) Get(key string) (res string, err error) {
	return mgr.client.Get(key).Result()
}

func (mgr *manager) Set(key string, value string) (err error) {
	return mgr.client.Set(key, value, 0).Err()
}

func (mgr *manager) GetByte(key string) (res []byte, err error) {
	return mgr.client.Get(key).Bytes()
}

func (mgr *manager) SetByte(key string, value []byte) (err error) {
	return mgr.client.Set(key, value, 0).Err()
}
