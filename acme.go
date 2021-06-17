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
	Zone    string
}

func NewACME(acmeManagerTemplate certmagic.ACMEManager, zone string) ACME {
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
		Zone:    zone,
	}
}

func (a ACME) OnStartup() error {
	httpPort := fmt.Sprintf(":%d", a.Manager.AltHTTPPort)
	tlsalpnPort := fmt.Sprintf(":%d", a.Manager.AltTLSALPNPort)
	tlsConfig := a.Config.TLSConfig()
	var err error
	if !a.Manager.DisableHTTPChallenge {
		go func() {
			_, err = tls.Listen("tcp", tlsalpnPort, tlsConfig)
		}()
	}
	if !a.Manager.DisableTLSALPNChallenge {
		go func() { err = http.ListenAndServe(httpPort, a.Manager.HTTPChallengeHandler(http.NewServeMux())) }()
	}
	/*go func() {
		dns.HandleFunc(a.Zone, func(w dns.ResponseWriter, r *dns.Msg) {
			state := request.Request{W: w, Req: r}
			var zone string
			if len(r.Question) > 0 {
				zone = r.Question[0].Name
			}
			if checkDNSChallenge(zone) {
				err = solveDNSChallenge(context.Background(), zone, state)
			}
		})
		server := &dns.Server{Addr: ":80", Net: "tcp"}
		err = server.ListenAndServe()
		defer server.Shutdown()
	}()*/
	return err
}

func (a ACME) IssueCert(zones []string) error {
	err := a.Config.ManageSync(zones)
	return err
}

func (a ACME) GetCert(zone string) error {
	err := a.Config.ObtainCertSync(context.Background(), zone)
	return err
}
