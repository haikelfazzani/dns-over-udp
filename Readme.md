# DNS forward over UDP
🌐 DNS over UDP: receiving a request from a client and then forwarding that request to DNS Resolver to obtain the answer.

### Build and Run

```shell
go build main.go && sudo ./main
```

### Configuration

```go
// default config
var config = &DNSConfig{
	HostsFilePath: "/etc/hosts",
	Laddr:         ":53",
	ListRaddr:     "8.8.8.8:53",
	UseWildCard:   true,
	Timeout:       15,
}
```
