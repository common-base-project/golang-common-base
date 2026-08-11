package sys

import "golang-common-base/app/models/base"

type UserRole struct {
	base.Model
	UserID uint64 `gorm:"column:user_id;not null;index;uniqueIndex:idx_sys_user_role;" json:"user_id"`
	RoleID uint64 `gorm:"column:role_id;not null;index;uniqueIndex:idx_sys_user_role;" json:"role_id"`
}

func (UserRole) TableName() string {
	return "sys_user_role"
}

type RoleMenu struct {
	base.Model
	RoleID uint64 `gorm:"column:role_id;not null;index;uniqueIndex:idx_sys_role_menu;" json:"role_id"`
	MenuID uint64 `gorm:"column:menu_id;not null;index;uniqueIndex:idx_sys_role_menu;" json:"menu_id"`
}

func (RoleMenu) TableName() string {
	return "sys_role_menu"
}
