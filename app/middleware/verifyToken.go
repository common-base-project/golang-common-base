package middleware

import (
	"regexp"
	"golang-common-base/pkg/utils"

	"github.com/gin-gonic/gin"
)

var Token = utils.AccessToken{}

func CheckToken() func(context *gin.Context) {
	return func(context *gin.Context) {
		if !isMustApi(context) {
			return
		}
		// 获取头部的 token
		tokenString := context.GetHeader(utils.TokenNameInHeader) // token
		requestId := context.GetHeader(utils.RequestID)           // logID
		Token.RequestID = requestId

		if b := Token.ValidateToken(context, tokenString); b {
			context.Next()
			return
		} else {
			context.JSON(200, map[string]interface{}{"errno": 100001, "errmsg": "access token无效"})
			context.Abort()
			return
		}
	}
}

// 定义无需登陆检测的接口
func isMustApi(context *gin.Context) bool {
	return context.Request.URL.Path != "/api/v1/login" &&
		context.Request.URL.Path != "/api/v1/logout" &&
		!ignoreMatchErr(`/api/v1/reservations/([0-9]+)/checkin`, context.Request.URL.Path) &&
		!ignoreMatchErr(`/api/v1/upload`, context.Request.URL.Path)
}

func ignoreMatchErr(pattern, str string) bool {
	match, _ := regexp.MatchString(pattern, str)
	return match
}
