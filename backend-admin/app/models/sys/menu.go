package sys

import "golang-common-base/app/models/base"

type Menu struct {
	base.Model
	ParentID  uint64 `gorm:"column:parent_id;default:0;" json:"parent_id"`
	Name      string `gorm:"column:name;not null;type:varchar(64);" json:"name"`
	Path      string `gorm:"column:path;type:varchar(255);" json:"path"`
	Component string `gorm:"column:component;type:varchar(255);" json:"component"`
	Perms     string `gorm:"column:perms;type:varchar(255);" json:"perms"`
	Icon      string `gorm:"column:icon;type:varchar(64);" json:"icon"`
	Type      int    `gorm:"column:type;default:0;" json:"type"`
	Sort      int    `gorm:"column:sort;default:0;" json:"sort"`
	Visible   int    `gorm:"column:visible;default:1;" json:"visible"`
}

func (Menu) TableName() string {
	return "sys_menu"
}
