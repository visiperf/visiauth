## Project goals

This package provides functions and service to help you making authentication verifications. You can retrieve and decode bearer tokens, such as User access token or M2M token.

## Usage

### Retrieve token

#### gRPC Context

```go
// context metadata => Authorization: Bearer ...

token, err := visiauth.RetrieveTokenFromContext(ctx)
if err != nil {
  ...
}
```

#### HTTP Request

```go
// request header => Authorization: Bearer ...

token, err := visiauth.RetrieveTokenFromRequest(req)
if err != nil {
  ...
}
```

#### PubSub message attribute

```go
// message attribute => Authorization: Bearer ...

token, err := visiauth.RetrieveTokenFromPubSubMessageAttribute(req)
if err != nil {
  ...
}
```

### Decode token

```go
service := visiauth.NewService(redis.NewJwkFetcher(), neo4j.NewUserRepository())

identity, err = service.DecodeAccessToken(r.Context(), token)
if err != nil {
  ...
}
```

## Identities

Two types of identities are provided by service : User and Application.

### User access token (User)

If token was generated using Client credentials flow, identity will be a user.

```go
user, ok := identity.(*visiauth.User)
if !ok {
  ...
}
```

### M2M Token (Application)

If token was generated using Client ID and Client secret, identity will be an application.

```go
app, ok := identity.(*visiauth.App)
if !ok {
  ...
}
```

