package routers

import (
	"fmt"
	handler "golang-common-base/app/handler/sys"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SysRouter(g *gin.Engine) {
	sysAPI := g.Group(fmt.Sprintf("%s/sys", viper.GetString(`api.version`)))
	{
		user := sysAPI.Group("/user")
		{
			user.GET("/info", handler.UserInfoHandler)
			user.GET("/page", handler.UserPageHandler)
			user.GET("", handler.UserListHandler)
			user.GET("/export", handler.UserExportCompatHandler)
			user.GET("/:id", handler.UserDetailCompatHandler)
			user.POST("", handler.UserUpsertCompatHandler)
			user.PUT("", handler.UserUpsertCompatHandler)
			user.PUT("/:id", handler.UserUpsertCompatHandler)
			user.PUT("/password", handler.UserPasswordCompatHandler)
			user.DELETE("", handler.UserBatchDeleteCompatHandler)
			user.DELETE("/:id", handler.DeleteUserHandler)
			user.PUT("/:id/roles", handler.SetUserRolesHandler)
			user.GET("/:id/roles", handler.GetUserRolesHandler)
		}

		role := sysAPI.Group("/role")
		{
			role.GET("/page", handler.RolePageCompatHandler)
			role.GET("/list", handler.RoleListCompatHandler)
			role.GET("", handler.RoleListCompatHandler)
			role.GET("/:id", handler.RoleDetailCompatHandler)
			role.POST("", handler.RoleUpsertCompatHandler)
			role.PUT("", handler.RoleUpsertCompatHandler)
			role.PUT("/:id", handler.RoleUpsertCompatHandler)
			role.DELETE("", handler.RoleBatchDeleteCompatHandler)
			role.DELETE("/:id", handler.DeleteRoleHandler)
			role.PUT("/:id/menus", handler.SetRoleMenusHandler)
			role.GET("/:id/menus", handler.GetRoleMenusHandler)
		}

		menu := sysAPI.Group("/menu")
		{
			menu.GET("/nav", handler.MenuNavHandler)
			menu.GET("/permissions", handler.PermissionListHandler)
			menu.GET("/list", handler.MenuListCompatHandler)
			menu.GET("/select", handler.MenuSelectCompatHandler)
			menu.GET("", handler.MenuListCompatHandler)
			menu.GET("/:id", handler.MenuDetailCompatHandler)
			menu.POST("", handler.MenuUpsertCompatHandler)
			menu.PUT("", handler.MenuUpsertCompatHandler)
			menu.PUT("/:id", handler.MenuUpsertCompatHandler)
			menu.DELETE("", handler.MenuBatchDeleteCompatHandler)
			menu.DELETE("/:id", handler.DeleteMenuHandler)
		}

		dict := sysAPI.Group("/dict")
		{
			dict.GET("/type/all", handler.DictTypeAllHandler)
			dict.GET("/type/page", handler.DictTypePageCompatHandler)
			dict.GET("/type/:id", handler.DictTypeDetailCompatHandler)
			dict.POST("/type", handler.DictTypeUpsertCompatHandler)
			dict.PUT("/type", handler.DictTypeUpsertCompatHandler)
			dict.DELETE("/type", handler.DictTypeDeleteCompatHandler)

			dict.GET("/data/page", handler.DictDataPageCompatHandler)
			dict.GET("/data/:id", handler.DictDataDetailCompatHandler)
			dict.POST("/data", handler.DictDataUpsertCompatHandler)
			dict.PUT("/data", handler.DictDataUpsertCompatHandler)
			dict.DELETE("/data", handler.DictDataDeleteCompatHandler)
		}

		dept := sysAPI.Group("/dept")
		{
			dept.GET("/list", handler.DeptListCompatHandler)
			dept.GET("", handler.DeptListCompatHandler)
			dept.GET("/:id", handler.DeptDetailCompatHandler)
			dept.POST("", handler.DeptUpsertCompatHandler)
			dept.PUT("", handler.DeptUpsertCompatHandler)
			dept.PUT("/:id", handler.DeptUpsertCompatHandler)
			dept.DELETE("", handler.DeptBatchDeleteCompatHandler)
			dept.DELETE("/:id", handler.DeleteDeptHandler)
		}

		region := sysAPI.Group("/region")
		{
			region.GET("/tree", handler.RegionTreeCompatHandler)
		}

		params := sysAPI.Group("/params")
		{
			params.GET("/page", handler.ParamsPageCompatHandler)
			params.GET("/:id", handler.ParamsDetailCompatHandler)
			params.POST("", handler.ParamsUpsertCompatHandler)
			params.PUT("", handler.ParamsUpsertCompatHandler)
			params.DELETE("", handler.ParamsDeleteCompatHandler)
		}

		schedule := sysAPI.Group("/schedule")
		{
			schedule.GET("/page", handler.SchedulePageCompatHandler)
			schedule.GET("/:id", handler.ScheduleDetailCompatHandler)
			schedule.POST("", handler.ScheduleUpsertCompatHandler)
			schedule.PUT("", handler.ScheduleUpsertCompatHandler)
			schedule.DELETE("", handler.ScheduleDeleteCompatHandler)
			schedule.PUT("/pause", handler.SchedulePauseCompatHandler)
			schedule.PUT("/resume", handler.ScheduleResumeCompatHandler)
			schedule.PUT("/run", handler.ScheduleRunCompatHandler)
		}

		scheduleLog := sysAPI.Group("/scheduleLog")
		{
			scheduleLog.GET("/page", handler.ScheduleLogPageCompatHandler)
			scheduleLog.GET("/:id", handler.ScheduleLogDetailCompatHandler)
		}

		log := sysAPI.Group("/log")
		{
			log.GET("/login/page", handler.LogLoginPageCompatHandler)
			log.GET("/login/export", handler.LogLoginExportCompatHandler)
			log.GET("/operation/page", handler.LogOperationPageCompatHandler)
			log.GET("/operation/export", handler.LogOperationExportCompatHandler)
			log.GET("/error/page", handler.LogErrorPageCompatHandler)
			log.GET("/error/export", handler.LogErrorExportCompatHandler)
		}

		oss := sysAPI.Group("/oss")
		{
			oss.GET("/page", handler.OssPageCompatHandler)
			oss.GET("/info", handler.OssInfoCompatHandler)
			oss.POST("", handler.OssConfigCompatHandler)
			oss.DELETE("", handler.OssDeleteCompatHandler)
			oss.POST("/upload", handler.OssUploadCompatHandler)
		}
	}
}
