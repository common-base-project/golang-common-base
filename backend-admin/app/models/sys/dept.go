package sys

import "golang-common-base/app/models/base"

type Dept struct {
	base.Model
	ParentID uint64 `gorm:"column:parent_id;default:0;" json:"parent_id"`
	Name     string `gorm:"column:name;not null;type:varchar(64);" json:"name"`
	Leader   string `gorm:"column:leader;type:varchar(64);" json:"leader"`
	Phone    string `gorm:"column:phone;type:varchar(32);" json:"phone"`
	Email    string `gorm:"column:email;type:varchar(128);" json:"email"`
	Sort     int    `gorm:"column:sort;default:0;" json:"sort"`
	Status   int    `gorm:"column:status;default:1;" json:"status"`
}

func (Dept) TableName() string {
	return "sys_dept"
}
