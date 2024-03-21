# DNS forward over UDP

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
	Raddr:         "8.8.8.8:53",
	UseWildCard:   true,
	Timeout:       15,
}
```

### Todo

- [x] Support wildcards
- [ ] support for multiple hosts file
- [ ] Authentification
- [ ] block by country
- [ ] block domain format IP
- [ ] DNS over HTTPS
