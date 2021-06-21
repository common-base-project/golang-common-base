package request

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
)

var orgcharts_url = "http://127.0.0.1:8888/api/v1/employees"

type RespOrgcharts struct {
	Code     int         `json:"errno"`
	ErrorMsg string      `json:"errmsg"`
	Data     interface{} `json:"data"`
}

// 获取用户信息根据用户 ldap
func GetUserInfoByLdap(ldap string) (data interface{}, err error) {
	client := resty.New()
	resp, err := client.R().SetQueryParam("ldap", ldap).Get(orgcharts_url)
	if err != nil {
		return nil, err
	}

	var cResp RespOrgcharts
	err = json.Unmarshal(resp.Body(), &cResp)
	if err != nil {
		log.Printf("JSON解析失败-%s", err.Error())
		return nil, err
	}
	if cResp.Code != 0 {
		log.Printf("获取用户信息失败-%d-%s", cResp.Code, cResp.ErrorMsg)
		return nil, err
	}
	return cResp.Data, nil
}

// 获取用户和部门信息
func GetUserAndDepInfoByLdap(ldap string) (u map[string]interface{}, d map[string]interface{}, e error) {
	data, err := GetUserInfoByLdap(ldap)
	if err != nil {
		log.Println("err")
		return nil, nil, err
	}

	if data == nil {
		return nil, nil, err
	}
	//fmt.Print(data)
	m := data.(map[string]interface{})
	user_rp := m["employees_list"].([]interface{})
	if len(user_rp) <= 0 {
		return nil, nil, err
	}

	user := user_rp[0].(map[string]interface{})
	//fmt.Print(user)
	dep1_rp := user["dept_info"].([]interface{})
	if dep1_rp == nil || len(dep1_rp) <= 0 {
		return user, nil, nil
	}
	dep := dep1_rp[0].(map[string]interface{})
	//fmt.Print(dep)
	return user, dep, nil
}

// 获取部门数据，根据部门ID
