package message

import (
	"encoding/json"
	"facechat/config"
	"facechat/proto"
	message_proto "facechat/proto/message"
	"facechat/tools"
	"fmt"
	"net"
	"strings"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/mhb8436/grpc-lb/registry"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var RedisClient *redis.Client
var RedisSessClient *redis.Client

type MessageRpcServer struct {
	message_proto.UnimplementedMessageServer
	addr string
	S    *grpc.Server
}

func (logic *MessageLogic) InitPublishRedisClient() (err error) {
	redisOpt := tools.RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}
	RedisClient = tools.GetRedisInstance(redisOpt)
	if pong, err := RedisClient.Ping().Result(); err != nil {
		logrus.Infof("RedisCli Ping Result pong: %s, err: %s", pong, err)
	}
	RedisSessClient = RedisClient
	return err
}

func (logic *MessageLogic) InitRpcServer() (err error) {
	var network, addr string
	rpcAddressList := strings.Split(config.Conf.Message.MessageBase.RpcAddress, ",")
	for _, bind := range rpcAddressList {
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitMessageLogicRpc PaseNetwork err: %s", err)
		}
		logrus.Infof("Message logic start run at-> %s:%s", network, addr)
		go logic.createRpcServer(network, addr)
	}
	return nil
}

func (logic *MessageLogic) createRpcServer(network string, addr string) {
	logrus.Infoln("createRpcServer begin %s:%s", network, addr)
	InstanceId := uuid.New().String()
	registrar := tools.AddRegistry(
		InstanceId,
		config.Conf.Message.MessageBase.SdName,
		config.Conf.Message.MessageBase.SdVersion,
		config.Conf.Message.MessageBase.SdDir,
		addr,
	)
	defer registrar.Unregister(&registry.ServiceInfo{
		InstanceId: InstanceId,
	})
	logic.runRpcServer(network, addr)
}

func (logic *MessageLogic) runRpcServer(network string, addr string) {
	s := grpc.NewServer()
	rs := &MessageRpcServer{
		addr: addr,
		S:    s,
	}
	listener, err := net.Listen(network, rs.addr)
	if err != nil {
		logrus.Panicf("failed to listen: %v", err)
	}
	defer rs.S.GracefulStop()
	message_proto.RegisterMessageServer(rs.S, rs)
	rs.S.Serve(listener)
	logrus.Infof("rpc listening on:%s", rs.addr)
}

func (logic *MessageLogic) RedisPublishChannel(serverId string, fromUserId string, toUserId string, roomId string, msg []byte) (err error) {
	redisMsg := &proto.RedisMsg{
		Op:         config.OpSingleSend,
		ServerId:   serverId,
		FromUserId: fromUserId,
		UserId:     toUserId,
		RoomId:     roomId,
		Msg:        msg,
	}
	redisMsgStr, err := json.Marshal(redisMsg)
	if err != nil {
		logrus.Errorf("logic.RedisPublishChannel Marshal err:%s", err.Error())
		return err
	}
	redisChannel := config.QueueName
	if err := RedisClient.LPush(redisChannel, redisMsgStr).Err(); err != nil {
		logrus.Errorf("logic,lpush err:%s", err.Error())
		return err
	}
	return
}

func (logic *MessageLogic) RedisUnReadMsgSave(userId string, msg []byte) (err error) {
	unReadMsg := &proto.UnReadMsg{
		UserId: userId,
		Msg:    msg,
	}
	unReadMsgStr, err := json.Marshal(unReadMsg)
	if err != nil {
		logrus.Errorf("logic.RedisUnReadMsgSave Marshal err:%s", err.Error())
		return err
	}
	redisKey := fmt.Sprintf("%s-%d", config.UnReadQueuePrefix, userId)
	if err := RedisClient.LPush(redisKey, unReadMsgStr).Err(); err != nil {
		logrus.Errorf("logic, unread, lpush err:%s", err.Error())
		return err
	}
	return
}

func (logic *MessageLogic) RedisUnReadMsgGet(userId string) (msg []string, err error) {
	redisKey := fmt.Sprintf("%s-%d", config.UnReadQueuePrefix, userId)
	results, err := RedisClient.LRange(redisKey, 0, -1).Result()

	if err != nil {
		logrus.Infof("redis unread message readed  err : %v", err)
	}

	for _, result := range results {
		m := &proto.RedisMsg{}
		if err := json.Unmarshal([]byte(result), m); err != nil {
			logrus.Infof("json.Unmarshal err:%v", err)
		}
		msg = append(msg, string(m.Msg))
	}
	RedisClient.Del(redisKey)
	return
}
