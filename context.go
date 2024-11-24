package core

import (
	"context"
)

func toContext(ctx context.Context, v *Instance) context.Context {
	if FromContext(ctx) != v {
		ctx = context.WithValue(ctx, helloKey, v)		
	}
	return ctx
}

func FromContext(ctx context.Context) *Instance {
	if s, ok := ctx.Value(helloKey).(*Instance); ok {
		return s
	}
	return nil
}

// MustFromContext 从给定的上下文返回一个实例，如果不存在则会引起恐慌
func MustFromContext(ctx context.Context) *Instance {
	v := FromContext(ctx)
	if v == nil {
		panic("V is not in context.")
	}
	return v
}

type HelloKey int

type Tag string

const helloKey HelloKey = 1