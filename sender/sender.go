package sender

import (
	"facechat/config"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Sender struct {
}

func New() *Sender {
	return new(Sender)
}

func (s *Sender) Run() {
	senderConfig := config.Conf.Sender
	runtime.GOMAXPROCS(senderConfig.SenderBase.CpuNum)

	if err := s.InitQueueRedisClient(); err != nil {
		logrus.Panicf("sender init InitQueueRedisClient fail, err:%s", err.Error())
	}
	if err := s.InitConnectRpcClient(); err != nil {
		logrus.Panicf("sender init InitConnectRpcClient fail, err:%s", err.Error())
	}
	if err := s.InitMessageLogicRpcClient(); err != nil {
		logrus.Panicf("InitMessageLogicRpcClient err:%s", err)
	}
	s.GoPush()
}
