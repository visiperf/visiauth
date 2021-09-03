package visiauth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapChainToCertificateConverter(t *testing.T) {
	tests := []struct {
		name        string
		converter   *WrapChainToCertificateConverter
		chain       string
		certificate string
	}{{
		name:      "success",
		converter: NewWrapChainToCertificateConverter(),
		chain:     `azerty`,
		certificate: `-----BEGIN CERTIFICATE-----
		azerty
		-----END CERTIFICATE-----`,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, strings.ReplaceAll(test.certificate, "\t", ""), test.converter.ConvertChainToCertificate(test.chain))
		})
	}
}
