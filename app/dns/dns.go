package dns

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
	//clients  *clientsConfig
	tag      string
	ctx context.Context
}



func New(ctx context.Context, config *core.DnsHandlerConfig) (*Handler, error) {
	//v := core.MustFromContext(ctx)
	var dnsTag core.Tag
	handler := &Handler{
		
		ctx: context.WithValue(ctx, dnsTag, config.Servers),
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


func init() {
	common.RegisterConfig((*core.DnsHandlerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*core.DnsHandlerConfig))
	})
}
