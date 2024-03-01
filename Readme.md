# DNS forward

```go
// default config
var config = &DNSConfig{
	HostsFilePath: "/etc/hosts",
	Laddr:         ":53",
	Raddr:         "8.8.8.8:53",
	UseWildCard:   true,
}
```