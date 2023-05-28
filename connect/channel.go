package connect

import "github.com/gorilla/websocket"

type ChannelMsg struct {
	Ver       int
	Operation string
	SeqId     string
	Body      []byte
}

type Channel struct {
	broadcast chan *ChannelMsg
	userId    string
	conn      *websocket.Conn
}

func NewChannel(size int) (c *Channel) {
	c = new(Channel)
	c.broadcast = make(chan *ChannelMsg, size)
	return
}

func (ch *Channel) Push(msg *ChannelMsg) (err error) {
	select {
	case ch.broadcast <- msg:
	default:
	}
	return
}
