package auth

import (
	"golang-common-base/app/models/base"
	"golang-common-base/pkg/connection"
	"golang-common-base/pkg/response/code"
	"log"
)

/*
  @Author : Mustang
*/

type User struct {
	base.Model
	Username       string `gorm:"column:username; not null; unique; type:varchar(45);" json:"username" form:"username"`
	Nickname       string `gorm:"column:nickname; type:varchar(45);" json:"nickname" form:"nickname"`
	Email          string `gorm:"column:email; type:varchar(128);" json:"email" form:"email"`
	Hrbp           string `gorm:"column:hrbp; type:varchar(128);" json:"hrbp" form:"hrbp"`
	VP             string `gorm:"column:vp; type:varchar(128);" json:"vp" form:"vp"`
	SubLeader      string `gorm:"column:sub_leader; type:varchar(128);" json:"sub_leader" form:"sub_leader"`
	Position       string `gorm:"column:position; type:varchar(128);" json:"position" form:"position"`
	Depart         int    `gorm:"column:depart_id; type:int(11);" json:"depart_id" form:"depart_id"`
	Status         string `gorm:"column:status; type:varchar(256);" json:"status" form:"status"`
	ReportTo       string `gorm:"column:report_to; type:varchar(256);" json:"report_to" form:"report_to"`
	OnBoardingTime string `gorm:"column:onboardingtime; type:varchar(256);" json:"onboardingtime" form:"onboardingtime"`
}

func (User) TableName() string {
	return "auth_user"
}

func (g *User) CountUser() (usercount int64, err error) {
	usercount = 0
	if err = connection.DB.Self.Model(&User{}).Where("username = ?", g.Username).Count(&usercount).Error; err != nil {
		log.Printf("统计用户数异常: %s", err)
	}
	return usercount, err
}

func (g *User) CreateUser() (err error) {
	if err = connection.DB.Self.Save(g).Error; err != nil {
		log.Println(code.CreateUserError)
	}
	return err
}

func (g *User) UpdateUser() (err error) {
	if err = connection.DB.Self.Model(&g).Where("username = ?", g.Username).Updates(g).Error; err != nil {
		log.Println(code.UpdateUserError)
	}
	return err
}

func (g *User) GetUser() (err error) {
	if err = connection.DB.Self.Set("gorm:auto_preload", true).Where(g).First(g).Error; err != nil {
		log.Println("Couldn's get User.")
	}
	return err
}
