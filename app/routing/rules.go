package routing

import (
	"context"
	"fmt"
	"sync"
	//"errors"

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



func New(ctx context.Context, config *core.RoutingHandlerConfig) (*Handler, error) {
	//v := core.MustFromContext(ctx)
	var rulesTag core.Tag = "rulesTag"
	handler := &Handler{
		
		ctx: context.WithValue(ctx, rulesTag, "rules"),
	}

	return handler, nil
}

func (handler *Handler) Start() error {
	//fmt.Println("in proxy-vmess-inbound-inbound.go func (handler *Handler) Start()")
	return nil
}

func (handler *Handler) Close() error {
	fmt.Println("in proxy-vmess-inbound-inbound.go func (handler *Handler) Close()")
	return nil
}

func (handler *Handler) Type() interface{} {
	return fmt.Sprint("in proxy-vmess-inbound-inbound.go func (handler *Handler) Type()")
}

func (handler *Handler) Getctx() context.Context {
	return handler.ctx
}

func init() {
	common.RegisterConfig((*core.RoutingHandlerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*core.RoutingHandlerConfig))
	})
}
