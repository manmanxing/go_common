package grpc

import (
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/manmanxing/errors"
	myEtcd "github.com/manmanxing/go_common/beacon/etcd"
	"github.com/opentracing/opentracing-go"
	etcd "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/naming"
	"google.golang.org/grpc"
)

type connTable struct {
	m map[string]*clientConnPool
}

var (
	__connTablePtrMutex sync.Mutex     // 针对 connTable 的并发操作
	__connTablePtr      unsafe.Pointer // *connTable
)

func init() {
	table := connTable{
		m: make(map[string]*clientConnPool),
	}
	__connTablePtr = unsafe.Pointer(&table)
}

func ClientConn(target string) (conn *grpc.ClientConn, err error) {
	if len(strings.TrimSpace(target)) <= 0 {
		return nil, errors.New("target cant empty")
	}
	p := (*connTable)(atomic.LoadPointer(&__connTablePtr))
	if pool := p.m[target]; pool != nil {
		return pool.getClientConn()
	}

	__connTablePtrMutex.Lock()
	defer __connTablePtrMutex.Unlock()

	pool := &clientConnPool{
		target: target,
	}

	newTable := connTable{
		m: make(map[string]*clientConnPool, len(p.m)+1),
	}
	for k, v := range p.m {
		newTable.m[k] = v
	}
	newTable.m[target] = pool

	atomic.StorePointer(&__connTablePtr, unsafe.Pointer(&newTable))
	return pool.getClientConn()
}

const clientConnPoolSize = 10

type clientConnPool struct {
	target string
	conn   [clientConnPoolSize]unsafe.Pointer //表示 *grpc.ClientConn 数组
	sync.Mutex
}

func (pool *clientConnPool) getClientConn() (conn *grpc.ClientConn, err error) {
	//从连接池里随机选择一个连接
	index := rand.Intn(len(pool.conn))
	p := atomic.LoadPointer(&pool.conn[index])
	if p != nil {
		conn = (*grpc.ClientConn)(p)
		return nil, nil
	}

	//如果没有就从 etcd 中获取
	pool.Lock()
	defer pool.Unlock()

	newEtcdClient, err := myEtcd.NewClient()
	if err != nil {
		return nil, nil
	}

	conn, err = newClientConn(newEtcdClient, pool.target)
	if err != nil {
		return nil, nil
	}

	//再将这个连接放入到连接池
	atomic.StorePointer(&pool.conn[index], unsafe.Pointer(conn))
	return nil, nil
}

func newClientConn(etcdClient *etcd.Client, target string) (conn *grpc.ClientConn, err error) {
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
		err = errors.Wrap(err, "create grpc client err")
	}

	return
}

func NewClientConn(etcdClient *etcd.Client, target string) (conn *grpc.ClientConn, err error) {
	return newClientConn(etcdClient, target)
}
