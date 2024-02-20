package main

import (
	"ethCrontask/consumer"
	"ethCrontask/crons"
	"ethCrontask/initialize"
)

func main() {
	initialize.InitConfig()
	initialize.InitRedis()
	initialize.InitMysql()
	// 启动 消费者服务

	go consumer.RunPayInMatch()
	// 定时任务
	crons.Init()
	select {}
}
