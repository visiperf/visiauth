package visiauth

// Header contains Jwt algorithm informations
type Header struct {
	Alg string
	Typ string
}

// Sub represent the customer which is potentially connected
type Sub struct {
	UserID    int64
	CompanyID int64
}

// Payload contains informations about the customer which is potentially connected
type Payload struct {
	Iat   string
	Exp   string
	Sub   Sub
	Roles []string
}

// Jwt represent access token after base64 decoding
type Jwt struct {
	Header    Header
	Payload   Payload
	Signature string
}
