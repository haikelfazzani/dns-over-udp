# DNS forward over UDP

### Build and Run

```shell
go build . && sudo ./main
```

### Configuration

```go
// default config
var config = &DNSConfig{
	HostsFilePath: "/etc/hosts",
	Laddr:         ":53",
	Raddr:         "8.8.8.8:53",
	RfallbackAddr: "8.8.4.4:53",
	UseWildCard:   true,
	Timeout:       5,
}
```

### Todo

- [x] Support wildcards
- [ ] Authentification
- [ ] block by country
- [ ] block domain format IP
- [ ] DNS over HTTPS
