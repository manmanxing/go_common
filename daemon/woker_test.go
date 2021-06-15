package daemon

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/manmanxing/go_common/daemon/handler"
)

//下面有三个worker，workerName都一样，表示三个节点同时运行相同任务的worker。测试方式如下：
//第一步：任意选择一个worker执行，判断worker是否开始运行："start":"ok"，，判断是否选主成功：isleader：ok，workerName: 启动的worker名称
//第二步：再次选择剩下的worker执行，不会出现 isleader：ok，说明选主失败
//关闭第一步中的worker，会发现 worker2会出现 isleader：ok，说明选主成功

func TestNewWorker1(t *testing.T) {
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	w1 := NewWorker(handler.Handler,
		WithWg(wg),
		WithCtx(ctx),
		WithName("worker1"),
		WithBusySleepTime(time.Second*1),
		WithIdleSleepTime(time.Second*2),
	)
	go w1.Work()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	s := <-c
	fmt.Println("signal.Notify", s)
	cancel()
	wg.Wait()
}

func TestNewWorker2(t *testing.T) {
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	w1 := NewWorker(handler.Handler,
		WithWg(wg),
		WithCtx(ctx),
		WithName("worker1"),
		WithBusySleepTime(time.Second*1),
		WithIdleSleepTime(time.Second*2),
	)
	go w1.Work()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	s := <-c
	fmt.Println("signal.Notify", s)
	cancel()
	wg.Wait()
}

func TestNewWorker3(t *testing.T) {
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	w1 := NewWorker(handler.Handler,
		WithWg(wg),
		WithCtx(ctx),
		WithName("worker1"),
		WithBusySleepTime(time.Second*1),
		WithIdleSleepTime(time.Second*2),
	)
	go w1.Work()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	s := <-c
	fmt.Println("signal.Notify", s)
	cancel()
	wg.Wait()
}
