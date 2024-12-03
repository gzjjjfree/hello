package internet

import (
	"context"
	"errors"
	"fmt"
	"github.com/gzjjjfree/hello/common/net"
)

var (
	transportListenerCache = make(map[string]ListenFunc)
)

type ConnHandler func(Connection)

func RegisterTransportListener(protocol string, listener ListenFunc) error {
	fmt.Println("in ttansport-internet-tcp_hub.go func RegisterTransportListener")
	if _, found := transportListenerCache[protocol]; found {
		return errors.New(" listener already registered")
	}
	transportListenerCache[protocol] = listener
	fmt.Println("in ttansport-internet-tcp_hub.go func RegisterTransportListener transportListenerCache[protocol]: ", transportListenerCache[protocol])
	return nil
}

type Listener interface {
	Close() error
	Addr() net.Addr
}

type ListenFunc func(ctx context.Context, address net.Address, port net.Port, handler ConnHandler) (Listener, error)

type ClientsConfig struct {
	Protocol string
	Port     uint32
	Address   net.Address
}

func ListenTCP(ctx context.Context, settings *ClientsConfig, handler ConnHandler) (Listener, error) {
	fmt.Println("in ttansport-internet-tcp_hub.go func ListenTCP")
	if settings == nil {
		return nil, errors.New("failed to create default stream settings")		
	}
	
	if settings.Address.Family().IsDomain() && settings.Address.Domain() == "localhost" {
		fmt.Println("in ttansport-internet-tcp_hub.go func ListenTCP address = net.LocalHostIP")
		settings.Address = net.LocalHostIP
	}

	if settings.Address.Family().IsDomain() {
		return nil, errors.New("domain address is not allowed for listening: ")
	}

	protocol := "tcp"
	
	listenFunc := transportListenerCache[protocol]
	fmt.Println("in ttansport-internet-tcp_hub.go func ListenTCP protocol: ", listenFunc)
	if listenFunc == nil {
		return nil, errors.New(" listener not registered")
	}
	port, err := net.PortFromInt(settings.Port)
	if err != nil {
		port = net.Port(54322)
	}
	listener, err := listenFunc(ctx, settings.Address, port, handler)
	if err != nil {
		return nil, errors.New("failed to listen on address")
	}
	fmt.Println("in ttansport-internet-tcp_hub.go return ListenTCP")
	return listener, nil
}

func ListenSystem(ctx context.Context, addr net.Addr, sockopt *SocketConfig) (net.Listener, error) {
	fmt.Println("in ttansport-internet-tcp_hub.go func ListenSystem")
	return effectiveListener.Listen(ctx, addr, sockopt)
}

func ListenSystemPacket(ctx context.Context, addr net.Addr, sockopt *SocketConfig) (net.PacketConn, error) {
	fmt.Println("in ttansport-internet-tcp_hub.go func ListenSystemPacket")
	return effectiveListener.ListenPacket(ctx, addr, sockopt)
}