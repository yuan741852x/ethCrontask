package initialize

import (
	"ethCrontask/global"
	"fmt"
	"github.com/go-redis/redis"
)

func InitRedis() {
	redisInfo := global.ServerConfig.RedisInfo
	sdn := fmt.Sprintf("%s:%d", redisInfo.Host, redisInfo.Port)
	global.Rdb0 = redis.NewClient(&redis.Options{
		Addr:     sdn,
		Password: redisInfo.Password,
		DB:       0,
	})
	global.Rdb1 = redis.NewClient(&redis.Options{
		Addr:     sdn,
		Password: redisInfo.Password,
		DB:       1,
	})
	global.Rdb2 = redis.NewClient(&redis.Options{
		Addr:     sdn,
		Password: redisInfo.Password,
		DB:       2,
	})
}
