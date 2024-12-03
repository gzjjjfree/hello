//go:build !confonly
// +build !confonly

package socks

import (
	"context"
	"io"
	"time"
	"fmt"
	"errors"

	//core "github.com/gzjjjfree/hello"
	"github.com/gzjjjfree/hello/common"
	"github.com/gzjjjfree/hello/common/buf"
	//"github.com/gzjjjfree/hello/common/log"
	"github.com/gzjjjfree/hello/common/net"
	"github.com/gzjjjfree/hello/common/protocol"
	//udp_proto "github.com/gzjjjfree/hello/common/protocol/udp"
	"github.com/gzjjjfree/hello/common/session"
	"github.com/gzjjjfree/hello/common/signal"
	"github.com/gzjjjfree/hello/common/task"
	//"github.com/gzjjjfree/hello/features"
	"github.com/gzjjjfree/hello/features/policy"
	"github.com/gzjjjfree/hello/features/routing"
	"github.com/gzjjjfree/hello/transport/internet"
	//"github.com/gzjjjfree/hello/transport/internet/udp"
)

// Server is a SOCKS 5 proxy server
type Server struct {
	config        *ServerConfig
	policyManager policy.Manager
}

// NewServer creates a new Server object.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	fmt.Println("in proxy-socks-serrverr.go func NewServer")
	//v := core.MustFromContext(ctx)
	s := &Server{
		config:        config,
		//policyManager: v.GetFeature(policy.ManagerType()).(policy.Manager),
	}
	return s, nil
}

func (s *Server) policy() policy.Session {
	fmt.Println("in proxy-socks-serrverr.go func (s *Server) policy()")	
	p := policy.SessionDefault()
	
	//fmt.Println("P is: ", p)
	return p
}

// Network implements proxy.Inbound.
func (s *Server) Network() []net.Network {
	fmt.Println("in proxy-socks-serrverr.go func (s *Server) Network()")
	list := []net.Network{net.Network_TCP}
	if s.config.UdpEnabled {
		list = append(list, net.Network_UDP)
	}
	return list
}

// Process implements proxy.Inbound.
func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher routing.Dispatcher) error {
	fmt.Println("in proxy-socks-serrverr.go func (s *Server) Process ctx: ", ctx)	

	switch network {
	case net.Network_TCP:
		return s.processTCP(ctx, conn, dispatcher)
	case net.Network_UDP:
		return errors.New("can't use UDP")
	default:
		return fmt.Errorf("unknown network: %v", network)
	}
}

func (s *Server) processTCP(ctx context.Context, conn internet.Connection, dispatcher routing.Dispatcher) error {
	fmt.Println("in proxy-socks-serrverr.go func (s *Server) processTCP")
	
	inbound := session.InboundFromContext(ctx)
	if inbound == nil || !inbound.Gateway.IsValid() {
		return errors.New("inbound gateway not specified")
	}

	svrSession := &ServerSession{
		config:        s.config,
		address:       inbound.Gateway.Address,
		port:          inbound.Gateway.Port,
		clientAddress: inbound.Source.Address,
	}
// BufferedReader 结构体的字段 Reader 是一个接口，NewReader(conn) 判断读取的 conn 是否是标准的 io.Reader，并进行转换
	reader := &buf.BufferedReader{Reader: buf.NewReader(conn)}
	request, err := svrSession.Handshake(reader, conn)
	// request 一个包含Socks5版本，连接类型，目标地址的请求头
	if err != nil {		
		return errors.New("failed to read request")
	}
	// 为网络连接设置读取超时时间
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		errors.New("failed to clear deadline")
	}

	if request.Command == protocol.RequestCommandTCP {
		fmt.Println("in proxy-socks-server.go func processTCP request.Command == protocol.RequestCommandTCP")
		dest := net.TCPDestination(request.Address, request.Port)  //request.Destination()
		fmt.Printf("TCP Connect in ID: %v request to %v\n", session.IDFromContext(ctx), dest)
		
		return s.transport(ctx, reader, conn, dest, dispatcher)
	}

	return nil
}

func (*Server) handleUDP(c io.Reader) error {
	// The TCP connection closes after this method returns. We need to wait until
	// the client closes it.
	return common.Error2(io.Copy(buf.DiscardBytes, c))
}

func (s *Server) transport(ctx context.Context, reader io.Reader, writer io.Writer, dest net.Destination, dispatcher routing.Dispatcher) error {
	fmt.Println("in proxy-socks-server.go func (s *Server) transport dest is: ", dest)
	ctx, cancel := context.WithCancel(ctx)
	// plcy 默认策略，超时及缓存
	plcy := policy.SessionDefault()
	// timer 定时控件
	timer := signal.CancelAfterInactivity(ctx, cancel, plcy.Timeouts.ConnectionIdle)
	
	ctx = policy.ContextWithBufferPolicy(ctx, plcy.Buffer)
	// link 以目标地址的读写通路
	fmt.Println("in proxy-socks-server.go func (s *Server) transport before link")
	link, err := dispatcher.Dispatch(ctx, dest)
	fmt.Println("in proxy-socks-server.go func (s *Server) transport after link")
	if err != nil {
		return err
	}
	// fmt.Println("in proxy-socks-serrverr.go link is: ", *link)
	requestDone := func() error {
		// 请求连接函数，请求后设置下行连接超时设置 DownlinkOnly = time.Second * 1
		defer timer.SetTimeout(plcy.Timeouts.DownlinkOnly)
		// 把来源端的请求内容读到 link 通路中
		if err := buf.Copy(buf.NewReader(reader), link.Writer, buf.UpdateActivity(timer)); err != nil {
			return errors.New("failed to transport all TCP request")
		}

		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(plcy.Timeouts.UplinkOnly)

		v2writer := buf.NewWriter(writer)
		if err := buf.Copy(link.Reader, v2writer, buf.UpdateActivity(timer)); err != nil {
			return errors.New("failed to transport all TCP response")
		}

		return nil
	}
	fmt.Println("in proxy-socks-serrverr.go task.OnSuccess(requestDone, task.Close(link.Writer))")
	var requestDonePost = task.OnSuccess(requestDone, task.Close(link.Writer))
	if err := task.Run(ctx, requestDonePost, responseDone); err != nil {
		common.Interrupt(link.Reader)
		common.Interrupt(link.Writer)
		return errors.New("connection ends")
	}
	fmt.Println("in proxy-socks-serrverr.go func (s *Server) transport END")
	return nil
}



func init() {
	fmt.Println("in proxy-socks-serrverr.go func init()")
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
