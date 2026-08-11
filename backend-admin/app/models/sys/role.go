package sys

import "golang-common-base/app/models/base"

type Role struct {
	base.Model
	Name   string `gorm:"column:name;not null;type:varchar(64);" json:"name"`
	Code   string `gorm:"column:code;not null;unique;type:varchar(64);" json:"code"`
	Remark string `gorm:"column:remark;type:varchar(255);" json:"remark"`
	Status int    `gorm:"column:status;default:1;" json:"status"`
}

func (Role) TableName() string {
	return "sys_role"
}
