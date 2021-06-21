package curd

import (
	"fmt"
	"golang-common-base/pkg/connection"
)

/*
  @Author : Mustang Kong
*/

type Param struct {
	Name       string
	Models     interface{}
	Param      interface{}
	WhereMap   map[string]interface{}
	WhereValue string
}

func whereDB(p *Param) (err error) {
	db := connection.DB.Self

	if p.WhereValue != "" {
		db = db.Where("name = ?", p.WhereValue)
	}

	if p.WhereMap != nil {
		for key, value := range p.WhereMap {
			db = db.Where(fmt.Sprintf("%v = ?", key), value)
		}
	}

	var dataC int64
	dataCount := 0
	err = db.Model(p.Models).Count(&dataC).Error
	if err != nil {
		err = fmt.Errorf("查询%s数据失败，%v", p.Name, err)
		return
	}
	dataCount = int(dataCount)
	if dataCount > 0 {
		err = fmt.Errorf("`%s`数据筛选出现问题，请确认", p.Name)
		return
	}

	return
}

func Create(p *Param) (err error) {

	err = whereDB(p)
	if err != nil {
		return
	}

	err = connection.DB.Self.Save(p.Param).Error
	if err != nil {
		return
	}

	return nil
}
