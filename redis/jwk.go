package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/visiperf/visiauth/v2"
)

type JwkFetcher struct {
	client *redis.Client
}

func NewJwkFetcher() *JwkFetcher {
	return &JwkFetcher{
		client: redis.NewClient(&redis.Options{
			Addr: env.Redis.Addr,
		}),
	}
}

func (f *JwkFetcher) FetchJwk(ctx context.Context, kid string) (*visiauth.Jwk, error) {
	var jwk visiauth.Jwk
	if err := f.client.Get(ctx, kid).Scan(&jwk); err != nil {
		return nil, err
	}

	return &jwk, nil
}
