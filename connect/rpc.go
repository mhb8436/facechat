package connect

import (
	"context"
	"facechat/config"
	connect_proto "facechat/proto/connect"
	message_proto "facechat/proto/message"
	"facechat/tools"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/mhb8436/grpc-lb/registry"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var once sync.Once
var once2 sync.Once
var MessageLogicClient message_proto.MessageClient

type RpcConnect struct {
}

type ConnectRpcServer struct {
	connect_proto.UnimplementedConnectServer
	addr string
	S    *grpc.Server
}

func (c *Connect) InitConnectWebsocketRpcServer() (err error) {
	var network, addr string
	rpcAddressList := strings.Split(config.Conf.Connect.ConnectRpcAddressWebsockets.Address, ",")
	for _, bind := range rpcAddressList {
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitConnectWebsocketRpcServer ParseNetwork err:%s", err)
		}
		logrus.Infof("InitConnectWebsocketRpcServer start run at -> %s:%s", network, addr)
		go c.createRpcServer(network, addr)
	}
	return nil
}

func (c *Connect) createRpcServer(network string, addr string) {
	logrus.Infoln("createRpcServer begin %s:%s", network, addr)
	registrar := tools.AddRegistry(
		fmt.Sprintf("serverId=%s&serverType=ws", c.ServerId),
		config.Conf.Connect.ConnectRpcAddressWebsockets.SdName,
		config.Conf.Connect.ConnectRpcAddressWebsockets.SdVersion,
		config.Conf.Connect.ConnectRpcAddressWebsockets.SdDir,
		addr,
	)
	defer registrar.Unregister(&registry.ServiceInfo{
		InstanceId: fmt.Sprintf("serverId=%s&serverType=ws", c.ServerId),
	})
	c.runRpcServer(network, addr)
}

func (c *Connect) runRpcServer(network string, addr string) {
	s := grpc.NewServer()
	rs := &ConnectRpcServer{
		addr: addr,
		S:    s,
	}
	listener, err := net.Listen(network, rs.addr)
	if err != nil {
		logrus.Panicf("failed to listen: %v", err)
	}
	defer rs.S.GracefulStop()
	connect_proto.RegisterConnectServer(rs.S, rs)
	rs.S.Serve(listener)
	logrus.Infof("rpc listening on:%s", rs.addr)
}

func (c *Connect) InitMessageLogicRpcClient() (err error) {
	once2.Do(func() {
		logrus.Infoln("InitMessageLogicRpcClient begin")
		c := tools.NewRpcClient(
			config.Conf.Message.MessageBase.SdName,
			config.Conf.Message.MessageBase.SdVersion,
			config.Conf.Message.MessageBase.SdDir,
		)

		logrus.Infoln("InitMessageLogicRpcClient get con", c)
		MessageLogicClient = message_proto.NewMessageClient(c)

		if MessageLogicClient == nil {
			logrus.Fatalf("get message logic client nil")
		}
		logrus.Infoln("InitMessageLogicRpcClient end")
	})
	if MessageLogicClient == nil {
		logrus.Fatalf("get message logic rpc client nil")
	}
	return
}

func (rpc *RpcConnect) Connect(req *message_proto.ConnectReq) (uid string, err error) {
	resp, err := MessageLogicClient.Connect(context.Background(), req)
	logrus.Infoln("MessageLogicClient", "Connect", "result", resp, err)
	if err != nil {
		uid = "-"
		return
	}
	uid = resp.UserId
	err = nil
	return
}

func (rpc *RpcConnect) DisConnect(req *message_proto.DisConnectReq) (has bool, err error) {
	resp, err := MessageLogicClient.DisConnect(context.Background(), req)
	logrus.Infoln("MessageLogicClient", "DisConnect", "result", resp, err)
	if err != nil {
		has = false
	}
	has = resp.Has
	err = nil
	return
}

func (s *ConnectRpcServer) Push(ctx context.Context, req *connect_proto.PushMsgReq) (resp *connect_proto.SuccessResp, err error) {
	var (
		bucket  *Bucket
		channel *Channel
	)
	logrus.Info("rpc Push : %v", req)
	if req == nil {
		logrus.Errorf("rpc Push() args:(%v)", req)
		return
	}
	bucket = DefaultServer.Bucket(req.UserId)
	if channel = bucket.Channel(req.UserId); channel != nil {
		channelMsg := ChannelMsg{
			Ver:       config.MsgVersion,
			Operation: fmt.Sprintf("%d", config.OpSingleSend),
			SeqId:     tools.GetSnowflakeId(),
			Body:      req.Msg.Body,
		}
		err = channel.Push(&channelMsg)
		if err != nil {
			logrus.Errorf("occured err while push channel err : %v", err)
			return
		}
	}
	return &connect_proto.SuccessResp{
		Code: config.SuccessReplyCode,
		Msg:  config.SuccessReplyMsg,
	}, nil

}

func (c *Connect) SendUnReadMsg(userId string) (err error) {

	req := &message_proto.GetUnReadMsgReq{
		UserId: userId,
	}
	resp, err := MessageLogicClient.GetUnReadMsg(context.Background(), req)
	if err != nil {
		logrus.Infof("unread msg reaed err:%s", err.Error())
	}

	bucket := DefaultServer.Bucket(req.UserId)
	channel := bucket.Channel(req.UserId)

	for _, s := range resp.Msg {
		channelMsg := ChannelMsg{
			Ver:       config.MsgVersion,
			Operation: fmt.Sprintf("%d", config.OpSingleSend),
			SeqId:     tools.GetSnowflakeId(),
			Body:      []byte(s),
		}
		err = channel.Push(&channelMsg)
		if err != nil {
			logrus.Errorf("occured err while push channel err : %v", err)
			continue
		}
	}

	return nil
}
