package internet

import (
	"context"
	"runtime"
	"syscall"
	"errors"
	"fmt"

	"github.com/pires/go-proxyproto"

	"github.com/gzjjjfree/hello/common/net"
)

func (dl *DefaultListener) Listen(ctx context.Context, addr net.Addr, sockopt *SocketConfig) (net.Listener, error) {
	fmt.Println("in transport-internet-system_listener.go func (dl *DefaultListener) Listen")
	var lc net.ListenConfig
	var l net.Listener
	var err error
	var network, address string
	switch addr := addr.(type) {
	case *net.TCPAddr:
		network = addr.Network()
		address = addr.String()
		lc.Control = getControlFunc(ctx, sockopt, dl.controllers)
	case *net.UnixAddr:
		lc.Control = nil
		network = addr.Network()
		address = addr.Name
		if (runtime.GOOS == "linux" || runtime.GOOS == "android") && address[0] == '@' {
			// linux abstract unix domain socket is lockfree
			if len(address) > 1 && address[1] == '@' {
				// but may need padding to work with haproxy
				fullAddr := make([]byte, len(syscall.RawSockaddrUnix{}.Path))
				copy(fullAddr, address[1:])
				address = string(fullAddr)
			}
		} else {
			// normal unix domain socket needs lock
			locker := &FileLocker{
				path: address + ".lock",
			}
			err := locker.Acquire()
			if err != nil {
				return nil, err
			}
			ctx = context.WithValue(ctx, address, locker) // nolint: golint,staticcheck
		}
	}

	l, err = lc.Listen(ctx, network, address)
	if sockopt != nil && sockopt.AcceptProxyProtocol {
		policyFunc := func(upstream net.Addr) (proxyproto.Policy, error) { return proxyproto.REQUIRE, nil }
		l = &proxyproto.Listener{Listener: l, Policy: policyFunc}
	}
	return l, err
}

func (dl *DefaultListener) ListenPacket(ctx context.Context, addr net.Addr, sockopt *SocketConfig) (net.PacketConn, error) {
	fmt.Println("in transport-internet-system_listener.go func (dl *DefaultListener) ListenPacket")
	var lc net.ListenConfig

	lc.Control = getControlFunc(ctx, sockopt, dl.controllers)

	return lc.ListenPacket(ctx, addr.Network(), addr.String())
}

// RegisterListenerController adds a controller to the effective system listener.
// The controller can be used to operate on file descriptors before they are put into use.
//
// v2ray:api:beta
func RegisterListenerController(controller func(network, address string, fd uintptr) error) error {
	fmt.Println("in transport-internet-system_listener.go func RegisterListenerController")
	if controller == nil {
		return errors.New("nil listener controller")
	}

	effectiveListener.controllers = append(effectiveListener.controllers, controller)
	return nil
}
