package acme

import (
	"context"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type AcmeHandler struct {
	Next     plugin.Handler
	provider *Provider
	*AcmeConfig
}

type AcmeConfig struct {
	Zone                    string
	IpAddr                  net.IP
	AuthoritativeNameServer string
}

const dnsChallengeString = "_acme-challenge."
const CertificateAuthority = "letsencrypt.org."

func (h AcmeHandler) Name() string { return pluginName }

func (h AcmeHandler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	a := new(dns.Msg)
	a.SetReply(state.Req)
	a.Answer = []dns.RR{}
	for _, question := range r.Question {
		zone := strings.ToLower(question.Name)
		if checkDNSChallenge(zone) {
			if question.Qtype == dns.TypeSOA {
				rr := new(dns.SOA)
				rr.Ns = h.AuthoritativeNameServer
				rr.Mbox = CertificateAuthority
				rr.Hdr = dns.RR_Header{Name: question.Name, Rrtype: dns.TypeSOA, Class: state.QClass()}
				rr.Serial = uint32(1)
				rr.Expire = uint32(60)
				rr.Minttl = uint32(60)
				a.Answer = append(a.Answer, rr)
			} else if question.Qtype == dns.TypeTXT {
				err := h.solveDNSChallenge(ctx, zone, state, a)
				if err != nil {
					return 0, err
				}
			} else if question.Qtype == dns.TypeNS {
				rr := new(dns.NS)
				rr.Ns = h.AuthoritativeNameServer
				rr.Hdr = dns.RR_Header{Name: question.Name, Rrtype: dns.TypeNS, Class: state.QClass()}
				a.Answer = append(a.Answer, rr)
			} else if question.Qtype == dns.TypeA {
				rr := new(dns.A)
				rr.A = h.IpAddr
				rr.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeA, Class: state.QClass()}
				a.Answer = append(a.Answer, rr)
			}
		}
		if zone == h.Zone+"." {
			if question.Qtype == dns.TypeSOA {
				rr := new(dns.SOA)
				rr.Ns = h.AuthoritativeNameServer
				rr.Mbox = CertificateAuthority
				rr.Hdr = dns.RR_Header{Name: question.Name, Rrtype: dns.TypeSOA, Class: state.QClass()}
				rr.Serial = uint32(1)
				rr.Expire = uint32(60)
				rr.Minttl = uint32(60)
				a.Answer = append(a.Answer, rr)
			} else if question.Qtype == dns.TypeCAA {
				rr := new(dns.CAA)
				rr.Tag = "issue"
				rr.Value = "letsencrypt.org"
				rr.Hdr = dns.RR_Header{Name: question.Name, Rrtype: dns.TypeCAA, Class: state.QClass()}
				a.Answer = append(a.Answer, rr)
			} else if question.Qtype == dns.TypeA {
				rr := new(dns.A)
				rr.A = h.IpAddr
				rr.Hdr = dns.RR_Header{Name: question.Name, Rrtype: dns.TypeA, Class: state.QClass()}
				a.Answer = append(a.Answer, rr)
			}
		}
	}
	if len(a.Answer) != 0 {
		err := w.WriteMsg(a)
		if err != nil {
			log.Error("acmeHandler.ServeDNS w.WriteMsg error: ", err)
			return 0, err
		}
	}
	return h.Next.ServeDNS(ctx, w, r)
}

func checkDNSChallenge(zone string) bool {
	return strings.HasPrefix(zone, dnsChallengeString)
}

func (h *AcmeHandler) solveDNSChallenge(ctx context.Context, zone string, state request.Request, a *dns.Msg) error {
	a.Authoritative = true
	records, err := h.provider.GetRecords(ctx, zone)
	if err != nil {
		return err
	}
	rrs := []dns.RR{}
	for _, record := range records {
		rr := new(dns.TXT)
		rr.Txt = []string{record.Value}
		rr.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeTXT, Class: state.QClass(), Ttl: uint32(record.TTL)}
		rrs = append(rrs, rr)
	}
	a.Answer = append(a.Answer, rrs...)
	return err
}
