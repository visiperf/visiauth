package visiperf

type jwt struct {
	Header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	Payload struct {
		Iat string `json:"iat"`
		Exp string `json:"exp"`
		Sub struct {
			UserID    int64 `json:"userId"`
			CompanyID int64 `json:"groupId"`
		} `json:"sub"`
		Roles []string `json:"roles"`
	}
	Signature string
}

func (jwt *jwt) isValid(secret string) error {
	return nil
}

func (jwt *jwt) isExpired() error {
	return nil
}

func (jwt *jwt) isUnlimited() bool {
	return true
}

func (jwt *jwt) generateSignature() (string, error) {
	return "", nil
}

func (jwt *jwt) toString() (string, error) {
	return "", nil
}
