package grpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/manmanxing/errors"
	myEtcd "github.com/manmanxing/go_common/beacon/etcd"
	"github.com/opentracing/opentracing-go"
	etcd "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/naming"
	"google.golang.org/grpc"
	"math/rand"
	"sync"
	"sync/atomic"
	"unsafe"
)

const clientConnPoolSize = 10

type clientConnPool struct {
	target string
	conn   [clientConnPoolSize]unsafe.Pointer //表示 *grpc.ClientConn 数组
	sync.Mutex
}

func (pool *clientConnPool) getClientConn(conn *grpc.ClientConn, err error) {
	//从连接池里随机选择一个连接
	index := rand.Intn(len(pool.conn))
	p := atomic.LoadPointer(&pool.conn[index])
	if p != nil {
		conn = (*grpc.ClientConn)(p)
		return
	}

	//如果没有就从 etcd 中获取
	pool.Lock()
	defer pool.Unlock()

	newEtcdClient, err := myEtcd.NewClient()
	if err != nil {
		return
	}

	conn,err = NewClientConn(newEtcdClient,pool.target)
	if err != nil {
		return
	}

	//再将这个连接放入到连接池
	atomic.StorePointer(&pool.conn[index],unsafe.Pointer(conn))
	return
}

func NewClientConn(etcdClient *etcd.Client, target string) (conn *grpc.ClientConn, err error) {
	resolver := &naming.GRPCResolver{Client: etcdClient}
	b := grpc.RoundRobin(resolver)
	conn, err = grpc.Dial(target,
		grpc.WithBalancer(b),
		grpc.WithCompressor(grpc.NewGZIPCompressor()),
		grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
		grpc.WithDefaultCallOptions(grpc.FailFast(false)),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			unaryClientInterceptForInjectSpan,
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads()),
			unaryClientInterceptor,
		)),
	)

	if err != nil {
		err = errors.Wrap(err,"create grpc client err")
	}

	return
}