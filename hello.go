package core

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/gzjjjfree/hello/features"
)

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

func initInstanceWithConfig(config *Config, server *Instance) (bool, error) {
	if err := addInboundHandlers(server, config.Inbounds); err != nil {
		return true, err
	}
	return false, nil
}

func addInboundHandlers(server *Instance, configs []*InboundHandlerConfig) error {
	for _, inboundConfig := range configs {
		if err := AddInboundHandler(server, inboundConfig); err != nil {
			return err
		}
	}
	return nil
}

func AddInboundHandler(server *Instance, config *InboundHandlerConfig) error {
	rawHandler, err := CreateObject(server, config)
	if err != nil {
		return err
	}
	_, ok := rawHandler.(features.Feature)
	if !ok {
		return errors.New("not an InboundHandler")
	}
	return nil
}
