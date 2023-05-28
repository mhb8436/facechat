package connect

import (
	"facechat/config"
	message_proto "facechat/proto/message"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (c *Connect) InitWebsocket() error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c.serveWs(DefaultServer, w, r)
	})
	err := http.ListenAndServe(config.Conf.Connect.ConnectWebsocket.Bind, nil)
	return err
}

func (c *Connect) serveWs(server *Server, w http.ResponseWriter, r *http.Request) {
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  server.Options.ReadBufferSize,
		WriteBufferSize: server.Options.WriteBufferSize,
	}

	upGrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("serveWs err:%s", err.Error())
		return
	}

	var ch = NewChannel(server.Options.BroadcastSize)
	ch.conn = conn

	// create channel and login
	bearerToken := r.Header.Get("Authorization")
	var accessToken string
	if len(strings.Split(bearerToken, " ")) == 2 {
		accessToken = strings.Split(bearerToken, " ")[1]
	}
	connectReq := message_proto.ConnectReq{
		AccessToken: accessToken,
		ServerId:    c.ServerId,
	}
	logrus.Infof("serveWs serverId : %s", c.ServerId)
	userId, err := server.operator.Connect(&connectReq)
	logrus.Infof("serveWs userId : %s", userId)
	ch.userId = userId
	if err != nil {
		logrus.Errorf("s.operator.Connect err:%v", err.Error())
		return
	}
	if userId == "-" {
		logrus.Errorf("[W]Invalid Access Token, userId empty!")
		return
	}
	logrus.Infof("websocket rpc call return userId:%d", userId)
	b := server.Bucket(userId)
	err = b.Put(userId, ch)
	if err != nil {
		logrus.Errorf("conn open err:%s", err.Error())
		ch.conn.Close()
	}

	// send data
	go server.writePump(ch, c)
	// receive data
	go server.readPump(ch, c)
	// unread mssage send
	time.Sleep(2)
	go c.SendUnReadMsg(ch.userId)
}
