/*
  @Author : Mustang Kong
*/

package main

import (
	"fmt"
	"golang-common-base/app/models"
	"golang-common-base/app/router"
	"golang-common-base/pkg/config"
	_ "golang-common-base/pkg/config"
	"golang-common-base/pkg/connection"
	"golang-common-base/pkg/logger"
	"golang-common-base/pkg/service/auth_rsync"
	"os"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

// init all
func init() {
	fmt.Println("#########1", os.Getenv("ENV_SERVER_MODE"))
	config.EnvMode = os.Getenv("ENV_SERVER_MODE")
	config.Initial()
	fmt.Println("#########2", config.EnvMode)
	logger.Initial()
	connection.Initial()
}

// @title golang-common-base API docs
// @version 0.0.1
// @contact.name Mustang Kong
// @contact.email mustang2247@gmail.com
// http://localhost:9080/api/v1/swagger/index.html
func main() {
	// 同步用户和部门数据
	go auth_rsync.Main()

	// 同步数据结构
	models.AutoMigrateTable()

	g := gin.New()

	if config.EnvMode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else if config.EnvMode == "staging" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 加载路由
	router.Load(g)

	// 运行程序
	err := g.Run(viper.GetString(`server.port`))
	if err != nil {
		logger.Error("启动失败")
		panic(fmt.Sprintf("程序启动失败：%v", err))
	}

	defer connection.DB.Close()
}
