package message

import (
	"facechat/config"
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MessageLogic struct {
	ServerId string
}

func New() *MessageLogic {
	return new(MessageLogic)
}

func (messageLogic *MessageLogic) Run() {
	messageLogicConfig := config.Conf.Message

	runtime.GOMAXPROCS(messageLogicConfig.MessageBase.CpuNum)
	messageLogic.ServerId = fmt.Sprintf("message-logic-%s", uuid.New().String())
	if err := messageLogic.InitPublishRedisClient(); err != nil {
		logrus.Panicf("message logic publishRedisClient fail, err:%s", err.Error())
	}

	if err := messageLogic.InitRpcServer(); err != nil {
		logrus.Panicf("user logic rpc server fail, err:%s", err.Error())
	}
}
