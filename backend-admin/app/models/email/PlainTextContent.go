package email

import "golang-common-base/app/models/base"

/**
普通邮件 model
*/
type EmailTextContent struct {
	base.Model
	From    string `gorm:"column:mail_from; not null; type:varchar(128);" json:"mail_from" form:"mail_from"` // 发送者
	To      string `gorm:"column:mail_to; not null; type:varchar(128);" json:"mail_to" form:"mail_to"`       // 接受者
	CC      string `gorm:"column:mail_cc; not null; type:varchar(128);" json:"mail_cc" form:"mail_cc"`       // 抄送
	Subject string `gorm:"column:subject; not null; type:varchar(512);" json:"subject" form:"subject"`       // 主题
	Content string `gorm:"column:content; not null; type:text;" json:"content" form:"content"`               // 消息体
}

func (EmailTextContent) TableName() string {
	return "t_plainText"
}
