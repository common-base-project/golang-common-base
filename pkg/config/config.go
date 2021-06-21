/*
  @Author : Mustang Kong
*/

package config

import (
	"flag"
	"fmt"
	"golang-common-base/pkg/logger"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// config path
var configDir string
var EnvMode string

type Config struct {
	Name string
}

func Initial() {
	fmt.Println("系统环境变量", EnvMode)
	// 获取配置文件路径
	//configDir = fmt.Sprintf("%s%s", settings.ObjectPath(), "/conf")
	//logger.Info("配置文件路径", configDir)
	var configDir string
	if EnvMode == "prod" {
		configDir = "./conf/config_prod.yaml"
	} else if EnvMode == "staging" {
		configDir = "./conf/config_stage.yaml"
	} else {
		os.Setenv("ENV_SERVER_MODE", "dev")
		configDir = "./conf/config_dev.yaml"
	}
	fmt.Println("配置文件路径", configDir)

	// 全局初始化viper
	cfg := flag.String("config", configDir, "配置文件的路径")
	//cfg := pflag.StringP("config_dev", "c", "", "配置文件的路径")
	//pflag.Parse()
	err := Init(*cfg)
	if err != nil {
		panic(err)
	}

	// 配置文件热加载
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Info(fmt.Sprintf("Config file changed: %s", e.Name))
	})
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}
	if err := c.initConfig(); err != nil {
		return err
	}
	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath(configDir)
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	// 热加载配置信息
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
