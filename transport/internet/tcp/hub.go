//go:build !confonly
// +build !confonly

package tcp

import (
	"context"
	gotls "crypto/tls"

	//"errors"
	"fmt"
	"strings"
	"time"

	//systemnet "net"
	//"github.com/gzjjjfree/hello/common"
	"github.com/gzjjjfree/hello/common/net"

	//"github.com/gzjjjfree/hello/common/session"
	"github.com/gzjjjfree/hello/transport/internet"
	"github.com/gzjjjfree/hello/transport/internet/tls"
)

// Listener is an internet.Listener that listens for TCP connections.
// Listener 是一个用于监听 TCP 连接的 internet.Listener
type Listener struct {
	listener   net.Listener
	tlsConfig  *gotls.Config
	authConfig internet.ConnectionAuthenticator
	addConn    internet.ConnHandler
	locker     *internet.FileLocker // for unix domain socket
}

func ListenTCP(ctx context.Context, address net.Address, port net.Port, handler internet.ConnHandler) (internet.Listener, error) {
	fmt.Println("in transport-internet-tcp-hub.go func ListenTCP")
	l := &Listener{}
	var lc net.ListenConfig
	addr := &net.TCPAddr{
		IP:   address.IP(),
		Port: int(port),
	}
	ln, err := lc.Listen(ctx, "tcp", addr.String())
	if err != nil {
		panic(err)
	}
	l.listener = ln
	l.addConn = handler
	go l.keepAccepting()
	return l, nil
}

/*
// ListenTCP creates a new Listener based on configurations.

	func ListenTCP(ctx context.Context, address net.Address, port net.Port, handler internet.ConnHandler) (internet.Listener, error) {
		fmt.Println("in transport-internet-tcp-hub.go func ListenTCP")
		l := &Listener{
			addConn: handler,
		}

		var listener net.Listener
		var err error
		var SocketSettings = &internet.SocketConfig{}
		var AcceptProxyProtocol = false
		if port == net.Port(0) { // unix
			listener, err = internet.ListenSystem(ctx, &net.UnixAddr{
				Name: address.Domain(),
				Net:  "unix",
			}, SocketSettings)
			if err != nil {
				return nil, err  // errors.New("failed to listen Unix Domain Socket on ")
			}
			fmt.Println("listening Unix Domain Socket on ", address)
			locker := ctx.Value(address.Domain())
			if locker != nil {
				l.locker = locker.(*internet.FileLocker)
			}
		} else {
			listener, err = internet.ListenSystem(ctx, &net.TCPAddr{
				IP:   address.IP(),
				Port: int(port),
			}, SocketSettings)
			if err != nil {
				return nil, err  // errors.New("failed to listen TCP on")
			}
			fmt.Println("listening TCP on ", address, ":", port)
		}

		if SocketSettings != nil && AcceptProxyProtocol {
			fmt.Println("accepting PROXY protocol")
		}

		l.listener = listener

		// 开启一个协程，监听接收的数据
		go l.keepAccepting()
		return l, nil
	}
*/
func (v *Listener) keepAccepting() {
	fmt.Println("in transport-internet-tcp-hub.go func (v *Listener) keepAccepting()")
	for {
		fmt.Println("等待转入下一个 TCP 连接的信号")
		// Accept 等待并返回下一个连接给监听器
		conn, err := v.listener.Accept()
		if err != nil {
			fmt.Println("测试是否在等待信号.......")
			errStr := err.Error()
			if strings.Contains(errStr, "closed") {
				break
			}
			fmt.Println("failed to accepted raw connections")
			if strings.Contains(errStr, "too many") {
				time.Sleep(time.Millisecond * 500)
			}
			continue
		}
		// Config 结构用于配置 TLS 客户端或服务器。
		if v.tlsConfig != nil {
			fmt.Println("in transport-internet-tcp-hub.go func (v *Listener) keepAccepting() v.tlsConfig != nil ")
			conn = tls.Server(conn, v.tlsConfig)
		}
		if v.authConfig != nil {
			fmt.Println("in transport-internet-tcp-hub.go func (v *Listener) keepAccepting() v.authConfig != nil ")
			conn = v.authConfig.Server(conn)
		}
		fmt.Println("转入下一个 TCP 连接")
		v.addConn(conn)
		//v.addConn(internet.Connection(conn))
	}
	fmt.Println("in transport-internet-tcp-hub.go func (v *Listener) keepAccepting()  END")
}

// Addr implements internet.Listener.Addr.
func (v *Listener) Addr() net.Addr {
	return v.listener.Addr()
}

// Close implements internet.Listener.Close.
func (v *Listener) Close() error {
	fmt.Println("in transport-internet-tcp-hub.go func (v *Listener) Close()")
	if v.locker != nil {
		v.locker.Release()
	}
	return v.listener.Close()
}

const protocolName = "tcp"

func init() {
	fmt.Println("in transport-internet-tcp-hub.go func init()")
	internet.RegisterTransportListener(protocolName, ListenTCP)
}
