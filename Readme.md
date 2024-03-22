# DNS forward over UDP

🌐 DNS over UDP: receiving a request from a client and then forwarding that
request to DNS Resolver to obtain the answer.

### Build and Run

```shell
go build main.go && sudo ./main
```

### Configuration

```go
// default config
var config = &DNSConfig{
	HostsFilePath: "/etc/hosts",
	Laddr:         ":5053",
	ListRaddr:     []string{"8.8.8.8:53", "8.8.4.4:53", "1.1.1.1:53", "9.9.9.9:53"},
	UseWildCard:   true,
	Timeout:       15,
}
```

### Test
```shell
go test
```

### License
Apache 3.0