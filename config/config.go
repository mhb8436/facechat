package config

import (
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once
var realPath string
var Conf *Config

type Config struct {
	Common  Common
	Api     ApiConfig
	User    UserConfig
	Connect ConnectConfig
	Sender  SenderConfig
	Message MessageConfig
}

const (
	SuccessReplyCode      = 0
	FailReplyCode         = 1
	SuccessReplyMsg       = "success"
	QueueName             = "facechat_queue"
	UnReadQueuePrefix     = "facechat_unread"
	RedisBaseValidTime    = 86400
	RedisPrefix           = "facechat_"
	RedisRoomPrefix       = "facechat_room_"
	RedisRoomOnlinePrefix = "facechat_room_online_count_"
	MsgVersion            = 1
	OpSingleSend          = 2 // single user
	OpRoomSend            = 3 // send to room
	OpRoomCountSend       = 4 // get online user count
	OpRoomInfoSend        = 5 // send info to room
	OpBuildTcpConn        = 6 // build tcp conn
)

func init() {
	Init()
}

func Init() {
	once.Do(func() {
		env := GetMode()
		realPath := getCurrentDir()
		configFilePath := realPath + "/" + env + "/"
		viper.AddConfigPath(configFilePath)
		viper.SetConfigType("toml")
		viper.SetConfigName("/api")
		err := viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("/user")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}

		viper.SetConfigName("/common")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}

		viper.SetConfigName("/connect")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}

		viper.SetConfigName("/sender")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("/message")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		Conf = new(Config)
		viper.Unmarshal(&Conf.Api)
		viper.Unmarshal(&Conf.User)
		viper.Unmarshal(&Conf.Common)
		viper.Unmarshal(&Conf.Connect)
		viper.Unmarshal(&Conf.Sender)
		viper.Unmarshal(&Conf.Message)
	})
}

func GetMode() string {
	env := os.Getenv("RUN_MODE")
	if env == "" {
		env = "dev"
	}
	return env
}

func getCurrentDir() string {
	_, fileName, _, _ := runtime.Caller(1)
	aPath := strings.Split(fileName, "/")
	dir := strings.Join(aPath[0:len(aPath)-1], "/")
	return dir
}
func GetGinRunMode() string {
	return "debug"
}

type CommonEtcd struct {
	Host              string `mapstructure:"host"`
	BasePath          string `mapstructure:"basePath"`
	ServerPathLogic   string `mapstructure:"serverPathLogic"`
	ServerPathConnect string `mapstructure:"serverPathConnect"`
	UserName          string `mapstructure:"userName"`
	Password          string `mapstructure:"password"`
	ConnectionTimeout int    `mapstructure:"connectionTimeout"`
}

type CommonRedis struct {
	RedisAddress  string `mapstructure:"redisAddress"`
	RedisPassword string `mapstructure:"redisPassword"`
	Db            int    `mapstructure:"db"`
}

type Common struct {
	CommonEtcd  CommonEtcd  `mapstructure:"common-etcd"`
	CommonRedis CommonRedis `mapstructure:"common-redis"`
}

type ApiBase struct {
	ListenPort int    `mapstructure:"listenPort"`
	ServeIp    string `mapstructure:"serveIp"`
	UploadDir  string `mapstructure:"uploadDir"`
}

type ApiConfig struct {
	ApiBase ApiBase `mapstructure:"api-base"`
}

type UserBase struct {
	CpuNum     int    `mapstructure:"cpuNum"`
	RpcAddress string `mapstructure:"rpcAddress"`
	SdName     string `mapstructure:"sdName"`
	SdVersion  string `mapstructure:"sdVersion"`
	SdDir      string `mapstructure:"sdDir"`
}

type UserConfig struct {
	UserBase UserBase `mapstructure:"user-base"`
}

type MessageBase struct {
	CpuNum     int    `mapstructure:"cpuNum"`
	RpcAddress string `mapstructure:"rpcAddress"`
	SdName     string `mapstructure:"sdName"`
	SdVersion  string `mapstructure:"sdVersion"`
	SdDir      string `mapstructure:"sdDir"`
}

type MessageConfig struct {
	MessageBase MessageBase `mapstructure:"message-base"`
}

type ConnectConfig struct {
	ConnectBase                 ConnectBase                 `mapstructure:"connect-base"`
	ConnectWebsocket            ConnectWebsocket            `mapstructure:"connect-websocket"`
	ConnectRpcAddressWebsockets ConnectRpcAddressWebsockets `mapstructure:"connect-rpcAddress-websockets"`
	ConnectBucket               ConnectBucket               `mapstructure:"connect-bucket"`
}

type ConnectBase struct {
	CertPath string `mapstructure:"certPath"`
	KeyPath  string `mapstructure:"keyPath"`
}

type ConnectWebsocket struct {
	Bind          string `mapstructure:"bind"`
	ReadBufSize   int    `mapstructure:"readBufSize"`
	WriteBufSize  int    `mapstructure:"writeBufSize"`
	MaxMsgSize    int    `mapstructure:"maxMsgSize"`
	BroadcastSize int    `mapstructure:"broadcastSize"`
	WriteWait     int    `mapstructure:"writeWait"`
	PongWait      int    `mapstructure:"pongWait"`
	PingPeriod    int    `mapstructure:"pingPeriod"`
}

type ConnectRpcAddressWebsockets struct {
	Address   string `mapstructure:"address"`
	SdName    string `mapstructure:"sdName"`
	SdVersion string `mapstructure:"sdVersion"`
	SdDir     string `mapstructure:"sdDir"`
}

type ConnectBucket struct {
	CpuNum      int `mapstructure:"cpuNum"`
	ChannelSize int `mapstructure:"channelSize"`
}

type SenderBase struct {
	CpuNum   int `mapstructure:"cpuNum"`
	PushChan int `mapstructure:"pushChan"`
}

type SenderConfig struct {
	SenderBase SenderBase `mapstructure:"sender-base"`
}
