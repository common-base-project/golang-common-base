package result

/*
  @Author : Mustang Kong
*/

import (
	"fmt"
	"golang-common-base/pkg/logger"
	httpCode "golang-common-base/pkg/response/code"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Errno  int         `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func Response(c *gin.Context, err error, data interface{}, resultText string) {
	code, message := httpCode.DecodeErr(err)

	if err != nil {
		message = fmt.Sprintf("%s，错误：%v", message, resultText)
	}

	if err == nil && resultText != "" {
		message = resultText
	}

	// write log
	if code != httpCode.Success.Errno {
		logger.Error(message)
	}

	// always return http.StatusOK
	c.JSON(http.StatusOK, ResponseData{
		Errno:  code,
		Errmsg: message,
		Data:   data,
	})
}
