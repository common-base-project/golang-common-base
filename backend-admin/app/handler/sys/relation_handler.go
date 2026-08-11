package sys

import (
	"strings"

	"golang-common-base/app/models/sys"
	"golang-common-base/pkg/connection"
	"golang-common-base/pkg/logger"
	"golang-common-base/pkg/response/code"
	resultResp "golang-common-base/pkg/response/response"

	"github.com/gin-gonic/gin"
)

type setUserRolesReq struct {
	RoleIDs []uint64 `json:"role_ids"`
}

type setRoleMenusReq struct {
	MenuIDs []uint64 `json:"menu_ids"`
}

func SetUserRolesHandler(c *gin.Context) {
	userID := c.Param("id")
	var req setUserRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, code.BindError, nil, err.Error())
		return
	}

	tx := connection.DB.Self.Begin()
	if tx.Error != nil {
		resultResp.Response(c, code.InternalServerError, nil, tx.Error.Error())
		return
	}

	if err := tx.Where("user_id = ?", userID).Delete(&sys.UserRole{}).Error; err != nil {
		tx.Rollback()
		resultResp.Response(c, code.UpdateCommonError, nil, err.Error())
		return
	}

	if len(req.RoleIDs) > 0 {
		rows := make([]sys.UserRole, 0, len(req.RoleIDs))
		for _, rid := range req.RoleIDs {
			rows = append(rows, sys.UserRole{UserID: toUint64(userID), RoleID: rid})
		}
		if err := tx.Create(&rows).Error; err != nil {
			tx.Rollback()
			resultResp.Response(c, code.UpdateCommonError, nil, err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		resultResp.Response(c, code.UpdateCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, gin.H{"user_id": userID, "role_ids": req.RoleIDs}, "成功设置用户角色")
}

func GetUserRolesHandler(c *gin.Context) {
	userID := c.Param("id")
	var rows []sys.UserRole
	if err := connection.DB.Self.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		resultResp.Response(c, code.SelectCommonError, nil, err.Error())
		return
	}

	roleIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		roleIDs = append(roleIDs, row.RoleID)
	}

	resultResp.Response(c, nil, gin.H{"user_id": userID, "role_ids": roleIDs}, "成功获取用户角色")
}

func SetRoleMenusHandler(c *gin.Context) {
	roleID := c.Param("id")
	var req setRoleMenusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, code.BindError, nil, err.Error())
		return
	}

	tx := connection.DB.Self.Begin()
	if tx.Error != nil {
		resultResp.Response(c, code.InternalServerError, nil, tx.Error.Error())
		return
	}

	if err := tx.Where("role_id = ?", roleID).Delete(&sys.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		resultResp.Response(c, code.UpdateCommonError, nil, err.Error())
		return
	}

	if len(req.MenuIDs) > 0 {
		rows := make([]sys.RoleMenu, 0, len(req.MenuIDs))
		for _, mid := range req.MenuIDs {
			rows = append(rows, sys.RoleMenu{RoleID: toUint64(roleID), MenuID: mid})
		}
		if err := tx.Create(&rows).Error; err != nil {
			tx.Rollback()
			resultResp.Response(c, code.UpdateCommonError, nil, err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		resultResp.Response(c, code.UpdateCommonError, nil, err.Error())
		return
	}

	if err := syncRoleMenuPolicies(toUint64(roleID), req.MenuIDs); err != nil {
		logger.Errorf("同步角色菜单到 Casbin 失败: %v", err)
		resultResp.Response(c, code.UpdateCommonError, nil, err.Error())
		return
	}

	resultResp.Response(c, nil, gin.H{"role_id": roleID, "menu_ids": req.MenuIDs}, "成功设置角色菜单")
}

func GetRoleMenusHandler(c *gin.Context) {
	roleID := c.Param("id")
	var rows []sys.RoleMenu
	if err := connection.DB.Self.Where("role_id = ?", roleID).Find(&rows).Error; err != nil {
		resultResp.Response(c, code.SelectCommonError, nil, err.Error())
		return
	}

	menuIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		menuIDs = append(menuIDs, row.MenuID)
	}

	resultResp.Response(c, nil, gin.H{"role_id": roleID, "menu_ids": menuIDs}, "成功获取角色菜单")
}

func toUint64(raw string) uint64 {
	var v uint64
	for i := 0; i < len(raw); i++ {
		if raw[i] < '0' || raw[i] > '9' {
			return 0
		}
		v = v*10 + uint64(raw[i]-'0')
	}
	return v
}

func syncRoleMenuPolicies(roleID uint64, menuIDs []uint64) error {
	var role sys.Role
	if err := connection.DB.Self.Where("id = ?", roleID).First(&role).Error; err != nil {
		return err
	}

	if connection.CasbinEnforcer == nil {
		return nil
	}

	_, _ = connection.CasbinEnforcer.RemoveFilteredPolicy(0, role.Code)

	if len(menuIDs) == 0 {
		return connection.CasbinEnforcer.SavePolicy()
	}

	var menus []sys.Menu
	if err := connection.DB.Self.Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
		return err
	}

	for _, menu := range menus {
		obj := strings.TrimSpace(menu.Path)
		if obj == "" {
			obj = strings.TrimSpace(menu.Perms)
		}
		if obj == "" {
			continue
		}
		_, _ = connection.CasbinEnforcer.AddPolicy(role.Code, obj, "(get|post|put|patch|delete|head|options)")
	}

	return connection.CasbinEnforcer.SavePolicy()
}
