/*
@Author : Mustang
*/
package auth

import (
	"golang-common-base/app/models/auth"
	"golang-common-base/pkg/connection"
	"golang-common-base/pkg/curd"
	"golang-common-base/pkg/pagination"
	"golang-common-base/pkg/response/code"
	resultResp "golang-common-base/pkg/response/response"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(c *gin.Context) {
	var User auth.User
	if err := c.ShouldBindJSON(&User); err != nil {
		resultResp.Response(c, code.BindError, nil, err.Error())
	}

	models := auth.User{}

	if err := curd.Create(&curd.Param{
		Name:       "用户",
		Models:     &models,
		Param:      &User,
		WhereValue: User.Username,
	}); err != nil {
		resultResp.Response(c, code.CreateUserError, nil, err.Error())
	}

	resultResp.Response(c, code.Success, User, "成功创建用户!")
}

func UpdateUserHandler(c *gin.Context) {
	var User auth.User
	if err := c.ShouldBindJSON(&auth.User{}); err != nil {
		resultResp.Response(c, code.BindError, nil, err.Error())
	}

	models := auth.User{}
	if err := curd.Update(&curd.Param{
		Name:       "用户",
		Models:     &models,
		Param:      &User,
		WhereValue: User.Username,
	}); err != nil {
		resultResp.Response(c, code.UpdateUserError, nil, err.Error())
	}

	resultResp.Response(c, code.Success, User, "成功更新用户信息!")
}

func DeleteUserHandler(c *gin.Context) {
	var User auth.User
	GID := c.Param("id")

	if err := connection.DB.Self.Delete(&auth.User{}, "id = ?", GID).Error; err != nil {
		resultResp.Response(c, code.DeleteUserError, nil, err.Error())
	}

	resultResp.Response(c, code.Success, User, "成功删除用户!")

}

func UserListHandler(c *gin.Context) {
	var data auth.User
	var userList []*auth.User
	result, err := pagination.Paging(&pagination.Param{
		C:  c,
		DB: connection.DB.Self,
	}, data, &userList)

	if err != nil {
		resultResp.Response(c, code.SelectUserError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, result, "成功获取用户列表")
}

func UserDetailHandler(c *gin.Context) {
	var User auth.User
	userID := c.Param("id")
	err := connection.DB.Self.Where("id = ?", userID).Find(&User).Error
	if err != nil {
		resultResp.Response(c, code.SelectUserError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, User, "获取用户详细信息")

}

// 获取用户部门对应的所有用户
func DeptUserListHandler(c *gin.Context) {
	var (
		err          error
		userInfo     auth.User
		deptUserList []auth.User
	)

	username := c.DefaultQuery("username", "")
	if username == "" {
		username = c.GetString("user")
	}

	err = connection.DB.Self.Model(&auth.User{}).
		Where("username = ?", username).
		Find(&userInfo).Error
	if err != nil {
		resultResp.Response(c, code.SelectUserError, nil, err.Error())
		return
	}

	err = connection.DB.Self.Model(&auth.User{}).
		Where("depart_id = ? and username != ?", userInfo.Depart, username).
		Find(&deptUserList).Error
	if err != nil {
		resultResp.Response(c, code.SelectUserError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, deptUserList, "")
}
