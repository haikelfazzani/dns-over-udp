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
	ListRaddr     []string
	UseWildCard   bool // *.google.com
	Timeout       uint8
}

var config = &DNSConfig{
	HostsFilePath: "/etc/hosts",
	Laddr:         ":5053",
	ListRaddr:     []string{"8.8.8.8:53", "8.8.4.4:53", "1.1.1.1:53", "9.9.9.9:53"},
	UseWildCard:   true,
	Timeout:       15,
}

func main() {
	server := &dns.Server{Addr: config.Laddr, Net: "udp"}
	dns.HandleFunc(".", handleDNSRequest)
	log.Printf("\n%sDNS Server listen on %s%s\n", ColorBlue, config.Laddr, ColorReset)
	log.Fatal(server.ListenAndServe())
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	for _, question := range r.Question {
		ip := checkHostsFile(question.Name)
		var rr dns.RR
		var rrtype uint16

		switch question.Qtype {
		case dns.TypeA:
			rrtype = dns.TypeA
		case dns.TypeAAAA:
			rrtype = dns.TypeAAAA
		default:
			m := new(dns.Msg)
			m.SetReply(r)
			m.SetRcode(r, dns.RcodeNameError)
			w.WriteMsg(m)
			continue
		}

		if ip != "" {
			if rrtype == dns.TypeA {
				rr = &dns.A{
					Hdr: dns.RR_Header{Name: question.Name, Rrtype: rrtype, Class: dns.ClassINET, Ttl: 0},
					A:   net.ParseIP(ip),
				}
			} else {
				rr = &dns.AAAA{
					Hdr:  dns.RR_Header{Name: question.Name, Rrtype: rrtype, Class: dns.ClassINET, Ttl: 0},
					AAAA: net.ParseIP(ip),
				}
			}
		} else {
			queryRemote(w, r)
		}

		m := new(dns.Msg)
		m.SetReply(r)
		m.Answer = append(m.Answer, rr)
		w.WriteMsg(m)
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
					fmt.Printf("%s[%s] Resolve From Hosts File: %s (%s) %s\n", ColorBlue, hosts_domain, host, ip, ColorReset)
					return ip
				}
			} else if host == hosts_domain {
				fmt.Printf("%s[hosts file] %s (%s) %s\n", ColorBlue, host, ip, ColorReset)
				return ip
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ""
}

func queryRemote(w dns.ResponseWriter, r *dns.Msg) {
	c := new(dns.Client)
	c.Timeout = time.Second * time.Duration(config.Timeout)

	for _, dnsServer := range config.ListRaddr {
		in, _, err := c.Exchange(r, dnsServer)
		if err == nil {
			fmt.Printf("\n> %s[%s] %s\n(%s)%s\n", ColorYellow, dnsServer, &r.Question[0], in.Answer[0], ColorReset)
			w.WriteMsg(in)
			return
		}
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Printf("> UDP read timeout: %s", err)
			log.Printf("%sRetrying query %s %s", ColorRed, dnsServer, ColorReset)
			continue
		}
		log.Println(err)
	}
}
