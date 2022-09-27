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

// CreateUserHandler 创建用户
// @Summary 创建用户
// @Description 创建用户
// @Tags user
// @Accept json
// @Produce json
// @Param param body auth.User true "用户数据"
// @Success 0
// @Router /api/v1/user [post]
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

// UpdateUserHandler 更新 user 数据
// @Summary 更新 user 数据
// @Description 更新 user 数据
// @Tags user
// @Accept json
// @Produce json
// @Param param body auth.User true "更新user数据"
// @Success 0 {string} json "{"code":0,"message":"更新 user 数据成功","data":"32"}"
// @Router /api/v1/user/update/:id [put]
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

// DeleteUserHandler 删除 user 数据
// @Summary 删除 user 数据
// @Description 删除 user 数据
// @Tags user
// @Accept json
// @Produce json
// @Success 0 {string} json "{"code":0,"message":"删除 user 成功","data":"32"}"
// @Router /api/v1/user/delete/:id [delete]
func DeleteUserHandler(c *gin.Context) {
	var User auth.User
	GID := c.Param("id")

	if err := connection.DB.Self.Delete(&auth.User{}, "id = ?", GID).Error; err != nil {
		resultResp.Response(c, code.DeleteUserError, nil, err.Error())
	}

	resultResp.Response(c, code.Success, User, "成功删除用户!")

}

// UserListHandler 获取 user 列表
// @Summary 获取 user 列表
// @Description 获取 user 列表
// @Tags user
// @Accept json
// @Produce json
// @Param page query int false "获取第几页的数据，默认为：1"
// @Param page_size query int false "每页展示多少行，默认为：10"
// @Param sort query int false "按照倒序或者顺序的方式排列，0或者-1为倒序，其他值为顺序"
// @Success 0 {string} json "{"code":0,"message":"获取分类列表成功","data":{}}"
// @Router /api/v1/user [get]
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

// UserDetailHandler 获取 user 详情
// @Summary 获取 user 详情
// @Description 获取 user 详情
// @Tags user
// @Accept json
// @Produce json
// @Param id query int true "user id"
// @Success 0 {string} json "{"code":0,"message":"获取 user 详情","data":{}}"
// @Router /api/v1/user/:id [get]
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

// DeptUserListHandler 获取用户部门对应的所有用户
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
