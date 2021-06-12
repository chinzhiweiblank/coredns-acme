package acme

import (
	"reflect"
	"testing"

	"github.com/caddyserver/certmagic"
	"github.com/coredns/caddy"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		shouldErr    bool
		acmeTemplate certmagic.ACMEManager
	}{
		{
			"Correct Config with correct challenge",
			`acme example.domain {
				http01 89
				tlsalpn 90
			}`,
			false,
			certmagic.ACMEManager{
				DisableHTTPChallenge:    false,
				DisableTLSALPNChallenge: false,
				AltHTTPPort:             89,
				AltTLSALPNPort:          90,
			},
		},
		{
			"Invalid challenge",
			`acme example.domain {
				invalid_challenge
			`,
			true,
			certmagic.ACMEManager{},
		},
		{
			"Missing domain argument",
			`acme {
				dns01
			}`,
			true,
			certmagic.ACMEManager{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := caddy.NewTestController("acme", test.input)
			if acmeTemplate, err := parseACME(c); (err != nil) != test.shouldErr {
				if err != nil && reflect.DeepEqual(test.acmeTemplate, acmeTemplate) {
					t.Errorf("Error: AcmeTemplate %+v is not configured as it should be %+v", acmeTemplate, test.acmeTemplate)
				}
				t.Errorf("Error: setup() error = %v, shouldErr %v", err, test.shouldErr)
			}
		})
	}
}
