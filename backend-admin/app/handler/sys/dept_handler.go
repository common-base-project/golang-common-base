package sys

import (
	"golang-common-base/app/models/sys"
	"golang-common-base/pkg/connection"
	"golang-common-base/pkg/pagination"
	"golang-common-base/pkg/response/code"
	resultResp "golang-common-base/pkg/response/response"

	"github.com/gin-gonic/gin"
)

func CreateDeptHandler(c *gin.Context) {
	var req sys.Dept
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, code.BindError, nil, err.Error())
		return
	}

	if err := connection.DB.Self.Create(&req).Error; err != nil {
		resultResp.Response(c, code.CreateCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, req, "成功创建部门")
}

func UpdateDeptHandler(c *gin.Context) {
	id := c.Param("id")
	var req sys.Dept
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, code.BindError, nil, err.Error())
		return
	}

	if err := connection.DB.Self.Model(&sys.Dept{}).Where("id = ?", id).Updates(&req).Error; err != nil {
		resultResp.Response(c, code.UpdateCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, req, "成功更新部门")
}

func DeleteDeptHandler(c *gin.Context) {
	id := c.Param("id")
	if err := connection.DB.Self.Delete(&sys.Dept{}, "id = ?", id).Error; err != nil {
		resultResp.Response(c, code.DeleteCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, nil, "成功删除部门")
}

func DeptListHandler(c *gin.Context) {
	var model sys.Dept
	var items []*sys.Dept
	result, err := pagination.Paging(&pagination.Param{C: c, DB: connection.DB.Self}, model, &items)
	if err != nil {
		resultResp.Response(c, code.SelectCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, result, "成功获取部门列表")
}

func DeptDetailHandler(c *gin.Context) {
	id := c.Param("id")
	var item sys.Dept
	if err := connection.DB.Self.Where("id = ?", id).First(&item).Error; err != nil {
		resultResp.Response(c, code.SelectCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, item, "成功获取部门详情")
}
