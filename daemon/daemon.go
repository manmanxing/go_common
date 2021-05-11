package daemon

import (
	"context"
	"sync"
	"time"

	"github.com/manmanxing/go_common/daemon/handler"
)

//通用worker

func InitDaemon(wg *sync.WaitGroup) (cancel context.CancelFunc) {
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	wg.Add(1)
	w1 := NewWorker(
		handler.Handler,
		WithName("worker1"),
		WithWg(wg),
		WithBusySleepTime(time.Millisecond*5),
		WithIdleSleepTime(time.Millisecond*4),
		WithMaxPanicRecoverCount(5),
		WithCtx(ctx),
	)
	go w1.Work()
	return
}
