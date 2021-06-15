package daemon

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"

	"sync"
	"time"

	"github.com/manmanxing/go_common/util"
)

//这里返回的 channel 表示已经成功选主为 leader
//所有worker可以监听这个channel，这种实现可以让worker阻塞等待节点成为leader，而不是轮询是否是leader节点。
func Campaign(c *clientv3.Client, parentCtx context.Context, wg *sync.WaitGroup, CampaignPrefix string) (success <-chan struct{}) {
	//这里选择当前机器的ip作为etcd的value
	ip := util.GetLocalIP().String()
	ctx, _ := context.WithCancel(parentCtx)

	if wg != nil {
		wg.Add(1)
	}

	hasLeader := make(chan struct{}, 1)

	go func() {
		defer func() {
			if wg != nil {
				wg.Done()
			}
		}()

		for {
			select {
			case <-ctx.Done(): //监听context是否取消
				return
			default:
				//没有取消将走下面的流程
			}
			//创建一个新的 session，并设置它的租期时间
			//session中keepAlive机制会一直续租，如果keepAlive断掉，session.Done会收到退出信号
			s, err := concurrency.NewSession(c, concurrency.WithTTL(5))
			if err != nil {
				fmt.Println("etcd new session err", err)
				time.Sleep(time.Second * 2) //防止重试频率太高
				continue
			}
			//根据 session 与 keyPrefix 开始选主
			e := concurrency.NewElection(s, CampaignPrefix)
			//调用Campaign方法，成为leader的节点会运行出来，非leader节点会阻塞在里面。
			//ctx.done 会导致这里报错
			if err = e.Campaign(ctx, ip); err != nil {
				fmt.Println("etcd campaign err", err)
				_ = s.Close()
				time.Sleep(time.Second * 2) //防止重试频率太高
				continue
			}
			//到这里表示选主成功
			fmt.Println("election leader success,ip:", ip)
			shouldBreak := false
			for !shouldBreak { //会不断的告诉其他的运行worker的服务，这里已经选主成功
				select {
				case hasLeader <- struct{}{}: //选主成功信号，channel 满了会阻塞
				case <-s.Done(): //如果与etcd断开了keepAlive，这里break，重新创建session，重新选举
					fmt.Println("session has done")
					shouldBreak = true
					break
				case <-ctx.Done(): //context被cancel
					ctxTmp, _ := context.WithTimeout(context.Background(), time.Second*1)
					_ = e.Resign(ctxTmp) //让 leader 开始新的选主
					_ = s.Close()        //关闭 session
					return
				}
			}
		}
	}()

	return hasLeader
}
