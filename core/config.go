package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"qianbei.com/config"
)

// 配置文件名称
const (
	defaultConfigName       = "config-local.yaml"
	defaultConfigSimName    = "config-sim.yaml"
	defaultConfigOnlineName = "config-online.yaml"
	defaultConfigDockerName = "config-docker.yaml"
)

var g_config *config.Config

func init() {
	v := viper.New()
	// 根据环境变量 设置配置文件
	switch os.Getenv("mode") {
	case "online":
		v.SetConfigFile(defaultConfigOnlineName)
	case "sim":
		v.SetConfigFile(defaultConfigSimName)
	case "docker":
		v.SetConfigFile(defaultConfigDockerName)
	default:
		v.SetConfigFile(defaultConfigName)
	}
	// 尝试读取 是否报错
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将文件内容写到变量中
	if err := v.Unmarshal(&g_config); err != nil {
		panic(fmt.Errorf("config changge error: %s \n", err))
	}
	// 监听文件
	v.WatchConfig()
	// 文件变化 重新加载内容
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&g_config); err != nil {
			panic(fmt.Errorf("config changge error: %s \n", err))
		}
	})
}

func Config() *config.Config {
	return g_config
}
