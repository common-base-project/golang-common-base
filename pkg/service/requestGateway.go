package service

import (
	"crypto/tls"
	"encoding/json"
	"golang-common-base/pkg/utils"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

/*
  @Author : Mustang Kong
*/

// 请求网关数据
func RequestGateway(URL string, params map[string]string) (cResp utils.MyResp, err error) {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	//var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	token := utils.AccessToken{ServiceName: "case", RequestTime: time.Now().Unix(), ExpTime: time.Now().Unix() + 5000}
	TokenValue := token.GenerateToken()
	req := client.R().
		SetQueryParams(params).
		SetHeader("Accept", "application/json").
		SetHeader(utils.TokenNameInHeader, TokenValue)
	resp, err := req.Get(URL)

	if err != nil {
		log.Println("获取数据失败!", err)
		return
	}
	err = json.Unmarshal(resp.Body(), &cResp)
	if err != nil {
		log.Printf("JSON解析失败-%s", err.Error())
		return
	}
	if cResp.Code != 0 {
		log.Printf("获取数据失败-%d-%s", cResp.Code, cResp.Msg)
		return
	}
	return
}

// 请求数据默认
func RequestGet(URL string, params map[string]string) ([]byte, error) {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	//var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	token := utils.AccessToken{ServiceName: "case", RequestTime: time.Now().Unix(), ExpTime: time.Now().Unix() + 5000}
	TokenValue := token.GenerateToken()
	req := client.R().
		SetQueryParams(params).
		SetHeader("Accept", "application/json").
		SetHeader(utils.TokenNameInHeader, TokenValue)
	resp, err := req.Get(URL)
	if err != nil {
		log.Println("获取数据失败!", err)
		return nil, err
	}
	return resp.Body(), nil
}

// 请求数据 post
func RequestPost(URL string, params interface{}) ([]byte, error) {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	//token := utils.AccessToken{ServiceName: "case", RequestTime: time.Now().Unix(), ExpTime: time.Now().Unix() + 5000}
	//TokenValue := token.GenerateToken()
	resp, err := client.R().
		SetBody(params).
		SetHeader("Accept", "application/json").
		//SetHeader(utils.TokenNameInHeader, TokenValue).
		Post(URL)
	//resp, err := req.Post(URL)
	if err != nil {
		log.Println("获取数据失败!", err)
		return nil, err
	}
	return resp.Body(), nil
}
