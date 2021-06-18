package acme

import (
	"sync"
	"testing"

	"github.com/caddyserver/certmagic"
	"github.com/libdns/libdns"
)

func TestACME(t *testing.T) {
	dnsProvider := Provider{
		recordMap: make(map[string][]libdns.Record),
		m:         sync.Mutex{},
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
	a := NewACME(acmeTemplate, zone)
	err := a.OnStartup()
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = a.IssueCert([]string{zone})
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = a.GetCert(zone)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
