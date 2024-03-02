package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// ANSI color escape codes
const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\x1b[34m"
	ColorReset  = "\033[0m"
)

type DNSConfig struct {
	HostsFilePath string
	Laddr         string
	Raddr         string
	UseWildCard   bool // *.google.com
	Timeout       uint8
}

var config = &DNSConfig{
	HostsFilePath: "/etc/hosts",
	Laddr:         ":53",
	Raddr:         "8.8.8.8:53",
	UseWildCard:   true,
	Timeout:       15,
}

func main() {
	StartServer()
}

func StartServer() {
	server := &dns.Server{Addr: config.Laddr, Net: "udp"}
	dns.HandleFunc(".", handleDNSRequest)
	log.Printf("\n%sDNS Proxy Server listen on %s%s\n", ColorBlue, config.Laddr, ColorReset)
	log.Fatal(server.ListenAndServe())
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	for _, question := range r.Question {
		switch question.Qtype {
		case dns.TypeA:
			ip := checkHostsFile(question.Name)
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

func checkHostsFile(host string) string {
	host = strings.TrimSuffix(host, ".")
	wildcardPattern := regexp.MustCompile(`^\s*\*\s*\.(.+)$`)

	file, err := os.Open(config.HostsFilePath)
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

			if config.UseWildCard && wildcardPattern.MatchString(hosts_domain) {

				domain := wildcardPattern.FindStringSubmatch(hosts_domain)[1]
				if strings.HasSuffix(host, "."+domain) {
					fmt.Printf("%s[%s] Resolve From Hosts File: %s (%s) %s\n", ColorYellow, hosts_domain, host, ip, ColorReset)
					return ip
				}
			} else if host == hosts_domain {
				fmt.Printf("%s[%s] Resolve From Hosts File: %s (%s) %s\n", ColorYellow, hosts_domain, host, ip, ColorReset)
				return ip
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s[] Resolve From %s: %s%s\n", ColorGreen, config.Raddr, host, ColorReset)
	return ""
}

func queryRemote(w dns.ResponseWriter, r *dns.Msg) {
	c := new(dns.Client)
	c.Timeout = time.Second * time.Duration(config.Timeout)

	in, _, err := c.Exchange(r, config.Raddr)

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Printf("> UDP read timeout: %s", err)
			log.Printf("%sRetrying query %s %s", ColorBlue, config.Raddr, ColorReset)
			w.WriteMsg(in)
			return
		}
		log.Fatal(err)

	}

	w.WriteMsg(in)
}
