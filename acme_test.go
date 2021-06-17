package acme

import (
	"testing"

	"github.com/caddyserver/certmagic"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/libdns/libdns"
)

func TestACME(t *testing.T) {
	dnsProvider := Provider{
		recordMap: make(map[string][]libdns.Record),
	}

	zone := "daedric.online"
	acmeTemplate := certmagic.ACMEManager{
		CA:                      certmagic.LetsEncryptStagingCA,
		Agreed:                  true,
		AltHTTPPort:             8089,
		AltTLSALPNPort:          8090,
		DisableHTTPChallenge:    false,
		DisableTLSALPNChallenge: false,
		DNS01Solver: &certmagic.DNS01Solver{
			DNSProvider: &dnsProvider,
		},
	}
	dnsConfig := dnsserver.Config{
		Zone: zone,
	}
	a := NewACME(acmeTemplate, zone)
	err := configureTLS(a, &dnsConfig)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(dnsConfig.TLSConfig.Certificates) == 0 {
		t.Errorf("Certificates were not configured for TLS")
	}
}
