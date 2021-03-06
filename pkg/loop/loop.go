package loop

import (
	"context"
	"errors"
	"runtime"
	"sync/atomic"
	"time"
)

// NewLoop creates a Loop struct, can be called from any thread
func NewLoop() *Loop {
	return &Loop{
		calls: []func(context.Context){},
	}
}

// Loop is the main struct that houses the calls
type Loop struct {
	index   int
	running int32
	ctx     context.Context
	cancel  context.CancelFunc
	calls   []func(context.Context)
}

// Run Starts and Blocks the loop until done is closed
// call this method from the thread you want it to be run on
func (l *Loop) Run(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&l.running, 0, 1) {
		return errors.New("Run Allready Called")
	}
	runtime.LockOSThread()
	if ctx == nil {
		ctx = context.Background()
	}
	l.ctx, l.cancel = context.WithCancel(ctx)
	for {
		select {
		case <-l.ctx.Done():
			l.Close()
			return nil
		default:
			if len(l.calls) == 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			if len(l.calls) <= l.index {
				l.index = 0
				continue
			}
			l.calls[l.index](ctx)
			l.index++
		}
	}
}

//Close cleans up the lock
func (l *Loop) Close() error {
	l.cancel()
	runtime.UnlockOSThread()
	return nil
}

// AddCall adds a call that should be run on loop thread, the caller is responsible for establishing channels to recieve any required values from the functions
// the function will be run in the order it is added to the queue, all calls are blocking and must return to run other call in the queue
// so if running multiple loops on the main thread, each func should execute the loop once to yeild to other calls
// i.e.
// errChan:=make(chan error)
// intChan:=make(chan int)
// l.AddCall(func(){
//	err,sum:=math.Add(1,2)
// 	if err!=nil{
//		errChan <- err
//		return
//	}
// 	intChan <- sum
// })
// [...] deal with channel results
func (l *Loop) AddCall(call func(context.Context)) {
	l.calls = append(l.calls, call)
	// l.callQueue <- call
}
