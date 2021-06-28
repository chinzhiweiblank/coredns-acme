package acme

import (
	"fmt"

	"github.com/miekg/dns"
)

func getAuthoritativeNameServer(zone string) (string, error) {
	resolvers := recursiveNameservers(nil)
	nameservers, err := lookupNameservers(zone, resolvers)
	if err != nil {
		return "", err
	}
	authoritativeNS := nameservers[len(nameservers)-1]
	return authoritativeNS, nil
}

func getExternalIpAddress(zone string) (string, error) {
	resolvers := recursiveNameservers(nil)
	r, err := dnsQuery(zone, dns.TypeA, resolvers, true)
	if err != nil {
		return "", fmt.Errorf("dns query for zone %v with resolvers %+v error: %+v", zone, resolvers, err)
	}
	for _, rr := range r.Answer {
		if a, ok := rr.(*dns.A); ok {
			return a.A.String(), nil
		}
	}
	return "", fmt.Errorf("no A record found for zone %s", zone)
}
