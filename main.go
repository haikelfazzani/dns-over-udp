package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// ANSI color escape codes
const (
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorReset = "\033[0m"
)

type DNSConfig struct {
	Laddr       string
	Raddr       string
	UseHosts    bool
	UseWildCard bool // *.google.com
}

var config = &DNSConfig{
	Laddr:       ":53",
	Raddr:       "8.8.8.8:53",
	UseHosts:    false,
	UseWildCard: true,
}

func main() {
	server := &dns.Server{Addr: config.Laddr, Net: "udp"}
	dns.HandleFunc(".", handleDNSRequest)
	log.Printf("Server listen on %s", config.Laddr)
	log.Fatal(server.ListenAndServe())
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	for _, question := range r.Question {
		switch question.Qtype {
		case dns.TypeA:
			ip := checkHosts(question.Name)
			if ip != "" {
				m := new(dns.Msg)
				m.SetReply(r)
				m.Answer = append(m.Answer, &dns.A{
					Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
					A:   net.ParseIP(ip),
				})
				w.WriteMsg(m)
			} else {
				queryRemote(w, r)
			}
		case dns.TypeMX:
			// Handle MX record queries
			// Here, you can add your logic to handle MX records
			// For simplicity, let's just respond with NXDOMAIN
			m := new(dns.Msg)
			m.SetReply(r)
			m.SetRcode(r, dns.RcodeNameError)
			w.WriteMsg(m)
		default:
			// Respond with NXDOMAIN for unsupported record types
			m := new(dns.Msg)
			m.SetReply(r)
			m.SetRcode(r, dns.RcodeNameError)
			w.WriteMsg(m)
		}
	}
}

func checkHosts(host string) string {
	// Remove the trailing dot from the host name
	host = strings.TrimSuffix(host, ".")

	// Open the hosts file
	file, err := os.Open("/etc/hosts")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && net.ParseIP(fields[0]) != nil {
			ip := fields[0]
			hosts_domain := fields[1]

			// Handle wildcard matching and equality
			if (config.UseWildCard && strings.HasSuffix(host, "."+hosts_domain) || hosts_domain == "*."+host) || host == hosts_domain {
				return ip
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Domain not found in /etc/hosts: %s%s%s\n", ColorGreen, host, ColorReset)
	return ""
}

func queryRemote(w dns.ResponseWriter, r *dns.Msg) {
	c := new(dns.Client)
	c.Timeout = time.Second * 5 // Set timeout to 5 seconds
	in, _, err := c.Exchange(r, config.Raddr)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Printf("UDP read timeout: %s", err)
			return
		}
		log.Fatal(err)
	}
	w.WriteMsg(in)
}
