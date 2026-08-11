package sys

import "golang-common-base/app/models/base"

type User struct {
	base.Model
	Username string `gorm:"column:username;not null;unique;type:varchar(64);" json:"username"`
	Password string `gorm:"column:password;not null;type:varchar(255);" json:"password"`
	Nickname string `gorm:"column:nickname;type:varchar(64);" json:"nickname"`
	Email    string `gorm:"column:email;type:varchar(128);" json:"email"`
	Mobile   string `gorm:"column:mobile;type:varchar(32);" json:"mobile"`
	DeptID   uint64 `gorm:"column:dept_id;default:0;" json:"dept_id"`
	Status   int    `gorm:"column:status;default:1;" json:"status"`
}

func (User) TableName() string {
	return "sys_user"
}
