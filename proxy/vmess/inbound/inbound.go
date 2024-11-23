package inbound

import (
	"context"
	"fmt"
	"sync"
	"errors"

	core "github.com/gzjjjfree/hello"
	"github.com/gzjjjfree/hello/common"
	//"github.com/gzjjjfree/hello/proxy/vmess/encoding"
)

type Handler struct {
	access          sync.RWMutex
	clients  *clientsConfig
	tag      string
	sniffing *sniffingConfig
	ctx context.Context
}

type clientsConfig struct {
	protocol string
	port     int32
	listen   string
}

type sniffingConfig struct {
	enabled      bool
	destOverride []string
}

func New(ctx context.Context, config *core.InboundHandlerConfig) (*Handler, error) {
	//v := core.MustFromContext(ctx)
	handler := &Handler{
		clients: &clientsConfig{
			protocol: config.Protocol,
			port:     config.Port,
			listen:   config.Listen,
		},
		tag: config.Tag,
		sniffing: &sniffingConfig{
			enabled:      config.Sniffing.Enabled,
			destOverride: config.Sniffing.DestOverride,
		},
		ctx: ctx,
	}

	return handler, nil
}

func (handler *Handler) Start() error{
	fmt.Println("in proxy-vmess-inbound-inbound.go func (handler *Handler) Start()")
	return errors.New("Start")
}

func (handler *Handler) Close() error{
	fmt.Println("in proxy-vmess-inbound-inbound.go func (handler *Handler) Close()")
	return errors.New("Close")
}

func (handler *Handler) Type() interface{}{
	return fmt.Sprint("in proxy-vmess-inbound-inbound.go func (handler *Handler) Type()")
}

var aeadForced = true
var aeadForced2022 = true

func init() {
	common.RegisterConfig((*core.InboundHandlerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*core.InboundHandlerConfig))
	})
}
