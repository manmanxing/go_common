package etcd

import (
	"fmt"
	"strings"
	"time"

	"path"

	"code.qschou.com/golang/errors"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
)

// Deprecated; Get 从 ETCD 获取某个 key 对应的 value.
func Get(key string) (value string, err error) {
	defer func() {
		// 这个方法很多时候在 init 调用, 此时日志尚未初始化, 所以输出到终端.
		fmt.Println("now", time.Now(), "beacon.etcd.Get", "key", key, "value", value, "error", err)
	}()

	clt, err := Client()
	if err != nil {
		return "", errors.Wrap(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := clt.KV.Get(ctx, key)
	if err != nil {
		return "", errors.Wrapf(err, "key=%s", key)
	}
	if len(resp.Kvs) == 0 {
		return "", errors.Errorf("The key=%s does not exist", key)
	}
	return string(resp.Kvs[0].Value), nil
}

// GetWithPrefix 从 ETCD 获取某个 前缀匹配 对应的 list.
func GetWithPrefix(key string) (result map[string]string, err error) {
	defer func() {
		// 这个方法很多时候在 init 调用, 此时日志尚未初始化, 所以输出到终端.
		fmt.Println("now", time.Now(), "beacon.etcd.Get.WithPrefix", "key", key, "result", result, "error", err)
	}()

	clt, err := Client()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	response, e := clt.Get(ctx, key, clientv3.WithPrefix())
	if e != nil {
		panic(e)
	}
	result = make(map[string]string)
	for _, v := range response.Kvs {
		k := string(v.Key)
		k = path.Base(k)
		result[k] = string(v.Value)
	}
	return
}

//过滤目录, 此方法必须以/结尾
func GetWithPrefixWithoutDir(key string) (result map[string]string, err error) {
	defer func() {
		// 这个方法很多时候在 init 调用, 此时日志尚未初始化, 所以输出到终端.
		fmt.Println("now", time.Now(), "beacon.etcd.Get.WithPrefix", "key", key, "result", result, "error", err)
	}()
	if !strings.HasSuffix(key, "/") {
		return nil, errors.New("key必须以'/'结尾")
	}
	clt, err := Client()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	response, e := clt.Get(ctx, key, clientv3.WithPrefix())
	if e != nil {
		panic(e)
	}
	result = make(map[string]string)
	for _, v := range response.Kvs {
		if strings.Contains(string(v.Value), "_dir_") {
			continue
		}
		k := string(v.Key)
		k = k[len(key):]
		result[k] = string(v.Value)
	}
	return
}

func Put(key, value string) (err error) {
	if strings.HasSuffix(key, "/") {
		err = errors.New("key不可以'/'结尾")
		return
	}
	clt, e := Client()
	if e != nil {
		err = errors.Wrap(e)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// 1 去抢锁（etcd里面的锁就是一个key）有重试机制
	kv := clientv3.NewKV(clt)
	// 创建事物
	txn := kv.Txn(ctx)
	//if 不存在key， then 设置它, else 抢锁失败
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, value)).
		Else(clientv3.OpGet(key))
	// 提交事务
	txnResp, e := txn.Commit()
	if e != nil {
		err = errors.Wrap(e)
		return
	}
	if !txnResp.Succeeded {
		err = errors.Errorf("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}
	return
}
