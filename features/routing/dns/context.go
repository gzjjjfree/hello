package dns



import (
	"github.com/gzjjjfree/hello/common/net"
	"github.com/gzjjjfree/hello/features/dns"
	"github.com/gzjjjfree/hello/features/routing"
	"fmt"
)

// ResolvableContext is an implementation of routing.Context, with domain resolving capability.
type ResolvableContext struct {
	routing.Context
	dnsClient   dns.Client
	resolvedIPs []net.IP
}

// GetTargetIPs overrides original routing.Context's implementation.
func (ctx *ResolvableContext) GetTargetIPs() []net.IP {
	if ips := ctx.Context.GetTargetIPs(); len(ips) != 0 {
		return ips
	}

	if len(ctx.resolvedIPs) > 0 {
		return ctx.resolvedIPs
	}

	if domain := ctx.GetTargetDomain(); len(domain) != 0 {
		var lookupFunc func(string) ([]net.IP, error) = ctx.dnsClient.LookupIP
		ipOption := &dns.IPOption{
			IPv4Enable: true,
			IPv6Enable: true,
		}

		if c, ok := ctx.dnsClient.(dns.ClientWithIPOption); ok {
			ipOption = c.GetIPOption()
			c.SetFakeDNSOption(false) // Skip FakeDNS.
		} else {
			fmt.Println("ctx.dnsClient doesn't implement ClientWithIPOption")
		}

		switch {
		case ipOption.IPv4Enable && !ipOption.IPv6Enable:
			if lookupIPv4, ok := ctx.dnsClient.(dns.IPv4Lookup); ok {
				lookupFunc = lookupIPv4.LookupIPv4
			} else {
				fmt.Println("ctx.dnsClient doesn't implement IPv4Lookup. Use LookupIP instead.")
			}
		case !ipOption.IPv4Enable && ipOption.IPv6Enable:
			if lookupIPv6, ok := ctx.dnsClient.(dns.IPv6Lookup); ok {
				lookupFunc = lookupIPv6.LookupIPv6
			} else {
				fmt.Println("ctx.dnsClient doesn't implement IPv6Lookup. Use LookupIP instead.")
			}
		}

		ips, err := lookupFunc(domain)
		if err == nil {
			ctx.resolvedIPs = ips
			return ips
		}
		fmt.Println("resolve ip for ", domain)
	}

	return nil
}

// ContextWithDNSClient creates a new routing context with domain resolving capability.
// Resolved domain IPs can be retrieved by GetTargetIPs().
func ContextWithDNSClient(ctx routing.Context, client dns.Client) routing.Context {
	return &ResolvableContext{Context: ctx, dnsClient: client}
}
