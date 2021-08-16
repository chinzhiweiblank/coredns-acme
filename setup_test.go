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
		zoneName     string
	}{
		{
			"Correct Config with only DNS challenge",
			`acme {
				domain test.domain
			}`,
			false,
			certmagic.ACMEManager{
				DisableHTTPChallenge:    true,
				DisableTLSALPNChallenge: true,
				AltHTTPPort:             80,
				AltTLSALPNPort:          443,
			},
			"test.domain",
		},
		{
			"Correct Config with correct challenge",
			`acme {
				domain test.domain
				challenge http port 89
				challenge tlsalpn port 8080
			}`,
			false,
			certmagic.ACMEManager{
				DisableHTTPChallenge:    false,
				DisableTLSALPNChallenge: false,
				AltHTTPPort:             89,
				AltTLSALPNPort:          90,
			},
			"test.domain",
		},
		{
			"Correct Config with default http01 port",
			`acme {
				domain test.domain
				challenge tlsalpn port 90
			}`,
			false,
			certmagic.ACMEManager{
				DisableHTTPChallenge:    false,
				DisableTLSALPNChallenge: false,
				AltHTTPPort:             80,
				AltTLSALPNPort:          90,
			},
			"test.domain",
		},
		{
			"Missing domain",
			`acme {
				http01 hello
			}`,
			true,
			certmagic.ACMEManager{},
			"",
		},
		{
			"Invalid port",
			`acme {
				domain test.domain
				challenge http port hello
			}`,
			true,
			certmagic.ACMEManager{},
			"test.domain",
		},
		{
			"Invalid challenge",
			`acme {
				domain test.domain
				invalid_challenge
			`,
			true,
			certmagic.ACMEManager{},
			"test.domain",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := caddy.NewTestController("acme", test.input)
			acmeTemplate, zoneName, err := parseACME(c)
			if (err != nil) != test.shouldErr {
				t.Errorf("Error: setup() error = %v, shouldErr %v", err, test.shouldErr)
				if !test.shouldErr && err != nil && reflect.DeepEqual(test.acmeTemplate, acmeTemplate) && test.zoneName == zoneName {
					t.Errorf("Error: AcmeTemplate %+v Zone %s is not configured as it should be %+v", acmeTemplate, zoneName, test.acmeTemplate)
				}
			}
		})
	}
}
