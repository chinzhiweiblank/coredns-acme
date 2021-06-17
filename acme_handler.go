package acme

import (
	"context"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type AcmeHandler struct {
	Next plugin.Handler
}

func (h AcmeHandler) Name() string { return pluginName }

func (h AcmeHandler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	var zone string
	if len(r.Question) > 0 {
		zone = r.Question[0].Name
	}
	if checkDNSChallenge(zone) {
		err := solveDNSChallenge(ctx, zone, state)
		if err != nil {
			return 0, err
		}
		err = configureTLS(A, Config)
		if err != nil {
			return 0, err
		}
	}
	return h.Next.ServeDNS(ctx, w, r)
}

func checkDNSChallenge(zone string) bool {
	dnsChallengeString := "_acme_challenge"
	return strings.HasPrefix(zone, dnsChallengeString)
}

func solveDNSChallenge(ctx context.Context, zone string, state request.Request) error {
	a := new(dns.Msg)
	a.SetReply(state.Req)
	a.Authoritative = true
	records, err := provider.GetRecords(ctx, zone)
	if err != nil {
		return err
	}
	rrs := []dns.RR{}
	for _, record := range records {
		rr := new(dns.TXT)
		rr.Txt = []string{record.Value}
		rr.Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeTXT, Class: state.QClass()}
		rrs = append(rrs, rr)
	}
	a.Answer = rrs
	state.W.WriteMsg(a)
	return nil
}
