package etcd

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"code.qschou.com/golang/errors"
	"code.qschou.com/qschou/go_center_common/etcd_auth"
	etcd "go.etcd.io/etcd/clientv3"
)

// Endpoints 获取 ETCD 的节点地址列表.
func Endpoints() (endpoints []string, err error) {
	defer func() {
		//这个方法很多时候在init 调用，此时日志尚未初始化，所以输出到终端
		fmt.Println("now", time.Now(), "ETCD_ADDR", endpoints, "err", err)
	}()
	v, ok := os.LookupEnv("ETCD_ADDR")
	if !ok {
		return nil, errors.New("The environment variable ETCD_ADDR not found")
	}
	return strings.Split(strings.TrimSpace(v), ";"), nil
}

// NewClient 创建一个新的 etcd.Client.
// 请注意应用程序需要在合适的时候关闭这个 Client;
// 一般情况下使用 Client() 函数获取 etcd.Client.
func NewClient(endpoints []string) (*etcd.Client, error) {
	clt, err := etcd.New(getEtcdConfig(endpoints))
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return clt, nil
}

var __ClientPointer unsafe.Pointer // *etcd.Client

// Client 返回一个 etcd.Client.
// 一般应用程序不需要关闭这个 Client, 或者在程序退出的时候关闭它.
func Client() (*etcd.Client, error) {
	for {
		p := (*etcd.Client)(atomic.LoadPointer(&__ClientPointer))
		if p != nil {
			return p, nil
		}
		endpoints, err := Endpoints()
		if err != nil {
			return nil, errors.Wrap(err)
		}
		clt, err := NewClient(endpoints)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		if !atomic.CompareAndSwapPointer(&__ClientPointer, nil, unsafe.Pointer(clt)) {
			clt.Close()
			continue
		}
		return clt, nil
	}
}

func getEtcdConfig(endpoints []string) etcd.Config {
	cfg := etcd.Config{
		Endpoints:        endpoints,
		AutoSyncInterval: time.Hour,
		DialTimeout:      time.Second * 5,
	}

	cfg.Username, cfg.Password = etcd_auth.GetEtcdConfig()
	return cfg
}
