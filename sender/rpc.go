package sender

import (
	"context"
	"errors"
	"facechat/config"
	connect_proto "facechat/proto/connect"
	message_proto "facechat/proto/message"
	"facechat/tools"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var once2 sync.Once

type Instance struct {
	ServerType string
	ServerId   string
	Client     connect_proto.ConnectClient
}

type RpcConnectClient struct {
	lock         sync.Mutex
	ServerInsMap map[string]Instance
}

var RClient = &RpcConnectClient{
	ServerInsMap: make(map[string]Instance),
}

var MessageLogicClient message_proto.MessageClient

func (rc *RpcConnectClient) GetRpcClientByServerId(serverId string) (c connect_proto.ConnectClient, err error) {
	rc.lock.Lock()
	defer rc.lock.Unlock()
	logrus.Infof("GetRpcClientByServerId [%s] => %s", serverId, rc.ServerInsMap[serverId])
	if _, ok := rc.ServerInsMap[serverId]; !ok {
		return nil, errors.New("no connect layer ip : " + serverId)
	}
	ins := rc.ServerInsMap[serverId]
	return ins.Client, nil
}

func (rc *RpcConnectClient) GetAllConnectRpcClient() (rpcClientList []connect_proto.ConnectClient) {
	for serverId, _ := range rc.ServerInsMap {
		c, err := rc.GetRpcClientByServerId(serverId)
		if err != nil {
			logrus.Infof("GetAllConnectRpcClient err:%s", err.Error())
			continue
		}
		rpcClientList = append(rpcClientList, c)
	}
	return
}

func getParamByKey(s string, key string) string {
	params := strings.Split(s, "&")
	for _, p := range params {
		kv := strings.Split(p, "=")
		if len(kv) == 2 && kv[0] == key {
			return kv[1]
		}
	}
	return ""
}

func (s *Sender) InitConnectRpcClient() (err error) {
	key := config.Conf.Connect.ConnectRpcAddressWebsockets.SdDir + "/" + config.Conf.Connect.ConnectRpcAddressWebsockets.SdName + "/" + config.Conf.Connect.ConnectRpcAddressWebsockets.SdVersion

	ServiceList, err := tools.GetAllServiceFromRepository(key)
	if err != nil {
		logrus.Info("GetAllServiceFromRepository get List error", err)
	}

	for index, service := range ServiceList {
		logrus.Infof("Addr : %s, Name: %s, InstanceId: %s, index: %d", service.Address, service.InstanceId, service.InstanceId, index)
		serverType := getParamByKey(service.InstanceId, "serverType")
		serverId := getParamByKey(service.InstanceId, "serverId")
		logrus.Infof("serviceName:%s, serverType is:%s,serverId is:%s", service.InstanceId, serverType, serverId)
		if serverType == "" || serverId == "" {
			continue
		}
		c := tools.NewRpcClientDirect(service.Address)
		client := connect_proto.NewConnectClient(c)
		ins := Instance{
			ServerType: serverType,
			ServerId:   serverId,
			Client:     client,
		}

		if _, ok := RClient.ServerInsMap[serverId]; !ok {
			RClient.ServerInsMap[serverId] = ins
		}
		logrus.Infof("RClient.ServerInsMap => %s", RClient.ServerInsMap)
	}
	go s.watchServiceChange(key)
	return
}

func (s *Sender) watchServiceChange(key string) {
	for {
		time.Sleep(30 * time.Second)
		ServiceList, err := tools.WatchServiceChanged(key)
		if err != nil {
			continue
		}
		for _, service := range ServiceList {
			serverType := getParamByKey(service.InstanceId, "serverType")
			serverId := getParamByKey(service.InstanceId, "serverId")
			logrus.Infof("serviceName:%s, serverType is:%s,serverId is:%s", service.InstanceId, serverType, serverId)
			if serverType == "" || serverId == "" {
				continue
			}
			c := tools.NewRpcClientDirect(service.Address)
			client := connect_proto.NewConnectClient(c)
			ins := Instance{
				ServerType: serverType,
				ServerId:   serverId,
				Client:     client,
			}
			RClient.lock.Lock()
			RClient.ServerInsMap[serverId] = ins
			RClient.lock.Unlock()
		}
	}
}

func (s *Sender) InitMessageLogicRpcClient() (err error) {
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

func (s *Sender) pushSingleToConnect(serverId string, userId string, msg []byte) {
	logrus.Infof("pushsingleToConnect serverId : %s, Body: %s", serverId, string(msg))
	if serverId == "" {
		// unread message logic needed
		s.unReadMsgSave(userId, msg)
		return
	}
	pushMsgReq := connect_proto.PushMsgReq{
		UserId: userId,
		Msg: &connect_proto.Msg{
			Ver:       config.MsgVersion,
			Operation: config.OpSingleSend,
			Seq:       tools.GetSnowflakeId(),
			Body:      msg,
		},
	}
	connectRpc, err := RClient.GetRpcClientByServerId(serverId)
	logrus.Infof("connectRpc : %s, err : %s", connectRpc, err)
	if err != nil {
		logrus.Infof("get rpc client err %v", err)
		return
	}
	resp, err := connectRpc.Push(context.Background(), &pushMsgReq)
	if err != nil {
		logrus.Infof("Push Call err %v", err)
		return
		// save to unread message to redis
	}
	logrus.Infof("reply %s", resp.Msg)
}

func (s *Sender) unReadMsgSave(userId string, msg []byte) {
	req := message_proto.SaveUnReadMsgReq{
		UserId: userId,
		Msg:    string(msg),
	}
	resp, err := MessageLogicClient.SaveUnReadMsg(context.Background(), &req)
	if err != nil {
		logrus.Infof("unReadMsgSave Call err %v", err)
	}
	logrus.Infof("reply %s", resp.Ok)
}
