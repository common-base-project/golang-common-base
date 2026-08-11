package sys

import (
	"golang-common-base/app/models/sys"
	"golang-common-base/pkg/connection"
	"golang-common-base/pkg/pagination"
	"golang-common-base/pkg/response/code"
	resultResp "golang-common-base/pkg/response/response"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(c *gin.Context) {
	var req sys.User
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, code.BindError, nil, err.Error())
		return
	}

	if err := connection.DB.Self.Create(&req).Error; err != nil {
		resultResp.Response(c, code.CreateCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, req, "成功创建系统用户")
}

func UpdateUserHandler(c *gin.Context) {
	id := c.Param("id")
	var req sys.User
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, code.BindError, nil, err.Error())
		return
	}

	if err := connection.DB.Self.Model(&sys.User{}).Where("id = ?", id).Updates(&req).Error; err != nil {
		resultResp.Response(c, code.UpdateCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, req, "成功更新系统用户")
}

func DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	if err := connection.DB.Self.Delete(&sys.User{}, "id = ?", id).Error; err != nil {
		resultResp.Response(c, code.DeleteCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, nil, "成功删除系统用户")
}

func UserListHandler(c *gin.Context) {
	var model sys.User
	var items []*sys.User
	result, err := pagination.Paging(&pagination.Param{C: c, DB: connection.DB.Self}, model, &items)
	if err != nil {
		resultResp.Response(c, code.SelectCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, result, "成功获取系统用户列表")
}

func UserDetailHandler(c *gin.Context) {
	id := c.Param("id")
	var item sys.User
	if err := connection.DB.Self.Where("id = ?", id).First(&item).Error; err != nil {
		resultResp.Response(c, code.SelectCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, item, "成功获取系统用户详情")
}
