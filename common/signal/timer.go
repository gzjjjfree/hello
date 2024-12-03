package signal

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gzjjjfree/hello/common"
	"github.com/gzjjjfree/hello/common/task"
)

type ActivityUpdater interface {
	Update()
}

type ActivityTimer struct {
	sync.RWMutex
	updated   chan struct{}
	checkTask *task.Periodic
	onTimeout func()
}

func (t *ActivityTimer) Update() {
	select {
	case t.updated <- struct{}{}:
	default:
	}
}

func (t *ActivityTimer) check() error {
	// 一直等待重置 t ，接到重置信号后，结束这个 t
	select {
	case <-t.updated:
	default:
		t.finish()
	}
	return nil
}

func (t *ActivityTimer) finish() {
	t.Lock()
	defer t.Unlock()
	// 如果 onTimeout 有函数 cancel , 执行后清空
	if t.onTimeout != nil {
		t.onTimeout()
		t.onTimeout = nil
	}
	// 如果有定期任务 checkTask 把它关闭
	if t.checkTask != nil {
		t.checkTask.Close()
		t.checkTask = nil
	}
}

func (t *ActivityTimer) SetTimeout(timeout time.Duration) {
	fmt.Println("in common-cignal-timer.go func (t *ActivityTimer) SetTimeout")
	// 响应超时为 0 时，结束超时控件
	if timeout == 0 {
		t.finish()
		return
	}

	
	checkTask := &task.Periodic{
		Interval: timeout,
		Execute:  t.check,
	}

	t.Lock()

	// 之前有定期任务时，先关闭
	if t.checkTask != nil {
		t.checkTask.Close()
	}
	// 设置检查定期任务
	t.checkTask = checkTask
	t.Unlock()
	// 先重置 t , 然后开始运行定期任务
	t.Update()
	common.Must(checkTask.Start())
}

func CancelAfterInactivity(ctx context.Context, cancel context.CancelFunc, timeout time.Duration) *ActivityTimer {
	timer := &ActivityTimer{
		updated:   make(chan struct{}, 1),
		onTimeout: cancel,
	}
	timer.SetTimeout(timeout)
	fmt.Println("in common-signal-tmer.go func CancelAfterInactivity return timer")
	return timer
}
