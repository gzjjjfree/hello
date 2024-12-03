package internet

import (
	"context"
	//"runtime"
	"syscall"
	//"errors"
	"fmt"

	//"github.com/pires/go-proxyproto"

	//"github.com/gzjjjfree/hello/common/net"
	//"github.com/gzjjjfree/hello/common/session"
)

var (
	effectiveListener = DefaultListener{}
)

type controller func(network, address string, fd uintptr) error

type DefaultListener struct {
	controllers []controller
}

func getControlFunc(ctx context.Context, sockopt *SocketConfig, controllers []controller) func(network, address string, c syscall.RawConn) error {
	fmt.Println("in transport-internet-system_listener.go func getControlFunc")
	return func(network, address string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			if sockopt != nil {
				if err := applyInboundSocketOptions(network, fd, sockopt); err != nil {
					fmt.Println("failed to apply socket options to incoming connection")
				}
			}

			setReusePort(fd) // nolint: staticcheck

			for _, controller := range controllers {
				if err := controller(network, address, fd); err != nil {
					fmt.Println("failed to apply external controller")
				}
			}
		})
	}
}


