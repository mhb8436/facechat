package proto

type RedisMsg struct {
	Op           int               `json:"op"`
	ServerId     string            `json:"serverId,omitempty"`
	RoomId       string            `json:"roomId,omitempty"`
	FromUserId   string            `json:"fromUserId,omitempty"`
	UserId       string            `json:"userId,omitempty"`
	Msg          []byte            `json:"msg"`
	Count        int               `json:"count"`
	RoomUserInfo map[string]string `json:"roomUserInfo"`
}

type SuccessReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type UnReadMsg struct {
	UserId string `json:"userId,omitempty"`
	Msg    []byte `json:"msg"`
}
