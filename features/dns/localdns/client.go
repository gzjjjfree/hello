package localdns

import (
	"fmt"

	"github.com/gzjjjfree/hello/features/dns"
	"github.com/gzjjjfree/hello/common/net"
)

// New create a new dns.Client that queries localhost for DNS.
// 新建一个 dns.Client，用于向 localhost 查询 DNS
func New() *Client {
	return &Client{}
}

// Client is an implementation of dns.Client, which queries localhost for DNS.
// Client 是 dns.Client 的一个实现，它向 localhost 查询 DNS。
type Client struct{}
// Type implements common.HasType.
func (*Client) Type() interface{} {
	return dns.ClientType()
}

// Start implements common.Runnable.
func (*Client) Start() error { return nil }

// Close implements common.Closable.
func (*Client) Close() error { return nil }

// LookupIP implements Client.
func (*Client) LookupIP(host string) ([]net.IP, error) {
	fmt.Println("in deatures-dns-localdns-client.go func (*Client) LookupIP")
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	parsedIPs := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		parsed := net.IPAddress(ip)
		if parsed != nil {
			parsedIPs = append(parsedIPs, parsed.IP())
		}
	}
	if len(parsedIPs) == 0 {
		return nil, dns.ErrEmptyResponse
	}
	return parsedIPs, nil
}

// LookupIPv4 implements IPv4Lookup.
func (c *Client) LookupIPv4(host string) ([]net.IP, error) {
	fmt.Println("in deatures-dns-localdns-client.go func  (c *Client) LookupIPv4")	
	ips, err := c.LookupIP(host)
	if err != nil {
		return nil, err
	}
	ipv4 := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		if len(ip) == net.IPv4len {
			ipv4 = append(ipv4, ip)
		}
	}
	if len(ipv4) == 0 {
		return nil, dns.ErrEmptyResponse
	}
	return ipv4, nil
}

// LookupIPv6 implements IPv6Lookup.
func (c *Client) LookupIPv6(host string) ([]net.IP, error) {
	fmt.Println("in deatures-dns-localdns-client.go func (c *Client) LookupIPv6")	
	ips, err := c.LookupIP(host)
	if err != nil {
		return nil, err
	}
	ipv6 := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		if len(ip) == net.IPv6len {
			ipv6 = append(ipv6, ip)
		}
	}
	if len(ipv6) == 0 {
		return nil, dns.ErrEmptyResponse
	}
	return ipv6, nil
}