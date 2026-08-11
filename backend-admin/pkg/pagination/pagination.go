package pagination

/*
 @Author : Mustang Kong
*/

import (
	"fmt"
	"golang-common-base/pkg/logger"
	"math"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Param struct {
	C       *gin.Context
	DB      *gorm.DB
	ShowSQL bool
}

type Paginator struct {
	TotalCount int         `json:"total_count"`
	TotalPage  int         `json:"total_page"`
	Data       interface{} `json:"data"`
	PageSize   int         `json:"page_size"`
	Page       int         `json:"page"`
}

type ListRequest struct {
	Page     int `json:"page" form:"page"`           // 页码
	PageSize int `json:"page_size" form:"page_size"` // 每页数据数量
	Sort     int `json:"sort" form:"sort"`
}

// Paging 分页
func Paging(p *Param, data interface{}, result interface{}, args ...interface{}) (*Paginator, error) {
	var (
		param     ListRequest
		paginator Paginator
		count     int64
		tableName string
	)

	if err := p.C.Bind(&param); err != nil {
		logger.Errorf("参数绑定失败，错误：%v", err)
		return nil, err
	}

	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}

	if param.Page < 1 {
		param.Page = 1
	}

	if param.PageSize == 0 {
		param.PageSize = 10
	}

	if param.Sort == 0 || param.Sort == -1 {
		db = db.Order("id desc")
	}

	if len(args) > 1 {
		tableName = fmt.Sprintf("`%s`.", args[1].(string))
	}

	if len(args) > 0 {
		for paramType, paramsValue := range args[0].(map[string]map[string]interface{}) {
			if paramType == "like" {
				for key, value := range paramsValue {
					db = db.Where(fmt.Sprintf("%v%v like ?", tableName, key), fmt.Sprintf("%%%v%%", value))
				}
			} else if paramType == "equal" {
				for key, value := range paramsValue {
					db = db.Where(fmt.Sprintf("%v%v = ?", tableName, key), value)
				}
			}
		}
	}

	err := db.Model(data).
		Scopes(
			Paginate(param.PageSize, param.Page),
		).
		Find(result).Limit(-1).Offset(-1).
		Count(&count).Error
	if err != nil {
		logger.Errorf("db error: %s", err)
		return nil, err
	}

	//done := make(chan bool, 1)
	//go countRecords(db, result, done, &count)
	//offset = (param.Page - 1) * param.PageSize
	//
	//err := db.Limit(param.PageSize).Offset(offset).Find(result).Error
	//if err != nil {
	//	logger.Errorf("数据查询失败，错误：%v", err)
	//	return nil, err
	//}
	//<-done

	paginator.TotalCount = int(count)
	paginator.Data = result
	paginator.Page = param.Page
	paginator.PageSize = param.PageSize
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(param.PageSize)))

	return &paginator, nil
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int64) {
	db.Model(anyType).Count(count)
	done <- true
}

func Paginate(pageSize, pageIndex int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageIndex - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(pageSize)
	}
}
