package common

import (
	"context"
	"reflect"
	"errors"
	"fmt"
	
)
// CreateObject 根据其配置创建一个对象。配置类型必须通过 RegisterConfig() 注册。
func CreateObject(ctx context.Context, config interface{}) (interface{}, error) {
	configType := reflect.TypeOf(config)
	
	creator, found := typeCreatorRegistry[configType]
	if !found {
		fmt.Println("in common-type.go func CreateObject !found")
		return nil, errors.New(configType.String() + " is not registered")
	}
	return creator(ctx, config)
}

// ConfigCreator 是一个通过配置创建对象的函数。
type ConfigCreator func(ctx context.Context, config interface{}) (interface{}, error)

var (
	typeCreatorRegistry = make(map[reflect.Type]ConfigCreator)
)
// RegisterConfig 注册一个全局配置创建者。配置可以为 nil，但必须有一个类型
func RegisterConfig(config interface{}, configCreator ConfigCreator) error {	
	configType := reflect.TypeOf(config)
	if _, found := typeCreatorRegistry[configType]; found {
		return errors.New(configType.Name() + " is already registered")
	}
	typeCreatorRegistry[configType] = configCreator
	return nil
}