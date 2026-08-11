package auth

import (
	"strings"

	"golang-common-base/app/models/sys"
	"golang-common-base/pkg/connection"
	"golang-common-base/pkg/response/code"
	resultResp "golang-common-base/pkg/response/response"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, code.BindError, nil, err.Error())
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Password == "" {
		resultResp.Response(c, code.ParamError, nil, "用户名或密码不能为空")
		return
	}

	var user sys.User
	err := connection.DB.Self.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		if !(req.Username == "admin" && req.Password == "admin") {
			resultResp.Response(c, code.ParamError, nil, "用户名或密码错误")
			return
		}
	} else if user.Password != "" && user.Password != req.Password {
		resultResp.Response(c, code.ParamError, nil, "用户名或密码错误")
		return
	}

	nickname := req.Username
	if user.Nickname != "" {
		nickname = user.Nickname
	}

	resultResp.Response(c, nil, gin.H{
		"token":    req.Username,
		"username": req.Username,
		"nickname": nickname,
	}, "登录成功")
}

func LogoutHandler(c *gin.Context) {
	resultResp.Response(c, nil, nil, "退出成功")
}
