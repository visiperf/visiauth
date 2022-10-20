package api

import (
	"net/http"
	"net/url"

	"github.com/visiperf/visiauth/v3/errors"
	"github.com/visiperf/visiauth/v3/neo4j"
	"github.com/visiperf/visiauth/v3/redis"
	"github.com/visiperf/visiauth/v3/renderer"
	"github.com/visiperf/visiauth/v3/visiauth"
)

const (
	accessTokenQueryParamsKey = "token"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	authenticable, err := func() (visiauth.Authenticable, error) {
		vs, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			return nil, errors.InvalidArgument(err.Error(), "INVALID_QUERY_PARAMS")
		}

		if len(vs.Get(accessTokenQueryParamsKey)) <= 0 {
			return nil, errors.InvalidArgument("token in query params is required", "TOKEN_QUERY_PARAMS_REQUIRED")
		}

		return visiauth.NewService(redis.NewJwkFetcher(), neo4j.NewUserRepository()).DecodeAccessToken(r.Context(), vs.Get(accessTokenQueryParamsKey))
	}()
	if err != nil {
		renderer.Error(err, w)
		return
	}

	renderer.Success(authenticable, http.StatusOK, w)
}
