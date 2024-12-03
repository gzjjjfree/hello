// +build !linux,!freebsd
// +build !confonly

package tcp

import (
	"github.com/gzjjjfree/hello/common/net"
	"github.com/gzjjjfree/hello/transport/internet"
)

func GetOriginalDestination(conn internet.Connection) (net.Destination, error) {
	return net.Destination{}, nil
}
