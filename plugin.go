package main

import (
    "io/ioutil"
    "net"
    "strings"
    "time"

    "github.com/coredns/coredns/plugin"
    "github.com/coredns/coredns/plugin/pkg/dnsutil"
    "github.com/miekg/dns"
)

type DynamicDomain struct {
    Next       plugin.Handler
    domains    map[string]bool
    updateTime time.Time
    forward1   string
    forward2   string
}

func (d *DynamicDomain) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
    if r.MsgHdr.Rcode != dns.RcodeSuccess {
        return dns.RcodeServerFailure, nil
    }

    query := dnsutil.TrimZone(r.Question[0].Name, d.DomainsZone)

    if _, ok := d.domains[query]; ok {
        // Forward to specific DNS server
        c := &dns.Client{}
        m := new(dns.Msg)
        m.SetQuestion(r.Question[0].Name, r.Question[0].Qtype)
        response, _, err := c.Exchange(m, d.forward1)
        if err != nil {
            return dns.RcodeServerFailure, err
        }
        w.WriteMsg(response)
        return dns.RcodeSuccess, nil
    }

    // Forward to default DNS servers
    c := &dns.Client{}
    m := new(dns.Msg)
    m.SetQuestion(r.Question[0].Name, r.Question[0].Qtype)
    response, _, err := c.Exchange(m, d.forward2)
    if err != nil {
        return dns.RcodeServerFailure, err
    }
    w.WriteMsg(response)
    return dns.RcodeSuccess, nil
}

func (d *DynamicDomain) Name() string {
    return "dynamicdomain"
}

func (d *DynamicDomain) UpdateDomains() {
    data, err := ioutil.ReadFile("domain.txt")
    if err != nil {
        log.Printf("Error reading domain file: %v", err)
        return
    }

    domains := strings.Split(string(data), "\n")
    newDomains := make(map[string]bool)
    for _, domain := range domains {
        domain = strings.TrimSpace(domain)
        if domain != "" {
            newDomains[domain] = true
        }
    }
    d.domains = newDomains
    d.updateTime = time.Now()
}

func (d *DynamicDomain) Start() {
    go func() {
        for {
            d.UpdateDomains()
            time.Sleep(5 * time.Minute) // Update every 5 minutes
        }
    }()
}
