package curd

import (
	"golang-common-base/pkg/connection"
)

/*
  @Author : Mustang Kong
*/

func Update(p *Param) (err error) {
	err = whereDB(p)

	err = connection.DB.Self.Model(p.Param).Save(p.Param).Error
	if err != nil {
		return
	}

	return nil
}
