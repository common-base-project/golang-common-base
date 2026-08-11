package base

/*
  @Author : Mustang Kong
*/

import (
	"golang-common-base/pkg/jsonTime"
)

type Model struct {
	Id        uint64             `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id" form:"id"`
	CreatedAt jsonTime.JSONTime  `gorm:"column:twf_created" json:"twf_created" form:"twf_created"`
	UpdatedAt jsonTime.JSONTime  `gorm:"column:twf_modified" json:"twf_modified" form:"twf_modified"`
	DeletedAt *jsonTime.JSONTime `gorm:"column:twf_deleted" sql:"index" json:"-" form:"twf_deleted"`
}
