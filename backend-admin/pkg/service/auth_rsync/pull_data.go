package auth_rsync

// import (
// 	"crypto/tls"
// 	"encoding/json"
// 	"log"
// 	"strings"
// 	"time"
// 	"golang-common-base/models/auth"
// 	"golang-common-base/pkg/connection"
// 	_ "golang-common-base/pkg/response/code"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/go-resty/resty/v2"
// )

// type MyResp struct {
// 	Code int         `json:"code"`
// 	Msg  string      `json:"message"`
// 	Data interface{} `json:"data"`
// }

// type AccessToken struct {
// 	RequestTime int64    `json:"request_time"`
// 	User        string   `json:"user"`
// 	Groups      []string `json:"groups"`
// 	ServiceName string   `json:"service_name"`
// }

// var secret interface{} = []byte("D&023u@981jwoIie_!@#*s;lij!poW2ireJLAn3)-")

// const TokenExpire = 5 * time.Second
// const TokenNameInHeader = "Access-Token"

// // 生成token返回tokenString用于设置http header
// func (at *AccessToken) GenerateToken() string {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
// 		"user":         at.User,
// 		"request_time": at.RequestTime,
// 		"groups":       at.Groups,
// 		"service_name": at.ServiceName,
// 	})
// 	if tokenString, err := token.SignedString(secret); err != nil {
// 		//common.Log.Errorf("Token signaure failed: %s", err.Error())
// 		return ""
// 	} else {
// 		return tokenString
// 	}
// }

// // GetUsers
// func GetUserApiHandler() (data interface{}, err error) {
// 	//URL := "http://127.0.0.1:8888/api/v1/employees?page=1&limit=99999"
// 	URL := "https://eim.xxx.cn/api/eim-hr/v1/lightemployees"

// 	client := resty.New()
// 	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
// 	token := AccessToken{ServiceName: "case", RequestTime: time.Now().Unix()}
// 	TokenValue := token.GenerateToken()
// 	req := client.R().
// 		SetQueryParams(map[string]string{
// 			"page":  "1",
// 			"limit": "1000000",
// 		}).
// 		SetHeader("Accept", "application/json").
// 		SetHeader(TokenNameInHeader, TokenValue)
// 	resp, err := req.Get(URL)

// 	if err != nil {
// 		log.Println("拉取用户数据失败!")
// 	}

// 	var cResp MyResp
// 	err = json.Unmarshal(resp.Body(), &cResp)
// 	if err != nil {
// 		log.Printf("JSON解析失败-%s", err.Error())
// 		return
// 	}
// 	if cResp.Code != 100000 {
// 		log.Printf("获取用户列表失败-%d-%s", cResp.Code, cResp.Msg)
// 		return
// 	}

// 	return cResp.Data, nil
// }

// // GetDeparts
// func GetDepartApiHandler() (data interface{}, err error) {
// 	//URL := "http://localhost:30001/api/v1/departments"
// 	//URL := "http://medusa:10000/api/v1/departments"
// 	URL := "https://eim.xxx.cn/api/eim-hr/v1/departments"

// 	client := resty.New()
// 	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
// 	//var cstSh, _ = time.LoadLocation("Asia/Shanghai")
// 	token := AccessToken{ServiceName: "case", RequestTime: time.Now().Unix()}
// 	TokenValue := token.GenerateToken()
// 	req := client.R().
// 		SetQueryParams(map[string]string{
// 			"page":  "1",
// 			"limit": "1000000",
// 		}).
// 		SetHeader("Accept", "application/json").
// 		SetHeader(TokenNameInHeader, TokenValue)
// 	resp, err := req.Get(URL)

// 	if err != nil {
// 		log.Println("拉取部门信息失败!")
// 	}

// 	var cResp MyResp
// 	err = json.Unmarshal(resp.Body(), &cResp)
// 	if err != nil {
// 		log.Printf("JSON解析失败-%s", err.Error())
// 		return
// 	}

// 	if cResp.Code != 100000 {
// 		log.Printf("获取部门列表失败-%d-%s", cResp.Code, cResp.Msg)
// 		return
// 	}

// 	return cResp.Data, nil
// }

// // 获取所有用户数据并入库
// func GetAllUsers() {
// 	users, err := GetUserApiHandler()
// 	//log.Println(users)
// 	if err != nil {
// 		log.Printf("Get user error: %s", err.Error())
// 		return
// 	}

// 	if users != nil {
// 		for _, value := range users.(map[string]interface{})["employees"].([]interface{}) {
// 			var (
// 				User auth.User
// 			)
// 			User.Username = value.(map[string]interface{})["username"].(string)
// 			User.Nickname = value.(map[string]interface{})["name"].(string)
// 			User.Position = value.(map[string]interface{})["position"].(string)
// 			User.Email = User.Username + "@xxx.com"
// 			User.ReportTo = value.(map[string]interface{})["report_to"].(string)
// 			User.OnBoardingTime = value.(map[string]interface{})["on_boarding_time"].(string)
// 			User.Hrbp = value.(map[string]interface{})["hrbp"].(string)
// 			User.SubLeader = value.(map[string]interface{})["sub_leader"].(string)
// 			User.VP = value.(map[string]interface{})["vp"].(string)

// 			var Departname string
// 			dpts := value.(map[string]interface{})["department_rec_name"].(string)
// 			if dpts != "" {
// 				dptList := strings.Split(dpts, "-")
// 				if len(dptList) > 0 {
// 					Departname = dptList[len(dptList)-1]
// 				}
// 			} else {
// 				Departname = ""
// 			}
// 			User.Status = "1"

// 			//DepartID
// 			var departValue auth.Depart
// 			_ = connection.DB.Self.Model(&auth.Depart{}).Where("cname = ?", Departname).Find(&departValue).Error

// 			User.Depart = int(departValue.Id)

// 			userCount, _ := User.CountUser()
// 			if userCount == 0 {
// 				if err := User.CreateUser(); err != nil {
// 					log.Printf("同步用户数据-创建用户失败: %s", err.Error())
// 				}
// 			} else if userCount == 1 {
// 				if err := User.UpdateUser(); err != nil {
// 					log.Printf("同步用户数据-更新用户失败: %s", err.Error())
// 				}
// 			}
// 		}
// 	}
// }

// // 获取所有部门并入库
// func GetAllDeparts() (err error) {

// 	departs, err := GetDepartApiHandler()
// 	if err != nil {
// 		log.Printf("Get user error: %s", err.Error())
// 		return
// 	}
// 	if departs != nil {
// 		for _, value := range departs.([]interface{}) {
// 			var (
// 				Depart auth.Depart
// 			)

// 			Depart.Name = value.(map[string]interface{})["name"].(string)
// 			Depart.Cname = value.(map[string]interface{})["c_name"].(string)
// 			Depart.Leader = value.(map[string]interface{})["leader"].(string)

// 			departcount, _ := Depart.CountDepart()
// 			if departcount == 0 {
// 				if err := Depart.CreateDepart(); err != nil {
// 					log.Printf("同步部门数据-创建部门失败: %s", err.Error())
// 				}
// 			} else if departcount == 1 {
// 				err = Depart.GetDepart()
// 				if err != nil {
// 					if err := Depart.UpdateDepart(); err != nil {
// 						log.Printf("同步部门数据-更新部门失败: %s", err.Error())
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return
// }
