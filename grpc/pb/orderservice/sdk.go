package orderservice

import (
	"errors"

	myEtcd "github.com/manmanxing/go_common/beacon/etcd"
	myGrpc "github.com/manmanxing/go_common/grpc"
	"google.golang.org/grpc"
)

const serverName = "order_server"

//client
func Client() (OrderServiceClient, error) {
	conn, err := myGrpc.ClientConn(serverName)
	if err != nil {
		return nil, err
	}
	return NewOrderServiceClient(conn), nil
}

// 访问自定义的client
func NewClient(etcdEndpoints []string) (OrderServiceClient, error) {
	etcdClient, err := myEtcd.GetClient(etcdEndpoints)
	if err != nil {
		return nil, err
	}
	conn, err := myGrpc.NewClientConn(etcdClient, serverName)
	if err != nil {
		return nil, err
	}
	return NewOrderServiceClient(conn), nil
}

//server
func Server(port int, srv OrderServiceServer) error {
	if port <= 0 {
		return errors.New("port cant <= 0")
	}

	if srv == nil {
		return errors.New("srv cant be nil")
	}

	newServer := register(port, srv)
	return newServer.Serve()
}

func register(port int, srv OrderServiceServer) *myGrpc.Server {
	r := func(s *grpc.Server) {
		RegisterOrderServiceServer(s, srv)
	}
	return myGrpc.NewServer(serverName, "", port, r)
}
