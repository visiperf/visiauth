# Authentication service for Golang

Package `visiauth` provide functions to help you making authentication verifications. Decode an Access token to get a JWT more easily with this amazing package ! :)

Table of contents
=================

  * [Install](#install)
  * [Usage](#usage)
      * [Access token to JWT](#access-token-to-jwt)

## Install

Use `go get` to install this package.

    go get github.com/visiperf/visiauth
    
## Usage

### Access token to JWT

The DecodeAccessToken(token string) (\*Jwt, error) function is used to decode an Access token, validate it and check expiration.

```go

package main

import "github.com/visiperf/visiauth/visiperf"

func main() {
	jwt, err := visiperf.NewAuthService("YOUR_SECRET").DecodeAccessToken("YOUR_TOKEN")
}

```