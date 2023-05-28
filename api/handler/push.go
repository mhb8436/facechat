package handler

import (
	"facechat/api/rpc"
	"facechat/config"
	message_proto "facechat/proto/message"
	"facechat/tools"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type FormPush struct {
	Msg      string `form:"msg" json:"msg" binding:"required"`
	ToUserId string `form:"toUserId" json:"toUserId" binding:"required"`
	RoomId   string `form:"roomId" json:"roomId" binding:"required"`
	// AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func Send(c *gin.Context) {
	var formPush FormPush
	if err := c.ShouldBindBodyWith(&formPush, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	bearerToken := c.Request.Header.Get("Authorization")
	var accessToken string
	if len(strings.Split(bearerToken, " ")) == 2 {
		accessToken = strings.Split(bearerToken, " ")[1]
	}
	msg := formPush.Msg
	roomId := formPush.RoomId
	req := &message_proto.SendReq{
		Code:        config.OpSingleSend,
		Msg:         msg,
		RoomUuid:    roomId,
		AccessToken: accessToken,
	}
	logrus.Infof("api send %s", msg)
	ok := rpc.MessageLogicObj.Send(req)
	logrus.Infof("api send %s result %s", msg, ok)
	var jsonData = map[string]interface{}{
		"ok": ok,
	}
	tools.SuccessWithMsg(c, "join room success", jsonData)

}
