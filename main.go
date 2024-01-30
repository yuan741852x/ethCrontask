package main

import (
	"ethCrontask/consumer"
	"ethCrontask/crons"
	"ethCrontask/initialize"
	"time"
)

func main() {
	initialize.InitConfig()
	initialize.InitRedis()
	initialize.InitMysql()
	// 启动 消费者服务
	for {
		go consumer.RunPayInMatch()
		time.Sleep(time.Millisecond * 1)
	}
	// 定时任务
	crons.Init()
}
