package inbound

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"errors"

	core "github.com/gzjjjfree/hello"
	"github.com/gzjjjfree/hello/common"
	"github.com/gzjjjfree/hello/common/session"

	//"github.com/gzjjjfree/hello/proxy/vmess/encoding"
	"github.com/gzjjjfree/hello/common/mux"
	"github.com/gzjjjfree/hello/common/net"
	"github.com/gzjjjfree/hello/features/routing"
	"github.com/gzjjjfree/hello/proxy"
	socks "github.com/gzjjjfree/hello/proxy/socks"
	"github.com/gzjjjfree/hello/transport/internet"
)

type HandlerConfig struct {
	access          sync.RWMutex
	clients  *internet.ClientsConfig
	tag      string
	sniffing *sniffingConfig
	hub  []internet.Listener
	proxy           proxy.Inbound
	ctx context.Context
	dispatcher      routing.Dispatcher
}



type sniffingConfig struct {
	enabled      bool
	destOverride []string
}

func New(ctx context.Context, config *core.InboundHandlerConfig) (*HandlerConfig, error) {
	fmt.Println("in proxy-vmess-inbound-inbound.go func New ctx: ", ctx)
	//v := core.MustFromContext(ctx)
	//var inboundTag core.Tag = "inboundTag"
	t := mux.NewServer(ctx)
	fmt.Println("the mux.NewServer type is: %T", reflect.TypeOf(t))
	fmt.Println("the mux.NewServer type is: %T", t)
	var s socks.Server
	handler := &HandlerConfig{
		clients: &internet.ClientsConfig{
			Protocol: config.Protocol,
			Port:     config.Port,
			Address:   net.ParseAddress(config.Listen),
		},
		tag: config.Tag,
		sniffing: &sniffingConfig{
			enabled:      config.Sniffing.Enabled,
			destOverride: config.Sniffing.DestOverride,
		},
		//ctx: session.ContextWithInbounds(ctx, config),
		ctx: ctx,
		proxy: &s,
		dispatcher: t,
	}

	return handler, nil
}

func (h *HandlerConfig) callback(conn internet.Connection) {
	//fmt.Println("in proxy-vmess-inbound-inbound.go func (h *HandlerConfig) callback")
	ctx, cancel := context.WithCancel(h.ctx)
	fmt.Println("in proxy-vmess-inbound-inbound.go func (h *HandlerConfig) callback ctx: ", ctx)
	sid := session.NewID()
	ctx = session.ContextWithID(ctx, sid)
	//fmt.Println("in proxy-vmess-inbound-inbound.go func (h *HandlerConfig) callback ctx: ", ctx)
// conn.RemoteAddr 返回远程网络地址（如果已知）
	ctx = session.ContextWithInbound(ctx, &session.Inbound{
		Source:  net.DestinationFromAddr(conn.RemoteAddr()),
		Gateway: net.TCPDestination(h.clients.Address, net.Port(h.clients.Port)),
		Tag:     h.tag,
	})

	content := new(session.Content)
	if h.sniffing.enabled {
		fmt.Println("in proxy-vmess-inbound-inbound.go func (h *HandlerConfig) callback handler.sniffing.enabled")
		content.SniffingRequest.Enabled = h.sniffing.enabled
		content.SniffingRequest.OverrideDestinationForProtocol = h.sniffing.destOverride
		content.SniffingRequest.MetadataOnly = false
	}
	ctx = session.ContextWithContent(ctx, content)
	//fmt.Println("in proxy-vmess-inbound-inbound.go func (h *HandlerConfig) callback ctx: ", ctx)
	if err := h.proxy.Process(ctx, net.Network_TCP, conn, h.dispatcher); err != nil {
		errors.New("connection ends")
	}
	cancel()
	if err := conn.Close(); err != nil {
		errors.New("failed to close connection")
	}
}

func (h *HandlerConfig) Start() error{
	fmt.Println("in proxy-vmess-inbound-inbound.go func (handler *Handler) Start()")
	ctx := context.Background()
	hub, err := internet.ListenTCP(ctx, h.clients, func(conn internet.Connection) {
		fmt.Println("in  app-proxyman-inbound-worker.go func (w *tcpWorker) Start() 开启一个协程 *tcpWorker-callback 等待数据")
		// 开启一个协程等待数据
		go h.callback(conn)
	}) 
	if err != nil {
		return err
	}
	h.hub = append(h.hub, hub)
	return nil
}



func (h *HandlerConfig) Close() error{
	fmt.Println("in proxy-vmess-inbound-inbound.go func (handler *Handler) Close()")
	
	for _, hub := range h.hub {

		
		if err := common.Close(hub); err != nil {
			fmt.Println("in proxy-vmess-inbound-inbound.go func (handler *Handler) Close() err: ", err)
		}
		
	}
	if err := common.Close(h.proxy); err != nil {
		fmt.Println("in proxy-vmess-inbound-inbound.go func (handler *Handler) Close() err: ", err)
	}
	
	return nil
}

func (h *HandlerConfig) Type() interface{}{
	return fmt.Sprint("in proxy-vmess-inbound-inbound.go func (handler *Handler) Type()")
}

func (h *HandlerConfig) Getctx() context.Context{
	return h.ctx
}

var aeadForced = true
var aeadForced2022 = true

func init() {
	fmt.Println("in proxy-vmess-inbound-inbound.go func init()")
	common.RegisterConfig((*core.InboundHandlerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*core.InboundHandlerConfig))
	})
}
