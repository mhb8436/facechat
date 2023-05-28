package message

import (
	"context"
	"encoding/json"
	"facechat/config"
	message_proto "facechat/proto/message"
	"facechat/tools"

	// user "facechat/user"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *MessageRpcServer) Connect(ctx context.Context, req *message_proto.ConnectReq) (*message_proto.ConnectResp, error) {
	logrus.Infoln("UserRpcServer Connect", req)
	u := new(CustomUser)
	u.Access_token = req.AccessToken
	oUser, err := u.FindUserByAccessToken()
	if err != nil {
		logrus.Infof("Connect AccessToken is invalid err=>%s", err)
	}
	oUser.UpdateConnect()
	logrus.Infof("Connect Ouser => %s", oUser.Uuid)
	// save redis userId:serverId
	userKey := tools.GetUserKey(fmt.Sprintf("%s", oUser.Uuid))
	validTime := config.RedisBaseValidTime * time.Second
	serverId := req.ServerId
	err = RedisClient.Set(userKey, serverId, validTime).Err()
	if err != nil {
		logrus.Warnf("userId:ServerId set err:%s", err)
	}
	return &message_proto.ConnectResp{
		UserId: oUser.Uuid,
	}, nil
}

func (s *MessageRpcServer) DisConnect(ctx context.Context, req *message_proto.DisConnectReq) (*message_proto.DisConnectResp, error) {
	userKey := tools.GetUserKey(fmt.Sprintf("%s", req.UserId))
	if err := RedisClient.Del(userKey); err != nil {
		logrus.Warnf("userId:ServerId del err:%s", err)
	}
	u := new(CustomUser)
	u.Uuid = req.UserId
	u.UpdateDisConnect()

	return &message_proto.DisConnectResp{
		Has: false,
	}, nil

}

func (s *MessageRpcServer) Send(ctx context.Context, req *message_proto.SendReq) (*message_proto.SendResp, error) {
	u := new(CustomUser)
	u.Access_token = req.AccessToken
	oUser, err := u.FindUserByAccessToken()
	if err != nil {
		logrus.Infof("Connect AccessToken is invalid err=>%s", err)
	}
	var bodyBytes []byte
	bodyBytes, err = json.Marshal(req.Msg)
	// fromUserKey := tools.GetUserKey(fmt.Sprintf("%d", oUser.Id))
	c := new(Chat)
	c.Uuid = req.RoomUuid
	logic := new(MessageLogic)
	userList, err := c.FindAllUserInRoom()
	for _, user := range userList {
		// if user.Id == oUser.Id {
		// 	continue
		// }
		toUserKey := tools.GetUserKey(fmt.Sprintf("%s", user.Uuid))
		serverIdStr := RedisClient.Get(toUserKey).Val()
		err = logic.RedisPublishChannel(serverIdStr, oUser.Uuid, user.Uuid, req.RoomUuid, bodyBytes)
		if err != nil {
			logrus.Errorf("logic,redis publish err: %s", err.Error())
			continue
		}
	}
	return &message_proto.SendResp{
		Ok: true,
	}, nil
}

func (s *MessageRpcServer) SaveUnReadMsg(ctx context.Context, req *message_proto.SaveUnReadMsgReq) (*message_proto.SaveUnReadMsgResp, error) {
	logic := new(MessageLogic)
	var bodyBytes []byte
	bodyBytes, err := json.Marshal(req.Msg)
	if err != nil {
		logrus.Infof("marshal err : %v", err)
	}
	err = logic.RedisUnReadMsgSave(req.UserId, bodyBytes)
	if err != nil {
		logrus.Infof("RedisUnReadMsgSave err : %v", err)
	}
	return &message_proto.SaveUnReadMsgResp{
		Ok: true,
	}, nil
}

func (s *MessageRpcServer) GetUnReadMsg(ctx context.Context, req *message_proto.GetUnReadMsgReq) (*message_proto.GetUnReadMsgResp, error) {
	logic := new(MessageLogic)
	results, err := logic.RedisUnReadMsgGet(req.UserId)
	if err != nil {
		logrus.Infof("unread msg readed err:%v", err)
	}
	return &message_proto.GetUnReadMsgResp{
		Msg: results,
	}, nil

}
