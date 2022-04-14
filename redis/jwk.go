package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/visiperf/visiauth/v2"
)

type JwkFetcher struct {
	options *redis.Options
}

func NewJwkFetcher() *JwkFetcher {
	return &JwkFetcher{
		options: &redis.Options{
			Addr:     env.Redis.Addr,
			Username: env.Redis.User,
			Password: env.Redis.Password,
		},
	}
}

func (f *JwkFetcher) FetchJwk(ctx context.Context, kid string) (*visiauth.Jwk, error) {
	client := redis.NewClient(f.options)
	defer client.Close()

	var jwk visiauth.Jwk
	if err := client.Get(ctx, kid).Scan(&jwk); err != nil {
		return nil, err
	}

	return &jwk, nil
}
