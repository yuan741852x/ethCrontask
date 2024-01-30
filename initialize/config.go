package initialize

import (
	"ethCrontask/global"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig() {
	v := viper.New()
	// //文件的路径如何设置
	v.SetConfigFile("config.yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	//这个对象如何进行全局使用
	err = v.Unmarshal(global.ServerConfig)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("总后台配置信息:", global.ServerConfig)
	//viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
	})
}
