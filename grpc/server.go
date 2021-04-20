package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/labstack/gommon/log"
	"github.com/manmanxing/go_common/beacon"
	myEtcd "github.com/manmanxing/go_common/beacon/etcd"
	"github.com/opentracing/opentracing-go"
	etcd "go.etcd.io/etcd/clientv3"
	etcdnaming "go.etcd.io/etcd/clientv3/naming"
	etcdrpctypes "go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	grpcnaming "google.golang.org/grpc/naming"
	"google.golang.org/grpc/reflection"
)

const ServiceLeaseTTL = 3

type ServiceRegister func(s *grpc.Server)

type Server struct {
	target     string
	host       string
	port       int
	grpcServer *grpc.Server

	closeNotifyCh chan struct{} //服务关闭 channel 通知

	mutex      sync.Mutex
	stopped    bool //标志grpc server 是否关闭
	etcdLease  etcd.Lease
	etcdClient *etcd.Client // 进程级别的 etcd.Client, 无需关闭
}

func NewServer(target, host string, port int, register ServiceRegister) *Server {
	return NewServerWithOption(
		target,
		host,
		port,
		register,
		grpc.RPCCompressor(grpc.NewGZIPCompressor()),
		grpc.RPCDecompressor(grpc.NewGZIPDecompressor()),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads()),
			unaryServerInterceptor,
		)),
		//这个连接最大的空闲时间，超过就释放，解决proxy等到网络问题（不通知grpc的client和server）
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
}

func NewServerWithOption(target, host string, port int, register ServiceRegister, opt ...grpc.ServerOption) *Server {
	srv := grpc.NewServer(opt...)
	register(srv)
	//开启refletion,方便grpcurl使用
	reflection.Register(srv)
	return &Server{
		target:        target,
		host:          host,
		port:          port,
		grpcServer:    srv,
		closeNotifyCh: make(chan struct{}),
	}
}

func (s *Server) Serve() (err error) {
	s.mutex.Lock()

	if s.stopped {
		s.mutex.Unlock()
		return grpc.ErrServerStopped
	}

	if s.etcdClient == nil {
		s.etcdClient, err = myEtcd.NewClient()
		if err != nil {
			s.mutex.Unlock()
			return err
		}
	}

	if len(strings.TrimSpace(s.host)) <= 0 {
		s.host, err = beacon.ServiceHost()
		if err != nil {
			s.mutex.Unlock()
			return err
		}
	}

	s.mutex.Unlock()

	//开始建立 tcp 连接
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	cancelLoopCh := make(chan struct{})
	defer func() {
		close(cancelLoopCh)
	}()
	//新起一个协程进行etcd续约
	go s.registerServiceToEtcdLoop(cancelLoopCh)
	//这里return之前，一定会有一个准备好的连接
	return s.grpcServer.Serve(lis)
}

func (s *Server) registerServiceToEtcdLoop(cancelLoopCh <-chan struct{}) {
	select {
	case <-cancelLoopCh:
		return
	case <-s.closeNotifyCh:
		return
	case <-time.After(time.Millisecond * 10): //10毫秒后listen 至少已经监听过一次连接了
	}

	leaseKeepAliveResponseChan, leaseId, err := s.registerServiceToEtcd()
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-cancelLoopCh:
			if err = s.etcdLeaseClose(leaseId); err != nil {
				log.Error("Server.etcdLeaseClose()", leaseId, "error", err.Error())
			}
			log.Info("server.etcd", "cancelLoopCh")
			return
		case <-s.closeNotifyCh:
			if err = s.etcdLeaseClose(leaseId); err != nil {
				log.Error("Server.etcdLeaseClose()", leaseId, "error", err.Error())
			}
			log.Info("server.etcd", "cancelLoopCh")
			return
		case _, ok := <-leaseKeepAliveResponseChan: //一般情况下，都是间隔500ms一次心跳
			if ok {
				log.Info("leaseKeepAliveResponseChan", ok)
				continue
			}
			//到这里说明已经收不到心跳了,那就关闭旧的租约，开启一个新的租约
			if err = s.etcdLeaseClose(leaseId); err != nil {
				log.Error("etcd close lease err ", err, " lease id ", leaseId)
				return
			}
			time.Sleep(time.Second)
			ch, id, err := s.registerServiceToEtcd()
			if err != nil {
				log.Error("Server.registerServiceToEtcd()", err.Error())
				continue
			}
			leaseKeepAliveResponseChan, leaseId = ch, id
		}
	}
}

func (s *Server) etcdLeaseClose(id etcd.LeaseID) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.etcdClient == nil {
		return nil
	}
	//请求 etcd 撤销租约
	ctx, cancel := context.WithTimeout(s.etcdClient.Ctx(), time.Second*5)
	defer cancel()
	_, err := s.etcdClient.Revoke(ctx, id)
	if err != nil && err != context.DeadlineExceeded && err != context.Canceled && err != etcdrpctypes.ErrLeaseNotFound {
		log.Error("etcd revoke err ", err.Error(), " etcd lease id ", id)
	}

	//请求 etcd 关闭租约
	err = s.etcdClient.Close()
	if err != nil {
		return err
	}
	s.etcdClient = nil
	return nil
}

func (s *Server) registerServiceToEtcd() (res <-chan *etcd.LeaseKeepAliveResponse, leaseId etcd.LeaseID, err error) {
	addr := ""
	defer func() {
		log.Info("registerServiceToEtcd", err, " leaseId", leaseId, " addr", addr)
	}()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	//如果etcd租约不为nil，先关闭，再置为nil
	if s.etcdLease != nil {
		if err = s.etcdLease.Close(); err != nil {
			return nil, 0, err
		}
		s.etcdLease = nil
	}
	//生成新的 etcd 租约
	s.etcdLease = etcd.NewLease(s.etcdClient)
	ctx, cancel := context.WithTimeout(s.etcdClient.Ctx(), time.Second*5)
	defer cancel()
	leaseGrantResponse, err := s.etcdLease.Grant(ctx, ServiceLeaseTTL)
	if err != nil {
		return nil, 0, err
	}
	//获取grpc名词解析器
	resolver := etcdnaming.GRPCResolver{Client: s.etcdClient}
	addr = fmt.Sprintf("%s:%d", s.host, s.port)
	idStr := fmt.Sprintf("%v", leaseGrantResponse.ID)
	update := grpcnaming.Update{
		Op:       grpcnaming.Add,
		Addr:     addr,
		Metadata: idStr,
	}

	//这里的 target 是 root/service/xxx(服务名)/leaseId
	target := "root/service/" + s.target + "/" + idStr
	if err = resolver.Update(ctx, target, update, etcd.WithLease(leaseGrantResponse.ID)); err != nil {
		return nil, 0, err
	}
	if res, err = s.etcdLease.KeepAlive(context.Background(), leaseGrantResponse.ID); err != nil {
		return nil, 0, err
	}
	return res, leaseGrantResponse.ID, nil
}

func (s *Server) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.grpcServer.RegisterService(sd, ss)
}

func (s *Server) GetServiceInfo() map[string]grpc.ServiceInfo {
	return s.grpcServer.GetServiceInfo()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.grpcServer.ServeHTTP(w, r)
}

func (s *Server) GracefulStop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.stopped {
		return
	}
	s.stopped = true
	close(s.closeNotifyCh)
	if s.etcdLease != nil {
		if err := s.etcdLease.Close(); err != nil {
			log.Error("Server.etcdLease.Close()", err.Error())
		} else {
			s.etcdLease = nil
		}
	}
	s.grpcServer.GracefulStop()
}
