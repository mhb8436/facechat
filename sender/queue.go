package sender

import (
	"facechat/config"
	"facechat/tools"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

func (s *Sender) InitQueueRedisClient() (err error) {
	redisOpt := tools.RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}
	RedisClient = tools.GetRedisInstance(redisOpt)
	if pong, err := RedisClient.Ping().Result(); err != nil {
		logrus.Infof("RedisClient Ping Result pong: %s,  err: %s", pong, err)
	}
	go func() {
		for {
			var result []string
			//10s timeout
			result, err = RedisClient.BRPop(time.Second*10, config.QueueName).Result()
			if err != nil {
				logrus.Infof("task queue block timeout,no msg err:%s", err.Error())
			}
			if len(result) >= 2 {
				s.Push(result[1])
			}
		}
	}()
	return
}
