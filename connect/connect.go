package connect

import (
	"facechat/config"
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var DefaultServer *Server

type Connect struct {
	ServerId string
}

func New() *Connect {
	return new(Connect)
}

func (c *Connect) Run() {
	connectConfig := config.Conf.Connect

	runtime.GOMAXPROCS(connectConfig.ConnectBucket.CpuNum)

	if err := c.InitMessageLogicRpcClient(); err != nil {
		logrus.Panicf("InitMessageLogicRpcClient err:%s", err)
	}
	Buckets := make([]*Bucket, connectConfig.ConnectBucket.CpuNum)
	for i := 0; i < connectConfig.ConnectBucket.CpuNum; i++ {
		Buckets[i] = NewBucket(BucketOptions{
			ChannelSize: connectConfig.ConnectBucket.ChannelSize,
		})
	}
	operator := new(DefaultOperator)
	DefaultServer = NewServer(Buckets, operator, ServerOptions{
		WriteWait:       time.Duration(config.Conf.Connect.ConnectWebsocket.WriteWait).Abs() * time.Second,
		PongWait:        time.Duration(config.Conf.Connect.ConnectWebsocket.PongWait).Abs() * time.Second,
		PingPeriod:      time.Duration(config.Conf.Connect.ConnectWebsocket.PingPeriod).Abs() * time.Second,
		MaxMessageSize:  int64(config.Conf.Connect.ConnectWebsocket.MaxMsgSize),
		ReadBufferSize:  config.Conf.Connect.ConnectWebsocket.ReadBufSize,
		WriteBufferSize: config.Conf.Connect.ConnectWebsocket.WriteBufSize,
		BroadcastSize:   config.Conf.Connect.ConnectWebsocket.BroadcastSize,
	})
	c.ServerId = fmt.Sprintf("%s-%s", "ws", uuid.New().String())
	logrus.Infof("Connect Run ServerId : %s", c.ServerId)
	if err := c.InitConnectWebsocketRpcServer(); err != nil {
		logrus.Panicf("InitConnectWebsocketRpcServer Fatal error: %s \n", err.Error())
	}
	if err := c.InitWebsocket(); err != nil {
		logrus.Panicf("Connect layer InitWebsocket() error: %s \n", err.Error())
	}

}
