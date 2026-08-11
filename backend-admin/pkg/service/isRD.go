package service

// import (
// 	"fmt"
// 	"strings"
// )

// /*
//   @Author : Mustang Kong
// */

// func IsRD(username string) (status bool, err error) {

// 	resp, err := RequestGateway("https://eim.xxx.cn/api/eim-hr/v1/employees/"+username, map[string]string{
// 		"permed_fields": "departments",
// 		"at_time":       "",
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Printf("工单创建者：%v，工单创建者类型：%v\n", username, strings.ToUpper(resp.Data.(map[string]interface{})["employee"].(map[string]interface{})["position_type"].(string)))

// 	if strings.ToUpper(resp.Data.(map[string]interface{})["employee"].(map[string]interface{})["position_type"].(string)) == "R" ||
// 		strings.ToUpper(resp.Data.(map[string]interface{})["employee"].(map[string]interface{})["position_type"].(string)) == "D" {
// 		status = true
// 	}
// 	return
// }
