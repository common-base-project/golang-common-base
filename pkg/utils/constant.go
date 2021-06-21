package utils

import (
	"time"
	"golang-common-base/pkg/logger"

	"github.com/gin-gonic/gin"
)

// 统计消耗时间
func CostTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前获取当前时间
		nowTime := time.Now()

		// 请求处理
		c.Next()

		// 处理后获取消耗时间
		costTime := time.Since(nowTime)
		url := c.Request.URL.String()
		// fmt.Printf("the request URL %s cost %v\n", url, costTime)
		logger.Infof("the request URL %s cost %v\n", url, costTime)
	}
}
