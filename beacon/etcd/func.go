package etcd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/manmanxing/errors"
	"github.com/manmanxing/go_common/util"
	"go.etcd.io/etcd/clientv3"
)

//ETCD 的一些相关操作

//根据key获取value
func GetValue(key string) (value string, err error) {
	defer func() {
		fmt.Println("Get ETCD value now ", time.Now().Format(util.TimeFormatDateTime), "key ", key, "err ", err)
	}()

	if len(strings.TrimSpace(key)) <= 0 {
		return "", nil
	}

	client, err := NewClient()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.KV.Get(ctx, key)

	if resp == nil || len(resp.Kvs) <= 0 {
		err = errors.Wrap(errors.New(fmt.Sprintf("The key=%s does not exist", key)))
		return "", err
	}

	return string(resp.Kvs[0].Value), nil
}

func Put(key, value string) (err error) {
	if len(strings.TrimSpace(key)) <= 0 {
		err = errors.Wrap(errors.New("The key is empty"))
		return
	}

	if strings.HasSuffix(key, "/") {
		err = errors.Wrap(errors.New("The key 不可以'/'结尾"))
		return
	}

	client, err := NewClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	//去抢锁更新
	txn := clientv3.NewKV(client).Txn(ctx)
	//if 不存在key， then 设置它, else 抢锁失败
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Then(clientv3.OpPut(key, value)).Else(clientv3.OpGet(key))
	//提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		err = errors.Wrap(err, "etcd commit transaction err")
		return err
	}
	if !txnResp.Succeeded {
		err = errors.Errorf("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}
	return nil
}
