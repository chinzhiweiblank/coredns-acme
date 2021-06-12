package acme

import (
	"crypto/tls"
	"strconv"
	"strings"

	"github.com/caddyserver/certmagic"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

const pluginName = "acme"

func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	acmeTemplate, err := parseACME(c)
	if err != nil {
		return plugin.Error("acme", err)
	}
	config := dnsserver.GetConfig(c)

	a := NewACME(acmeTemplate)
	err = configureTLS(a, config)
	if err != nil {
		return c.Errf("Unexpected error: %s", err.Error())
	}
	return nil
}

func setTLSDefaults(tlsConfig *tls.Config) {
	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.MaxVersion = tls.VersionTLS13
	tlsConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
	tlsConfig.PreferServerCipherSuites = true
}

func parseACME(c *caddy.Controller) (certmagic.ACMEManager, error) {
	acmeTemplate := certmagic.ACMEManager{
		DisableHTTPChallenge:    true,
		DisableTLSALPNChallenge: true,
	}
	var err error
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) != 1 {
			return acmeTemplate, c.Errf("Unexpected number of arguments: %#v", args)
		}
		for c.NextBlock() {
			challenge := strings.ToLower(c.Val())
			switch challenge {
			case HTTPChallenge:
				args := c.RemainingArgs()
				if len(args) > 1 {
					return acmeTemplate, c.Errf("Unexpected number of arguments: %#v", args)
				}
				httpPort := 80
				if len(args) == 1 {
					httpPort, err = strconv.Atoi(args[0])
					if err != nil {
						return acmeTemplate, c.Errf("HTTP port is not a string: %#v", args)
					}
				}
				acmeTemplate.AltHTTPPort = httpPort
				acmeTemplate.DisableHTTPChallenge = false
			case TLPSALPNChallenge:
				args := c.RemainingArgs()
				if len(args) > 1 {
					return acmeTemplate, c.Errf("Unexpected number of arguments: %#v", args)
				}
				var err error
				tlsAlpnPort := 443
				if len(args) == 1 {
					tlsAlpnPort, err = strconv.Atoi(args[0])
					if err != nil {
						return acmeTemplate, c.Errf("TlsAlpn port is not a string: %#v", args)
					}
				}
				acmeTemplate.AltTLSALPNPort = tlsAlpnPort
				acmeTemplate.DisableTLSALPNChallenge = false
			default:
				return acmeTemplate, c.Errf("Unexpected challenge: %s", challenge)
			}
		}
	}
	return acmeTemplate, nil
}

func configureTLS(a ACME, conf *dnsserver.Config) error {
	err := a.OnStartup()
	if err != nil {
		return err
	}
	zone := conf.Zone
	err = a.IssueCert([]string{zone})
	if err != nil {
		return err
	}
	err = a.GetCert(zone)
	if err != nil {
		return err
	}
	cert, err := a.Config.CacheManagedCertificate(zone)
	if err != nil {
		return err
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert.Certificate}}
	tlsConfig.ClientAuth = tls.NoClientCert
	tlsConfig.ClientCAs = tlsConfig.RootCAs

	setTLSDefaults(tlsConfig)

	conf.TLSConfig = tlsConfig
	return nil
}
