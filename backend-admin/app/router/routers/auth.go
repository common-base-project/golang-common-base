/*
@Author : Mustang
*/
package routers

import (
	"fmt"
	"golang-common-base/app/handler/auth"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// AuthRouter 登录退出
func AuthRouter(g *gin.Engine) {
	authRouter := g.Group(viper.GetString(`api.version`))
	{
		authRouter.POST("/login", auth.LoginHandler)
		authRouter.POST("/logout", auth.LogoutHandler)
	}
}

// 用户
func UserRouter(g *gin.Engine) {
	authRouterUser := fmt.Sprintf("%s%s", viper.GetString(`api.version`), "/user")
	authUser := g.Group(authRouterUser)
	{
		authUser.POST("", auth.CreateUserHandler)
		authUser.PUT("/update/:id", auth.UpdateUserHandler)
		authUser.DELETE("/:id", auth.DeleteUserHandler)
		authUser.GET("", auth.UserListHandler)
		authUser.GET("/:id", auth.UserDetailHandler)
		authUser.GET("/:id/dept-user", auth.DeptUserListHandler)
	}
}
