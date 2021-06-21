package tools

import (
	"reflect"
	"strings"
)

/*
  @Author : Mustang Kong
*/

func Struct2Map(s interface{}) map[string]interface{} {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		ts := t.Field(i)
		tagName := ts.Tag.Get("json")
		if tagName != "" {
			data[ts.Tag.Get("json")] = v.Field(i).Interface()
		} else {
			data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
		}
	}
	return data
}
