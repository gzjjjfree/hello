package core

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"example.com/gztest"

	"github.com/gzjjjfree/hello/common"
	"github.com/gzjjjfree/hello/features"
	//"github.com/gzjjjfree/hello/features/routing"
)

type Server interface { //Server 是 V2Ray 的一个实例，任何时候都最多只能有一个 Server 实例在运行。
	common.Runnable
}

func New(config *Config) (*Instance, error) {
	fmt.Println("in hello.go func New")
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
		k := Tag("outboundTag")
		if v := f.Getctx().Value(k); v != nil {
			fmt.Println("in hello.go func (s *Instance) Start() : ", reflect.TypeOf(f), "---", v)
		} else {
			fmt.Println("in hello.go func (s *Instance) Start() : ", reflect.TypeOf(f))
		}

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
	fmt.Println("in gzv2ray.go func initInstanceWithConfig")
	if err := addOutboundHandlers(server, config.Outbounds); err != nil {
		return true, err
	}
	fmt.Println("in gzv2ray.go func initInstanceWithConfig after add outboundhandlers")
	//if config.Dns != nil {
	//	if err := AddHandler(server, config.Dns); err != nil {
	//		return true, err
	//	}
	//}

	if config.Routing != nil {
		if err := AddHandler(server, config.Routing); err != nil {
			return true, err
		}
	}
	//if err := AddHandler(server, routing.Dispatcher); err != nil {
	//	return true, err
	//}
	fmt.Println("in gzv2ray.go func initInstanceWithConfig after add config.Routing")
	
	if err := addInboundHandlers(server, config.Inbounds); err != nil {
		return true, err
	}

	for _, f := range server.features {
		fmt.Println("in gzv2ray.go func initInstanceWithConfig  range allFeatures: ", reflect.TypeOf(f.Type()))
		
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

func AddHandler(server *Instance, config interface{}) error {
	rawHandler, err := CreateObject(server, config)
	if err != nil {
		return err
	}
	fmt.Println("in gzv2ray.go func AddHandler  range allFeatures: ", reflect.TypeOf(rawHandler))
	if feature, ok := rawHandler.(features.Feature); ok {
		server.features = append(server.features, feature)
		
		return nil
	}
	return fmt.Errorf("not an : %s", reflect.TypeOf(config))

}

// RequireFeatures is a helper function to require features from Instance in context.
// RequireFeatures 是一个辅助函数，用于在上下文中请求来自实例的特征
// See Instance.RequireFeatures for more information.
// 查看 Instance.RequireFeatures 以了解更多信息。
func RequireFeatures(ctx context.Context, callback interface{}) error {
	fmt.Println("in gzv2ray.go func RequireFeatures")
	v := MustFromContext(ctx)
	return v.RequireFeatures(callback)
}

// RequireFeatures registers a callback, which will be called when all dependent features are registered.
// RequireFeatures 注册一个回调，当所有依赖功能都注册后将调用该回调。
// The callback must be a func(). All its parameters must be features.Feature.
// 回调必须是一个 func()。它的所有参数必须是 features.Feature。
func (s *Instance) RequireFeatures(callback interface{}) error {
	fmt.Println("in gzv2ray.go (s *Instance) RequireFeatures")
	callbackType := reflect.TypeOf(callback)
	// 确认回调接口是函数
	if callbackType.Kind() != reflect.Func {
		panic("not a function")
	}

	var featureTypes []reflect.Type
	//fmt.Println("in gzv2ray.go (s *Instance) RequireFeatures featureTypes: ", featureTypes)
	// featureTypes 汇总回调函数各个参数 feature 的指针
	for i := 0; i < callbackType.NumIn(); i++ {
		fmt.Println("in gzv2ray.go (s *Instance) RequireFeatures I: ", i)
		featureTypes = append(featureTypes, reflect.PointerTo(callbackType.In(i)))
	}

	r := resolution{
		deps:     featureTypes,
		callback: callback,
	}

	if finished, err := r.resolve(s.features); finished {
		fmt.Println("in gzv2ray.go (s *Instance) RequireFeatures err")
		return err
	}
	fmt.Println("in gzv2ray.go (s *Instance) RequireFeatures r.deps: ", r.deps)
	//gztest.GetMessageReflectType(r.deps)
	// 把没有注册的依赖功能类型列表 r.deps 添加到实例 featureResolutions 中
	s.featureResolutions = append(s.featureResolutions, r)
	return nil
}

func (r *resolution) resolve(allFeatures []features.Feature) (bool, error) { // resoleve 解析接口
	fmt.Println("in hello.go func (r *resolution) resolve")
	gztest.GetMessageReflectType(r.deps)
	var fs []features.Feature
	// r 是  Feature 类型列表及回调函数
	for _, d := range r.deps {
		//fmt.Println("in gzv2ray.go func (r *resolution) resolve d: ", d)
		// 在功能中查找 deps 类型匹配
		f := getFeature(allFeatures, d)
		// 找到最后的参数类型 *stats.Manager 后，f 才会 != nill 才会执行后面的代码
		if f == nil { // 当无匹配时，返回
			fmt.Println("in gzv2ray.go func (r *resolution) resolve f == nil")
			return false, nil
		}
		fmt.Println("in gzv2ray.go func (r *resolution) resolve fs = append(fs, f)")
		// 把找到已注册的需要功能汇总到 fs
		fs = append(fs, f) // 当匹配时，汇总到变量 fs
	}
	//fmt.Println("in gzv2ray.go func (r *resolution) resolve callback := reflect.ValueOf(r.callback) len(fs): ", len(fs))
	//
	callback := reflect.ValueOf(r.callback) // ValueOf 返回一个新的值，该值初始化为接口 i 中存储的具体值。 ValueOf(nil) 返回零值。
	var input []reflect.Value
	// 需要 Feature 类型的列表
	callbackType := callback.Type()
	fmt.Println("in gzv2ray.go func (r *resolution) resolve callbackType: ", callbackType)
	for i := 0; i < callbackType.NumIn(); i++ { // NumIn 返回函数类型的输入参数数量。如果类型的 Kind 不是 Func，则会引起混乱。
		fmt.Println("in gzv2ray.go func (r *resolution) resolve i: ", i)
		pt := callbackType.In(i) // In 返回函数类型的第 i 个输入参数的类型。如果类型的 Kind 不是 Func，则会引起混乱。如果 i 不在 [0, NumIn()) 范围内，则会引起混乱。
		for _, f := range fs {
			fmt.Println("in gzv2ray.go func (r *resolution) resolve pt: ", pt)
			fmt.Println("in gzv2ray.go func (r *resolution) resolve reflect.TypeOf(f): ", reflect.TypeOf(f))
			fmt.Printf("%T\n", f)
			fmt.Printf("%T\n", pt)

			//gztest.GetMessageReflectType(f)
			//fmt.Println("in gzv2ray.go func (r *resolution) resolve f := range fs i: ", i, " ", reflect.TypeOf(f.Type()))
			// 判定已注册的需求功能是否能赋值给回调函数当参数
			if reflect.TypeOf(f).AssignableTo(pt) { // AssignableTo 报告该类型的值是否可以分配给类型 u
				fmt.Println("in gzv2ray.go func (r *resolution) reflect.TypeOf(f).AssignableTo(pt)")
				input = append(input, reflect.ValueOf(f)) // 把匹配的  具体值添加到 input
				break                                     // break 从头轮询，是因为有可能回调函数的参数有同类别，所以要再次轮询 ?
			}
		}
	}
	//fmt.Println("in gzv2ray.go func (r *resolution) resolve len(input)： ", len(input), callbackType.NumIn())
	if len(input) != callbackType.NumIn() {
		panic("Can't get all input parameters") // 内置函数 panic 会停止当前 goroutine 的正常执行。
	}

	var err error
	ret := callback.Call(input)                          // Call 使用输入参数 in 调用函数 v。例如，如果 len(in) == 3，则 v.Call(in) 表示 Go 调用 v(in[0], in[1], in[2])。input 就是参数列表
	errInterface := reflect.TypeOf((*error)(nil)).Elem() // TypeOf 返回表示 i 的动态类型的反射 [Type]。如果 i 是 nil 接口值，则 TypeOf 返回 nil。Elem 返回类型的元素类型。
	for i := len(ret) - 1; i >= 0; i-- {
		if ret[i].Type() == errInterface { // ret.Type 有错误或空时，检查回调返回值是否有错误
			v := ret[i].Interface() // 接口以 interface{} 形式返回 v 的当前值。它相当于：var i interface{} = (v 的底层值) 如果通过访问未导出的结构字段获取值，则会引起混乱。
			if v != nil {
				err = v.(error)
			}
			break
		}
	}

	return true, err
}

func getFeature(allFeatures []features.Feature, t reflect.Type) features.Feature {
	fmt.Println("in gzv2ray.go func getFeature(allFeatures []features.Feature, t reflect.Type) t: ", reflect.ValueOf(t))
	for _, f := range allFeatures {
		fmt.Println("in gzv2ray.go func getFeature(allFeatures []features.Feature, t reflect.Type)  range allFeatures: ", reflect.TypeOf(f.Type()))
		if reflect.TypeOf(f.Type()) == t {
			fmt.Println("in gzv2ray.go func getFeature(allFeatures []features.Feature, t reflect.Type) reflect.TypeOf(f.Type()) == t")
			return f
		}
	}
	return nil
}

func (s *Instance) GetFeature(t reflect.Type) features.Feature {
    return getFeature(s.features, t)
}
