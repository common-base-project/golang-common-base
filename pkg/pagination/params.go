package pagination

import (
	"golang-common-base/pkg/logger"

	"github.com/gin-gonic/gin"
)

/*
  @Author : Mustang Kong
*/

func RequestParams(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{}, 10)

	if c.Request.Form == nil {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			logger.Error(err)
		}
	}

	if len(c.Request.Form) > 0 {
		for key, value := range c.Request.Form {
			if key == "page" || key == "page_size" || key == "sort" {
				continue
			}
			params[key] = value[0]
		}
	}

	return params
}
