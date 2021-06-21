package tools

import (
	"encoding/json"
	"reflect"
)

/*
  @Author : Mustang Kong
  @Desc :
    数据结构：slice中存的是map
    解决的问题：判断slice中map的key对应的元素是否存在
*/

func Json2sliceMap(list interface{}) ([]map[string]interface{}, error) {
	var listParam []map[string]interface{}
	err := json.Unmarshal([]byte(list.(json.RawMessage)), &listParam)
	if err != nil {
		return nil, err
	}

	return listParam, nil
}

func IsSliceKey(list interface{}, key string) (status bool, err error) {
	switch reflect.TypeOf(list).String() {
	case "json.RawMessage":
		var listParam []map[string]interface{}
		listParam, err = Json2sliceMap(list)
		if err != nil {
			return false, err
		}
		for _, fieldValue := range listParam {
			if fieldValue["field_key"] == key {
				status = true
				return
			}
		}
	}
	return
}
