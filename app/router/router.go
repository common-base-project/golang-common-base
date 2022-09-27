package router

/*
  @Author : Mustang Kong
*/

import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files"       // swagger embed files

import (
	"fmt"
	"golang-common-base/app/middleware"
	"golang-common-base/app/router/routers"
	_ "golang-common-base/docs"
	result "golang-common-base/pkg/response/response"
	"golang-common-base/pkg/utils"
	"net/http"
	"time"

	"github.com/spf13/viper"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// Load 加载路由
func Load(g *gin.Engine) {
	// 404
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, result.ResponseData{
			Errno:  404,
			Errmsg: "API地址不存在",
			Data:   nil,
		})
	})

	// pprof router
	pprof.Register(g)

	//cors， 跨域
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowOrigins:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	// 注册zap相关中间件
	g.Use(cors.New(config))
	//g.Use(logger.GinLogger(), logger.GinRecovery(true))
	g.Use(utils.CostTime())

	// ========================文件配置===============================
	//filePath := viper.GetString("filePath")
	//_, err := tools.CreateDictByPath(filePath)
	//if err != nil {
	//	logger.Error("创建目录失败，请手动创建![%v]\n", err)
	//	return
	//}
	//logger.Infof("创建目录成功: %s", filePath)
	//
	//staticPath := fmt.Sprintf("%s%s", viper.GetString(`api.version`), "/upload")
	//// 静态文件地址 http://localhost:port/api/v1/upload/fileid.jpg
	//g.Static(staticPath, filePath)
	//
	//// g.POST("/api/v1/upload", upload.UploadMutiFileHandler)
	//// 设置文件大小，文件最大为10M (默认 32 MiB)
	//g.MaxMultipartMemory = 5000 << 20 // 500M
	// =======================================================

	// swagger api docs
	g.GET(fmt.Sprintf("%s%s", viper.GetString(`api.version`), "/swagger/*any"), ginSwagger.WrapHandler(swaggerFiles.Handler))

	// jwt 检查
	g.Use(middleware.CheckToken())

	// user
	routers.UserRouter(g)

	// email
	routers.EmailRouter(g)
}
