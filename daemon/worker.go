package daemon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/manmanxing/go_common/beacon/etcd"

	"github.com/labstack/gommon/log"
	"github.com/manmanxing/go_common/util"
)

/**
1. 自定义命名每个worker名称
2. 遇到panic自动重启，超过设置的panic次数会打印错误日志
3. optional参数，有默认参数和支持自定义参数
4. 可以根据是否空闲，休眠对应的时间段
5. 结合etcd实现分布式选主机制
6. 支持用context优雅的停止
*/

const (
	defaultWorkerName           = "default_worker"       //worker名称
	defaultIdleSleepTime        = time.Millisecond * 100 //空闲时worker休眠时间
	defaultBusySleepTime        = time.Second * 5        // 处理任务期间worker休眠时间
	defaultMaxPanicRecoverCount = 5                      //遇到panic就会自动重启 goroutine，该值表示从panic中恢复的最大次数，
)

type TaskFunc func() (idle bool, err error)
type Option func(info *WorkerInfo)

type WorkerInfo struct {
	Name                 string
	IdleSleepTime        time.Duration
	BusySleepTime        time.Duration
	MaxPanicRecoverCount int
	Task                 TaskFunc
	ctx                  context.Context
	wg                   *sync.WaitGroup
}

func (w *WorkerInfo) Work() {
	defer func() {
		if r := recover(); r != nil {
			stackInfo := util.StackInfo(util.Callers(3))
			log.Error("panic", r, "stack", stackInfo)
			if w.MaxPanicRecoverCount == 0 {
				log.Error("worker name:", w.Name, "panic too many times", "recover info:", r)
			} else {
				w.MaxPanicRecoverCount--
				go w.Work() //restart goroutine
			}
		}
		if w.wg != nil {
			w.wg.Done()
		}
	}()
	client, err := etcd.NewClient()
	if err != nil {
		fmt.Println("get etcd client err:", err)
		return
	}
	isLeader := Campaign(client, w.ctx, w.wg)
	for {
		select {
		case <-w.ctx.Done():
			fmt.Println(fmt.Sprintf("worker %s 关闭...", w.Name))
			return
		case <-isLeader: //监听选主成功信号，若选主失败，会阻塞
			//开始执行任务
			idle, e := w.Task()
			if e != nil {
				log.Error("worker name:", w.Name, "err", e.Error())
			}
			if idle {
				time.Sleep(w.IdleSleepTime)
			} else {
				time.Sleep(w.BusySleepTime)
			}
		}
		time.Sleep(time.Second * 2)
	}
}

func NewWorker(task TaskFunc, opts ...Option) (w *WorkerInfo) {
	w = initWorker(task)
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func initWorker(task TaskFunc) (w *WorkerInfo) {
	return &WorkerInfo{
		Name:                 defaultWorkerName,
		IdleSleepTime:        defaultIdleSleepTime,
		BusySleepTime:        defaultBusySleepTime,
		MaxPanicRecoverCount: defaultMaxPanicRecoverCount,
		Task:                 task,
		ctx:                  context.TODO(),
		wg:                   new(sync.WaitGroup),
	}
}

//应用自定义配置
func WithName(name string) Option {
	return func(info *WorkerInfo) {
		info.Name = name
	}
}

func WithIdleSleepTime(IdleSleepTime time.Duration) Option {
	return func(info *WorkerInfo) {
		info.IdleSleepTime = IdleSleepTime
	}
}

func WithBusySleepTime(BusySleepTime time.Duration) Option {
	return func(info *WorkerInfo) {
		info.BusySleepTime = BusySleepTime
	}
}

func WithMaxPanicRecoverCount(MaxPanicRecoverCount int) Option {
	return func(info *WorkerInfo) {
		info.MaxPanicRecoverCount = MaxPanicRecoverCount
	}
}

func WithCtx(ctx context.Context) Option {
	return func(info *WorkerInfo) {
		info.ctx = ctx
	}
}

func WithWg(wg *sync.WaitGroup) Option {
	return func(info *WorkerInfo) {
		info.wg = wg
	}
}
