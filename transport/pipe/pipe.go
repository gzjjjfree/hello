package pipe

import (
	"context"

	"github.com/gzjjjfree/hello/common/signal"
	"github.com/gzjjjfree/hello/common/signal/done"
	"github.com/gzjjjfree/hello/features/policy"
)

// Option for creating new Pipes.
type Option func(*pipeOption)

// WithoutSizeLimit returns an Option for Pipe to have no size limit.
func WithoutSizeLimit() Option {
	return func(opt *pipeOption) {
		opt.limit = -1
	}
}

// WithSizeLimit returns an Option for Pipe to have the given size limit.
// WithSizeLimit 返回 Pipe 的选项以具有给定的大小限制。
func WithSizeLimit(limit int32) Option {
	return func(opt *pipeOption) {
		opt.limit = limit
	}
}

// DiscardOverflow returns an Option for Pipe to discard writes if full.
func DiscardOverflow() Option {
	return func(opt *pipeOption) {
		opt.discardOverflow = true
	}
}

// OptionsFromContext returns a list of Options from context.
func OptionsFromContext(ctx context.Context) []Option {
	var opt []Option

	bp := policy.BufferPolicyFromContext(ctx)
	if bp.PerConnection >= 0 {
		opt = append(opt, WithSizeLimit(bp.PerConnection))
	} else {
		opt = append(opt, WithoutSizeLimit())
	}

	return opt
}

// New creates a new Reader and Writer that connects to each other.
// New 创建一个新的 Reader 和 Writer，使它们相互连接
func New(opts ...Option) (*Reader, *Writer) {
	p := &pipe{
		readSignal:  signal.NewNotifier(),
		writeSignal: signal.NewNotifier(),
		done:        done.New(),
		option: pipeOption{
			limit: -1,
		},
	}
// 把选项opts写入p.option
	for _, opt := range opts {
		opt(&(p.option))
	}

	return &Reader{
			pipe: p,
		}, &Writer{
			pipe: p,
		}
}
