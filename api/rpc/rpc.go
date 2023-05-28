package rpc

import (
	"context"
	"facechat/config"
	message_proto "facechat/proto/message"
	"facechat/tools"
	"sync"

	"github.com/sirupsen/logrus"
)

var once sync.Once
var once2 sync.Once
var MessageLogicClient message_proto.MessageClient

type UserInfo map[string]interface{}

type UserRpcLogic struct {
}

type MessageRpcLogic struct {
}

var UserRpcLogicObj *UserRpcLogic
var MessageLogicObj *MessageRpcLogic

func InitMessageLogicRpcClient() {
	once2.Do(func() {
		logrus.Infoln("InitMessageLogicRpcClient begin")
		c := tools.NewRpcClient(
			config.Conf.Message.MessageBase.SdName,
			config.Conf.Message.MessageBase.SdVersion,
			config.Conf.Message.MessageBase.SdDir,
		)
		// c := tools.NewRcpTestClient()
		logrus.Infoln("InitMessageLogicRpcClient get con", c)
		MessageLogicClient = message_proto.NewMessageClient(c)
		MessageLogicObj = new(MessageRpcLogic)
		if MessageLogicClient == nil {
			logrus.Fatalf("get message logic client nil")
		}
		logrus.Infoln("InitMessageLogicRpcClient end")
	})
	if MessageLogicClient == nil {
		logrus.Fatalf("get message logic rpc client nil")
	}
}

func (rpc *MessageRpcLogic) Send(req *message_proto.SendReq) (ok bool) {
	resp, err := MessageLogicClient.Send(context.Background(), req)
	logrus.Infoln("MessageLogicClient", "Send", "result", resp, err)
	if err != nil {
		return false
	}
	return true
}
