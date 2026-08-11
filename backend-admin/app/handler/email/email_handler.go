package email

import (
	"golang-common-base/app/models/email"
	"golang-common-base/pkg/connection"
	"golang-common-base/pkg/logger"
	email2 "golang-common-base/pkg/notify/email"
	"golang-common-base/pkg/pagination"
	"golang-common-base/pkg/response/code"
	. "golang-common-base/pkg/response/response"
	"golang-common-base/pkg/service"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
邮件处理
*/

// GetEmailListHandler @Summary 获取 email 列表
// @Description 获取 email 列表
// @Tags email
// @Accept json
// @Produce json
// @Param page query int false "获取第几页的数据，默认为：1"
// @Param page_size query int false "每页展示多少行，默认为：10"
// @Param sort query int false "按照倒序或者顺序的方式排列，0或者-1为倒序，其他值为顺序"
// @Success 0 {string} json "{"code":0,"message":"获取分类列表成功","data":{}}"
// @Router /api/v1/email/list [get]
func GetEmailListHandler(c *gin.Context) {
	SearchParams := map[string]map[string]interface{}{
		"like": pagination.RequestParams(c),
	}

	var data email.EmailTextContent
	var emailList []*email.EmailTextContent
	result, err := pagination.Paging(&pagination.Param{
		C:  c,
		DB: connection.DB.Self,
	}, &data, &emailList, SearchParams, "t_plainText")

	if err != nil {
		logger.Error(err.Error())
		Response(c, code.SelectCommonError, nil, err.Error())
		return
	}
	Response(c, nil, result, "")
}

// AddEmailHandler 添加 email 数据到数据库
// @Summary 添加 email 数据到数据库
// @Description 添加 email 数据到数据库
// @Tags email
// @Accept json
// @Produce json
// @Param param body email.EmailTextContent true "新建email数据"
// @Success 0
// @Router /api/v1/email/add [post]
func AddEmailHandler(c *gin.Context) {
	var text email.EmailTextContent
	err := c.ShouldBind(&text)
	if err != nil {
		Response(c, code.CreateCommonError, nil, err.Error())
		return
	}
	text.From = viper.GetString("email.support.sender")

	err = connection.DB.Self.Create(&text).Error
	if err != nil {
		Response(c, code.CreateCommonError, nil, err.Error())
		return
	}

	Response(c, nil, nil, "")
}

// UpdateEmailHandler 更新 email 数据
// @Summary 更新 email 数据
// @Description 更新分类	{"name":"test1234","key":"mus_test","child":{"0-":"test"}}"
// @Tags email
// @Accept json
// @Produce json
// @Param param body email.EmailTextContent true "更新email数据"
// @Success 0 {string} json "{"code":0,"message":"更新email 数据成功","data":"32"}"
// @Router /api/v1/email/update [put]
func UpdateEmailHandler(c *gin.Context) {
	var textContent email.EmailTextContent
	err := c.ShouldBind(&textContent)
	if err != nil {
		Response(c, code.UpdateCommonError, nil, err.Error())
		return
	}
	err = connection.DB.Self.Model(&email.EmailTextContent{}).
		Where("id = ?", textContent.Id).
		Updates(map[string]interface{}{
			"mail_to": textContent.To,
			"mail_cc": textContent.CC,
			"subject": textContent.Subject,
			"content": textContent.Content,
		}).Error
	if err != nil {
		Response(c, code.UpdateCommonError, nil, err.Error())
		return
	}
	Response(c, nil, nil, "")
}

// DeleteEmailHandler 删除 email 数据
// @Summary 删除 email 数据
// @Description 删除 email 数据
// @Tags email
// @Accept json
// @Produce json
// @Success 0 {string} json "{"code":0,"message":"删除email成功","data":"32"}"
// @Router /api/v1/email/delete/:contentId [delete]
func DeleteEmailHandler(c *gin.Context) {
	contentId := c.Param("contentId")
	err := connection.DB.Self.Where("id = ?", contentId).Delete(&email.EmailTextContent{}).Error
	if err != nil {
		Response(c, code.DeleteCommonError, nil, err.Error())
		return
	}
	Response(c, nil, nil, "")
}

// AddPushHandler 推送 email 数据到数据库
// @Summary 推送 email 数据到数据库
// @Description 推送 email 数据到数据库
// @Tags email
// @Accept json
// @Produce json
// @Param param body email.EmailTextContent true "新建email数据"
// @Success 0
// @Router /api/v1/email/push [post]
func AddPushHandler(c *gin.Context) {
	var text email.EmailTextContent
	err := c.ShouldBind(&text)
	if err != nil {
		Response(c, code.CreateCommonError, nil, err.Error())
		return
	}
	match, _ := regexp.MatchString(`(support|eim)@xxx\.com`, text.From)
	if !match {
		text.From = viper.GetString("email.support.sender")
	}

	err = connection.DB.Self.Create(&text).Error
	if err != nil {
		Response(c, code.CreateCommonError, nil, err.Error())
		return
	}

	go func() {
		var content = email2.PlainTextContent{
			From:    text.From,
			To:      text.To,
			CC:      text.CC,
			Subject: text.Subject,
			Content: text.Content,
		}
		email2.SendEmail(text.From, content)
		logger.Infof("发送email到: %s  %s", content.To, content.Subject)

		pushUrl := viper.GetString("pushUrl")
		resp, err := service.RequestPost(pushUrl, text)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Infof("发送到飞书: %s  %s", content.To, string(resp))
	}()

	Response(c, nil, nil, "")
}
