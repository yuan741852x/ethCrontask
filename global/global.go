package global

import (
	"ethCrontask/config"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Rdb0         *redis.Client
	Rdb1         *redis.Client
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Rdb2         *redis.Client
	Db           *gorm.DB
)
