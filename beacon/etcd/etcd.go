package etcd

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/manmanxing/errors"
	"github.com/manmanxing/go_common/util"
	etcd "go.etcd.io/etcd/clientv3"
)

const (
	address = "ETCD_ADDR"
)

// *etcd.Client
var __ClientPointer unsafe.Pointer

//获取 ETCD_ADDR 节点地址列表
func getEndPoints() (endPoints []string, err error) {
	defer func() {
		fmt.Println("getEndPoints now ", time.Now().Format(util.TimeFormatDateTime), "ETCD_ADDR ", endPoints, "err ", err)
	}()

	v, ok := os.LookupEnv(address)
	if !ok {
		return nil, errors.New("The environment variable ETCD_ADDR not found")
	}
	return strings.Split(strings.TrimSpace(v), ";"), nil
}

//根据 ETCD 节点地址列表获取 etcd 配置
func getEtcdConfig(endPoints []string) etcd.Config {
	if len(endPoints) <= 0 {
		panic("endPoints is empty")
	}

	cfg := etcd.Config{
		Endpoints:        endPoints,
		AutoSyncInterval: time.Hour,
		DialTimeout:      time.Second * 5,
	}

	return cfg
}

func getClient(endPoints []string) (*etcd.Client, error) {
	client, err := etcd.New(getEtcdConfig(endPoints))
	if err != nil {
		err = errors.Wrap(err, "get etcd client err")
		return nil, err
	}
	return client, nil
}

func GetClient(endPoints []string) (*etcd.Client, error) {
	return getClient(endPoints)
}

func NewClient() (*etcd.Client, error) {
	for {
		p := (*etcd.Client)(atomic.LoadPointer(&__ClientPointer))
		if p != nil {
			return p, nil
		}

		endPoints, err := getEndPoints()
		if err != nil {
			return nil, err
		}

		client, err := getClient(endPoints)
		if err != nil {
			return nil, err
		}
		//只要是没有发生交换，就会一直循环
		if !atomic.CompareAndSwapPointer(&__ClientPointer, nil, unsafe.Pointer(client)) {
			err = client.Close()
			if err != nil {
				return nil, err
			}
			continue
		}
		return client, nil
	}
}
