package sender

import (
	"encoding/json"
	"facechat/config"
	"facechat/proto"
	"math/rand"

	"github.com/sirupsen/logrus"
)

type PushParams struct {
	ServerId string
	UserId   string
	Msg      []byte
	RoomId   string
}

var pushChannel []chan *PushParams

func init() {
	pushChannel = make([]chan *PushParams, config.Conf.Sender.SenderBase.PushChan)
}

func (s *Sender) GoPush() {
	for i := 0; i < len(pushChannel); i++ {
		pushChannel[i] = make(chan *PushParams, config.Conf.Sender.SenderBase.PushChan)
		go s.processSinglePush(pushChannel[i])
	}
}

func (s *Sender) processSinglePush(ch chan *PushParams) {
	var arg *PushParams
	for {
		arg = <-ch
		s.pushSingleToConnect(arg.ServerId, arg.UserId, arg.Msg)
	}
}

func (s *Sender) Push(msg string) {
	logrus.Infof("Push orgin msg => " + msg)
	m := &proto.RedisMsg{}
	if err := json.Unmarshal([]byte(msg), m); err != nil {
		logrus.Infof("json.Unmarshal err:%v", err)
	}
	logrus.Infof("push msg info userId : %d, op : %d, msg: %s", m.UserId, m.Op, m.Msg)
	switch m.Op {
	case config.OpSingleSend:
		pushChannel[rand.Int()%config.Conf.Sender.SenderBase.PushChan] <- &PushParams{
			ServerId: m.ServerId,
			UserId:   m.UserId,
			Msg:      m.Msg,
		}
	}
}
