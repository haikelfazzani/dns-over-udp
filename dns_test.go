package main

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
)

func TestLookupIP(t *testing.T) {
	domains := []string{"google.com", "cloudflare.com", "meet.google.com"}

	for _, domain := range domains {
		t.Run(domain, func(t *testing.T) {
			m := new(dns.Msg)
			m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

			c := new(dns.Client)
			in, _, err := c.Exchange(m, "localhost:53")

			fmt.Print("\n", in.Answer)

			if err != nil {
				t.Fatalf("\nFailed to resolve IP for %s: %v", domain, err)
			}

			if len(in.Answer) == 0 {
				t.Fatalf("\nNo IPs resolved for %s", domain)
			}
		})
	}
}
