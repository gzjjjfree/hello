package core

import (
	"context"
	"github.com/gzjjjfree/hello/common"
)

func CreateObject(v *Instance, config interface{}) (interface{}, error) {
	//fmt.Println("in functions.go func CreateObject")
	var ctx context.Context
	if v != nil {
		//fmt.Println("in functions.go func CreateObject v != nil")
		ctx = toContext(v.ctx, v)
	}
	//fmt.Println("in functions.go func CreateObject return")
	return common.CreateObject(ctx, config)
}
