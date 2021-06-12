package acme

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/caddyserver/certmagic"
)

const (
	HTTPChallenge     = "http01"
	DNSChallenge      = "dns01"
	TLPSALPNChallenge = "tlsalpn"
)

type ACME struct {
	Manager *certmagic.ACMEManager
	Config  *certmagic.Config
}

type Config struct {
	Token  string
	User   string
	Server string
	Domain string
}

func NewACME(acmeManagerTemplate certmagic.ACMEManager) ACME {
	configTemplate := certmagic.NewDefault()
	cache := certmagic.NewCache(certmagic.CacheOptions{
		GetConfigForCert: func(cert certmagic.Certificate) (*certmagic.Config, error) {
			return configTemplate, nil
		},
	})
	config := certmagic.New(cache, *configTemplate)
	acmeManager := certmagic.NewACMEManager(config, acmeManagerTemplate)
	config.Issuers = append(config.Issuers, acmeManager)
	return ACME{
		Config:  config,
		Manager: acmeManager,
	}
}

func (a ACME) OnStartup() error {
	httpPort := fmt.Sprintf(":%d", a.Manager.AltHTTPPort)
	tlsalpnPort := fmt.Sprintf(":%d", a.Manager.AltTLSALPNPort)
	tlsConfig := a.Config.TLSConfig()
	var err error
	go func() {
		_, err = tls.Listen("tcp", tlsalpnPort, tlsConfig)
	}()
	go func() { err = http.ListenAndServe(httpPort, a.Manager.HTTPChallengeHandler(http.NewServeMux())) }()
	return err
}

func (a ACME) IssueCert(domains []string) error {
	err := a.Config.ManageSync(domains)
	return err
}

func (a ACME) GetCert(domain string) error {
	err := a.Config.ObtainCert(context.Background(), domain, false)
	return err
}

func (a ACME) RevokeCert(domains []string) error {
	for _, domain := range domains {
		ctx := context.Background()
		err := a.Config.RevokeCert(ctx, domain, 0, false)
		if err != nil {
			return err
		}
	}
	return nil
}
