package connect

import (
	message_proto "facechat/proto/message"

	"github.com/sirupsen/logrus"
)

type Operator interface {
	Connect(conn *message_proto.ConnectReq) (string, error)
	DisConnect(conn *message_proto.DisConnectReq) (err error)
}

type DefaultOperator struct {
	Operator
}

func (o *DefaultOperator) Connect(conn *message_proto.ConnectReq) (uid string, err error) {
	rpcConnect := new(RpcConnect)
	uid, err = rpcConnect.Connect(conn)
	return
}

func (o *DefaultOperator) DisConnect(disConn *message_proto.DisConnectReq) (err error) {
	rpcConnect := new(RpcConnect)
	has, err := rpcConnect.DisConnect(disConn)
	logrus.Infoln("DisConnect result:%s", has)
	return
}
