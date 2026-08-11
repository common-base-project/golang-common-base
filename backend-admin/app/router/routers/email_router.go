package routers

/*
   @Author : Mustang Kong
*/

import (
	"fmt"
	email2 "golang-common-base/app/handler/email"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// email 路由
func EmailRouter(g *gin.Engine) {
	routerEmail := fmt.Sprintf("%s%s", viper.GetString(`api.version`), "/email")
	email := g.Group(routerEmail)
	{
		// Upload
		email.GET("/list", email2.GetEmailListHandler)
		email.POST("/add", email2.AddEmailHandler)
		email.PUT("/update/:id", email2.UpdateEmailHandler)
		email.DELETE("/delete/:contentId", email2.DeleteEmailHandler)
		email.POST("/push", email2.AddPushHandler)
	}
}
