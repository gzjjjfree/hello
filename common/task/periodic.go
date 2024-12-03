package task

import (
	"sync"
	"time"
	"fmt"
)

// Periodic is a task that runs periodically.
// 定期是定期运行的任务
type Periodic struct {
	// Interval of the task being run
	Interval time.Duration
	// Execute is the task function
	Execute func() error

	access  sync.Mutex
	timer   *time.Timer
	running bool
}

func (t *Periodic) hasClosed() bool {
	t.access.Lock()
	defer t.access.Unlock()

	return !t.running
}

func (t *Periodic) checkedExecute() error {
	// 检查是否正在运行，运行中时直接返回
	if t.hasClosed() {
		return nil
	}

	// 没有在运行时，运行 Execute() == t.check
	if err := t.Execute(); err != nil {
		t.access.Lock()
		t.running = false
		t.access.Unlock()
		return err
	}

	t.access.Lock()
	defer t.access.Unlock()

	if !t.running {
		return nil
	}

	// Timer 类型表示单个事件。当 Timer 到期时，当前时间将在 C 上发送，除非 Timer 是由 [AfterFunc] 创建的。必须使用 [NewTimer] 或 AfterFunc 创建 Timer。
	// AfterFunc 等待持续时间过去，然后在其自己的 goroutine 中调用 f。它返回一个 [Timer]，可以使用其 Stop 方法取消调用。返回的 Timer 的 C 字段未使用，将为 nil。
	t.timer = time.AfterFunc(t.Interval, func() {
		t.checkedExecute()
	}) 

	return nil
}

// Start implements common.Runnable.
// 开始实现 common.Runnable。
func (t *Periodic) Start() error {
	fmt.Println("in common-task-periodic.go func (t *Periodic) Start()")
	// Lock 锁定 t。如果锁已在使用中，则调用 goroutine 将阻塞，直到互斥锁可用。
	t.access.Lock()
	// Unlock 解锁 m。如果 m 在进入 Unlock 时未锁定，则会出现运行时错误。
    // 锁定的 [Mutex] 与特定 goroutine 无关。允许一个 goroutine 锁定 Mutex，然后安排另一个 goroutine 解锁它
	// 如果正在运行，解锁返回
	if t.running {
		t.access.Unlock()
		return nil
	}
	// 如果没有运行，标记为运行，解锁后运行 checkedExecute(), 由于可以在不同位置修改 t ，所以修改前先锁定，完成后解锁
	t.running = true
	t.access.Unlock()

	if err := t.checkedExecute(); err != nil {
		t.access.Lock()
		t.running = false
		t.access.Unlock()
		return err
	}

	return nil
}

// Close implements common.Closable.
func (t *Periodic) Close() error {
	fmt.Println("in common-task-periodic.go func (t *Periodic) Close() ")
	t.access.Lock()
	defer t.access.Unlock()

	t.running = false
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}

	return nil
}
