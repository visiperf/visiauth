package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/visiperf/visiauth/v3"
	"github.com/visiperf/visiauth/v3/errors"
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
		if IsErrRedisNilMessage(err) {
			return nil, errors.NotFound("jwk", "JWK_NOT_FOUND")
		}

		return nil, errors.Internal(err)
	}

	return &jwk, nil
}
