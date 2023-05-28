package connect

import (
	"encoding/json"
	message_proto "facechat/proto/message"
	"facechat/tools"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Buckets   []*Bucket
	Options   ServerOptions
	bucketIdx uint32
	operator  Operator
}

type ServerOptions struct {
	WriteWait       time.Duration
	PongWait        time.Duration
	PingPeriod      time.Duration
	MaxMessageSize  int64
	ReadBufferSize  int
	WriteBufferSize int
	BroadcastSize   int
}

type Method string

const (
	NewReq Method = "new"
	ByeReq Method = "bye"
)

type Request struct {
	Type Method      `json:"type"`
	Data interface{} `json:"data"`
}

type NewReqBody struct {
	AccessToken string `json:"access_token"`
	ServerId    string `json:"server_id"`
}

func NewServer(b []*Bucket, o Operator, options ServerOptions) *Server {
	s := new(Server)
	s.Buckets = b
	s.Options = options
	s.bucketIdx = uint32(len(b))
	s.operator = o
	return s
}

func (s *Server) Bucket(userId string) *Bucket {
	userIdStr := fmt.Sprintf("%s", userId)
	idx := tools.CityHash32([]byte(userIdStr), uint32(len(userIdStr))) % s.bucketIdx
	return s.Buckets[idx]
}

func (s *Server) writePump(ch *Channel, c *Connect) {
	ticker := time.NewTicker(s.Options.PingPeriod)
	defer func() {
		ticker.Stop()
		ch.conn.Close()
	}()

	for {
		select {
		case message, ok := <-ch.broadcast:
			ch.conn.SetWriteDeadline(time.Now().Add(s.Options.WriteWait))
			if !ok {
				logrus.Warn("SetWriteDeadline not ok")
				ch.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := ch.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logrus.Warn("ch.conn.NextWriter err:%s", err.Error())
				return
			}
			logrus.Infof("message write body:%s", message.Body)
			w.Write(message.Body)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			ch.conn.SetWriteDeadline(time.Now().Add(s.Options.WriteWait))
			logrus.Infof("websocket.PingMessage :%v", websocket.PingMessage)
			if err := ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *Server) readPump(ch *Channel, c *Connect) {
	defer func() {
		logrus.Infof("start exec disConnect....")
		if ch.userId == "-" {
			logrus.Infof("userId is 0")
			ch.conn.Close()
			return
		}
		logrus.Infof("exec disConnect...")
		// disConnect to rpc
		disConnectReq := &message_proto.DisConnectReq{
			UserId: ch.userId,
		}
		s.Bucket(ch.userId).DeleteChannel(ch)
		if err := s.operator.DisConnect(disConnectReq); err != nil {
			logrus.Warnf("DisConnect err:%s", err.Error())
		}
		ch.conn.Close()
	}()
	ch.conn.SetReadLimit(s.Options.MaxMessageSize)
	ch.conn.SetReadDeadline(time.Now().Add(s.Options.PongWait))
	ch.conn.SetPongHandler(func(string) error {
		ch.conn.SetReadDeadline(time.Now().Add(s.Options.PongWait))
		return nil
	})

	for {
		_, message, err := ch.conn.ReadMessage()
		logrus.Infoln("readMessage", message, err)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("readPump ReadMessage err:%s", err.Error())
				return
			}
		}
		logrus.Infoln("readMessage is nil ", message == nil)
		if message == nil {
			return
		}
		logrus.Infof("get a message:%s", message)
		var body json.RawMessage
		request := &Request{
			Data: &body,
		}
		if err := json.Unmarshal([]byte(message), &request); err != nil {
			logrus.Errorf("message struct err:%s", err.Error())
			return
		}
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			logrus.Errorf("json unmarshal err:%v", err)
			return
		}
		logrus.Infof("request type is ", request.Type)
		if ch.userId != "-" {
			b := s.Bucket(ch.userId)
			err = b.Put(ch.userId, ch)
			if err != nil {
				logrus.Errorf("conn open err:%s", err.Error())
				ch.conn.Close()
			}
		}
		switch request.Type {
		case NewReq: // this code maybe will be changed because webrtc singling will be added
			var newBody *NewReqBody
			if err := json.Unmarshal(body, &newBody); err != nil {
				logrus.Errorf("New req unmarshal err:%v", err)
				return
			}
			if newBody == nil || newBody.AccessToken == "" {
				logrus.Errorf("s.operator.Connect has no AccessToken ")
				return
			}
			connectReq := message_proto.ConnectReq{
				AccessToken: newBody.AccessToken,
				// ServerId:    newBody.ServerId,
				ServerId: c.ServerId,
			}
			userId, err := s.operator.Connect(&connectReq)
			if err != nil {
				logrus.Errorf("s.operator.Connect err:%v", err.Error())
				return
			}
			if userId != "-" {
				logrus.Errorf("[S]Invalid Access Token, userId empty!")
				return
			}
			logrus.Infof("websocket rpc call return userId:%d", userId)
			b := s.Bucket(userId)
			err = b.Put(userId, ch)
			if err != nil {
				logrus.Errorf("conn open err:%s", err.Error())
				ch.conn.Close()
			}
		}

	}
}
