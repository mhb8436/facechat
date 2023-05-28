package tools

import (
	"facechat/config"
	// "fmt"

	"time"

	lb_common "github.com/mhb8436/grpc-lb/common"
	"github.com/mhb8436/grpc-lb/registry"
	etcd "github.com/mhb8436/grpc-lb/registry/etcd3"
	"github.com/sirupsen/logrus"
	etcdv3 "go.etcd.io/etcd/client/v3"

	// resolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// type RpcServer struct {
// 	addr string
// 	S    *grpc.Server
// }

// func NewRpcServer(addr string, r interface{}) interface{} {
// 	s := grpc.NewServer()
// 	r = &RpcServer{
// 		addr: addr,
// 		S:    s,
// 	}
// 	return r
// }

// func (s *RpcServer) Run() {
// 	listener, err := net.Listen("tcp", s.addr)
// 	if err != nil {
// 		logrus.Panicf("RcpServer Run Fail err:%s", err)
// 	}
// 	logrus.Printf("rcp listening on:%s", s.addr)
// 	s.S.Serve(listener)
// }

// func (s *RpcServer) Stop() {
// 	s.S.GracefulStop()
// }

func AddRegistry(nodeID string, sdName string, sdVersion string, sdDir string, addr string) *etcd.Registrar {
	etcdConfig := etcdv3.Config{
		Endpoints: []string{config.Conf.Common.CommonEtcd.Host},
	}

	service := &registry.ServiceInfo{
		InstanceId: nodeID,
		Name:       sdName,
		Version:    sdVersion,
		Address:    addr,
		Metadata:   metadata.Pairs(lb_common.WeightKey, "1"),
	}

	registrar, err := etcd.NewRegistrar(
		&etcd.Config{
			EtcdConfig:  etcdConfig,
			RegistryDir: sdDir,
			Ttl:         10 * time.Second,
		})
	if err != nil {
		logrus.Panicf("error:%s", err)
		return nil
	}

	go func() {
		registrar.Register(service)
	}()
	return registrar
}

func NewRpcClient(sdName string, sdVersion string, sdDir string) *grpc.ClientConn {
	logrus.Infoln("NewRpcClient", sdName, sdVersion, sdDir)
	etcdConfig := etcdv3.Config{
		Endpoints: []string{config.Conf.Common.CommonEtcd.Host},
	}

	etcd.RegisterResolver("etcd3", etcdConfig, sdDir, sdName, sdVersion)
	logrus.Infoln("NewRpcClient", "etcd.RegisterResolver")
	c, err := grpc.Dial("etcd3:///", grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	logrus.Infoln("NewRpcClient", "grpc Connection", c)
	if err != nil {
		logrus.Infof("grpc dial error : %s", err)
	}
	return c
}

func NewRpcClientDirect(addr string) *grpc.ClientConn {
	c, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Infof("grpc dial error : %s", err)
	}
	logrus.Infof("NewRpcClientDirect connected => %s", c)
	return c
}

func NewRcpTestClient() *grpc.ClientConn {
	c, err := grpc.Dial("localhost:6900", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Infof("grpc dial error : %s", err)
	}
	logrus.Infof("NewRcpTestClient connected", c)
	return c
}
