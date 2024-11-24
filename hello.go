package core

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/gzjjjfree/hello/features"
	"github.com/gzjjjfree/hello/common"
)

type Server interface { //Server 是 V2Ray 的一个实例，任何时候都最多只能有一个 Server 实例在运行。
	common.Runnable
}

func New(config *Config) (*Instance, error) {
	var server = &Instance{ctx: context.Background()}

	done, err := initInstanceWithConfig(config, server) 
	if done {                                           
		fmt.Println("in hello.go New err is: ", err)
		return nil, err
	}
	return server, nil
}

type Instance struct {
	access             sync.Mutex         // Mutex 是一种互斥锁。Mutex 的零值表示未锁定的互斥锁。首次使用后不得复制 Mutex。
	features           []features.Feature // {common.HasType common.Runnable} Runnable 是可以根据需要开始工作和停止的对象的接口。HasType 是知道其类型的对象的接口
	featureResolutions []resolution       // {deps []reflect.Type callback interface{}} Type 是 Go 类型的表示 callback 一个接口
	running            bool

	ctx context.Context // Context 类型，它携带跨 API 边界和进程之间的截止日期、取消信号和其他请求范围的值
}

type resolution struct {
	deps     []reflect.Type
	callback interface{}
}

func (s *Instance) Start() error {
	s.access.Lock()
	defer s.access.Unlock()

	s.running = true
	for _, f := range s.features {
		fmt.Println("in hello.go func (s *Instance) Start() : ", reflect.TypeOf(f))
		//k := Tag("outboundTag")
		//if v := f.Getctx().Value(k); v != nil {
		//	fmt.Println("in hello.go func (s *Instance) f.ctx : ", v)
		//}
		
		if err := f.Start(); err != nil {
			return err
		}
	}

	fmt.Println("GzV2Ray ", Version(), " started")

	return nil
}

func (s *Instance) Close() error {
	fmt.Println("in gzv2ray.go func (s *Instance) Close()")
	s.access.Lock()
	defer s.access.Unlock()

	s.running = false

	var errorsmsg []interface{}
	for _, f := range s.features {
		if err := f.Close(); err != nil {
			errorsmsg = append(errorsmsg, err)
		}
	}
	if len(errorsmsg) > 0 {
		return errors.New("failed to close all features")
	}

	return nil
}

func (s *Instance) Type() interface{} {
	return ServerType()
}

func ServerType() interface{} {
	return (*Instance)(nil)
}

func initInstanceWithConfig(config *Config, server *Instance) (bool, error) {
	if err := addOutboundHandlers(server, config.Outbounds); err != nil {
		return true, err
	}

	if config.Dns != nil {
		if err := AddHandler(server, config.Dns); err != nil {
			return true, err
		}
	}

	if config.Routing != nil {
		if err := AddHandler(server, config.Routing); err != nil {
			return true, err
		}
	}
	
	if err := addInboundHandlers(server, config.Inbounds); err != nil {
		return true, err
	}

	return false, nil
}

func addInboundHandlers(server *Instance, configs []*InboundHandlerConfig) error {
	for _, inboundConfig := range configs {
		if err := AddHandler(server, inboundConfig); err != nil {
			return err
		}
	}
	return nil
}

func addOutboundHandlers(server *Instance, configs []*OutboundHandlerConfig) error {
	for _, outboundConfig := range configs {
		//fmt.Println("in addOutboundHandlers index is: ", index)
		if err := AddHandler(server, outboundConfig); err != nil {
			return err
		}
	}
	return nil
}

func AddHandler(server *Instance, config interface {}) error {
	rawHandler, err := CreateObject(server, config)
	if err != nil {
		return err
	}
	if feature, ok := rawHandler.(features.Feature); ok {
        server.features = append(server.features, feature)
		return nil
	}
	return fmt.Errorf("not an : %s", reflect.TypeOf(config))
		
}




