package sys

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	sysModel "golang-common-base/app/models/sys"
	"golang-common-base/pkg/connection"
	resultResp "golang-common-base/pkg/response/response"

	"github.com/gin-gonic/gin"
)

type pageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type simpleNode struct {
	ID       uint64       `json:"id"`
	PID      uint64       `json:"pid"`
	Name     string       `json:"name"`
	Children []simpleNode `json:"children,omitempty"`
}

var (
	compatMu sync.RWMutex
	compatID uint64 = 100000

	dictTypes = []map[string]interface{}{
		{"id": uint64(1), "dictName": "性别", "dictType": "gender", "sort": 1, "remark": "", "createDate": nowText()},
	}
	dictData = []map[string]interface{}{
		{"id": uint64(1), "dictTypeId": uint64(1), "dictLabel": "男", "dictValue": "0", "sort": 1, "remark": "", "createDate": nowText()},
		{"id": uint64(2), "dictTypeId": uint64(1), "dictLabel": "女", "dictValue": "1", "sort": 2, "remark": "", "createDate": nowText()},
	}
	paramsData = []map[string]interface{}{
		{"id": uint64(1), "paramCode": "sys.demo", "paramValue": "enabled", "remark": "demo", "createDate": nowText()},
	}
	schedulesData   = []map[string]interface{}{}
	scheduleLogs    = []map[string]interface{}{}
	ossFiles        = []map[string]interface{}{}
	loginLogs       = []map[string]interface{}{}
	operationLogs   = []map[string]interface{}{}
	errorLogs       = []map[string]interface{}{}
	roleDeptBinding = map[uint64][]uint64{}
	ossConfig       = map[string]interface{}{
		"type": 1,
	}
)

func nowText() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func nextCompatID() uint64 {
	return atomic.AddUint64(&compatID, 1)
}

func parseUintFromAny(v interface{}) uint64 {
	switch t := v.(type) {
	case float64:
		return uint64(t)
	case int:
		return uint64(t)
	case int64:
		return uint64(t)
	case uint64:
		return t
	case string:
		n, _ := strconv.ParseUint(strings.TrimSpace(t), 10, 64)
		return n
	default:
		return 0
	}
}

func parseIntFromAny(v interface{}) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case int64:
		return int(t)
	case string:
		n, _ := strconv.Atoi(strings.TrimSpace(t))
		return n
	default:
		return 0
	}
}

func parseStringFromAny(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func parseIDList(v interface{}) []uint64 {
	rs := make([]uint64, 0)
	switch t := v.(type) {
	case []interface{}:
		for _, item := range t {
			id := parseUintFromAny(item)
			if id > 0 {
				rs = append(rs, id)
			}
		}
	case []uint64:
		for _, id := range t {
			if id > 0 {
				rs = append(rs, id)
			}
		}
	case []string:
		for _, item := range t {
			id := parseUintFromAny(item)
			if id > 0 {
				rs = append(rs, id)
			}
		}
	}
	return rs
}

func parseBodyIDs(c *gin.Context) []uint64 {
	var idsNum []uint64
	if err := c.ShouldBindJSON(&idsNum); err == nil {
		return idsNum
	}
	var idsStr []string
	if err := c.ShouldBindJSON(&idsStr); err == nil {
		return parseIDList(idsStr)
	}
	var one map[string]interface{}
	if err := c.ShouldBindJSON(&one); err == nil {
		if id := parseUintFromAny(one["id"]); id > 0 {
			return []uint64{id}
		}
	}
	return nil
}

func writeCSV(c *gin.Context, filename string, header []string, rows [][]string) {
	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	_ = w.Write(header)
	for _, row := range rows {
		_ = w.Write(row)
	}
	w.Flush()
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.String(http.StatusOK, buf.String())
}

func pageLimit(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return page, limit
}

func paginateSlice[T any](items []T, page, limit int) ([]T, int64) {
	total := int64(len(items))
	start := (page - 1) * limit
	if start >= len(items) {
		return []T{}, total
	}
	end := start + limit
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], total
}

func UserPageHandler(c *gin.Context) {
	var users []sysModel.User
	db := connection.DB.Self.Model(&sysModel.User{})
	if keyword := strings.TrimSpace(c.Query("username")); keyword != "" {
		db = db.Where("username like ?", "%"+keyword+"%")
	}
	var total int64
	db.Count(&total)
	page, limit := pageLimit(c)
	db = db.Order("id desc").Offset((page - 1) * limit).Limit(limit)
	if err := db.Find(&users).Error; err != nil {
		resultResp.Response(c, nil, pageResult{List: []interface{}{}, Total: 0}, "获取用户列表失败")
		return
	}

	deptIDs := make([]uint64, 0)
	for _, u := range users {
		if u.DeptID > 0 {
			deptIDs = append(deptIDs, u.DeptID)
		}
	}
	deptMap := map[uint64]string{}
	if len(deptIDs) > 0 {
		var depts []sysModel.Dept
		connection.DB.Self.Where("id in ?", deptIDs).Find(&depts)
		for _, d := range depts {
			deptMap[d.Id] = d.Name
		}
	}

	items := make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		items = append(items, map[string]interface{}{
			"id":         u.Id,
			"username":   u.Username,
			"deptId":     u.DeptID,
			"deptName":   deptMap[u.DeptID],
			"email":      u.Email,
			"mobile":     u.Mobile,
			"gender":     0,
			"status":     u.Status,
			"createDate": fmt.Sprintf("%v", u.CreatedAt),
		})
	}
	resultResp.Response(c, nil, pageResult{List: items, Total: total}, "成功")
}

func UserDetailCompatHandler(c *gin.Context) {
	id := c.Param("id")
	var user sysModel.User
	if err := connection.DB.Self.Where("id = ?", id).First(&user).Error; err != nil {
		resultResp.Response(c, nil, map[string]interface{}{}, "用户不存在")
		return
	}
	var roles []sysModel.UserRole
	connection.DB.Self.Where("user_id = ?", user.Id).Find(&roles)
	roleIDs := make([]uint64, 0, len(roles))
	for _, r := range roles {
		roleIDs = append(roleIDs, r.RoleID)
	}
	deptName := ""
	if user.DeptID > 0 {
		var dept sysModel.Dept
		if err := connection.DB.Self.Where("id = ?", user.DeptID).First(&dept).Error; err == nil {
			deptName = dept.Name
		}
	}
	resultResp.Response(c, nil, map[string]interface{}{
		"id":              user.Id,
		"username":        user.Username,
		"deptId":          user.DeptID,
		"deptName":        deptName,
		"password":        "",
		"confirmPassword": "",
		"realName":        user.Nickname,
		"gender":          0,
		"email":           user.Email,
		"mobile":          user.Mobile,
		"roleIdList":      roleIDs,
		"status":          user.Status,
	}, "成功")
}

func UserUpsertCompatHandler(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, nil, nil, "请求参数错误")
		return
	}
	id := parseUintFromAny(req["id"])
	if id == 0 {
		id = parseUintFromAny(c.Param("id"))
	}
	username := parseStringFromAny(req["username"])
	if username == "" {
		resultResp.Response(c, nil, nil, "用户名不能为空")
		return
	}
	password := parseStringFromAny(req["password"])
	if password == "" && id == 0 {
		password = "123456"
	}

	user := sysModel.User{
		Username: username,
		Password: password,
		Nickname: parseStringFromAny(req["realName"]),
		Email:    parseStringFromAny(req["email"]),
		Mobile:   parseStringFromAny(req["mobile"]),
		DeptID:   parseUintFromAny(req["deptId"]),
		Status:   parseIntFromAny(req["status"]),
	}
	if id == 0 {
		if user.Status == 0 {
			user.Status = 1
		}
		if err := connection.DB.Self.Create(&user).Error; err != nil {
			resultResp.Response(c, nil, nil, "创建用户失败")
			return
		}
		id = user.Id
	} else {
		updates := map[string]interface{}{
			"username": username,
			"nickname": user.Nickname,
			"email":    user.Email,
			"mobile":   user.Mobile,
			"dept_id":  user.DeptID,
			"status":   user.Status,
		}
		if password != "" {
			updates["password"] = password
		}
		if err := connection.DB.Self.Model(&sysModel.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			resultResp.Response(c, nil, nil, "更新用户失败")
			return
		}
	}

	roleIDs := parseIDList(req["roleIdList"])
	connection.DB.Self.Where("user_id = ?", id).Delete(&sysModel.UserRole{})
	for _, rid := range roleIDs {
		if rid == 0 {
			continue
		}
		connection.DB.Self.Create(&sysModel.UserRole{UserID: id, RoleID: rid})
	}
	resultResp.Response(c, nil, gin.H{"id": id}, "成功")
}

func UserBatchDeleteCompatHandler(c *gin.Context) {
	ids := parseBodyIDs(c)
	if len(ids) == 0 {
		resultResp.Response(c, nil, nil, "请选择要删除的用户")
		return
	}
	connection.DB.Self.Where("id in ?", ids).Delete(&sysModel.User{})
	connection.DB.Self.Where("user_id in ?", ids).Delete(&sysModel.UserRole{})
	resultResp.Response(c, nil, nil, "成功")
}

func UserPasswordCompatHandler(c *gin.Context) {
	resultResp.Response(c, nil, nil, "成功")
}

func UserExportCompatHandler(c *gin.Context) {
	var users []sysModel.User
	connection.DB.Self.Order("id desc").Find(&users)
	rows := make([][]string, 0, len(users))
	for _, u := range users {
		rows = append(rows, []string{strconv.FormatUint(u.Id, 10), u.Username, u.Nickname, u.Email, u.Mobile, strconv.Itoa(u.Status)})
	}
	writeCSV(c, "users.csv", []string{"id", "username", "realName", "email", "mobile", "status"}, rows)
}

func RolePageCompatHandler(c *gin.Context) {
	var roles []sysModel.Role
	db := connection.DB.Self.Model(&sysModel.Role{})
	if keyword := strings.TrimSpace(c.Query("name")); keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	var total int64
	db.Count(&total)
	page, limit := pageLimit(c)
	db = db.Order("id desc").Offset((page - 1) * limit).Limit(limit)
	if err := db.Find(&roles).Error; err != nil {
		resultResp.Response(c, nil, pageResult{List: []interface{}{}, Total: 0}, "获取角色列表失败")
		return
	}
	items := make([]map[string]interface{}, 0, len(roles))
	for _, r := range roles {
		items = append(items, map[string]interface{}{
			"id":         r.Id,
			"name":       r.Name,
			"remark":     r.Remark,
			"status":     r.Status,
			"createDate": fmt.Sprintf("%v", r.CreatedAt),
		})
	}
	resultResp.Response(c, nil, pageResult{List: items, Total: total}, "成功")
}

func RoleListCompatHandler(c *gin.Context) {
	var roles []sysModel.Role
	connection.DB.Self.Order("id desc").Find(&roles)
	items := make([]map[string]interface{}, 0, len(roles))
	for _, r := range roles {
		items = append(items, map[string]interface{}{"id": r.Id, "name": r.Name, "remark": r.Remark})
	}
	resultResp.Response(c, nil, items, "成功")
}

func RoleDetailCompatHandler(c *gin.Context) {
	id := parseUintFromAny(c.Param("id"))
	var role sysModel.Role
	if err := connection.DB.Self.Where("id = ?", id).First(&role).Error; err != nil {
		resultResp.Response(c, nil, map[string]interface{}{}, "角色不存在")
		return
	}
	var rms []sysModel.RoleMenu
	connection.DB.Self.Where("role_id = ?", id).Find(&rms)
	menuIDs := make([]uint64, 0, len(rms))
	for _, rm := range rms {
		menuIDs = append(menuIDs, rm.MenuID)
	}
	compatMu.RLock()
	deptIDs := roleDeptBinding[id]
	compatMu.RUnlock()
	if deptIDs == nil {
		deptIDs = []uint64{}
	}
	resultResp.Response(c, nil, map[string]interface{}{
		"id":         role.Id,
		"name":       role.Name,
		"remark":     role.Remark,
		"menuIdList": menuIDs,
		"deptIdList": deptIDs,
	}, "成功")
}

func RoleUpsertCompatHandler(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, nil, nil, "请求参数错误")
		return
	}
	id := parseUintFromAny(req["id"])
	if id == 0 {
		id = parseUintFromAny(c.Param("id"))
	}
	name := parseStringFromAny(req["name"])
	if name == "" {
		resultResp.Response(c, nil, nil, "角色名称不能为空")
		return
	}
	remark := parseStringFromAny(req["remark"])
	if id == 0 {
		role := sysModel.Role{Name: name, Code: fmt.Sprintf("role_%d", time.Now().UnixNano()), Remark: remark, Status: 1}
		if err := connection.DB.Self.Create(&role).Error; err != nil {
			resultResp.Response(c, nil, nil, "创建角色失败")
			return
		}
		id = role.Id
	} else {
		if err := connection.DB.Self.Model(&sysModel.Role{}).Where("id = ?", id).Updates(map[string]interface{}{"name": name, "remark": remark}).Error; err != nil {
			resultResp.Response(c, nil, nil, "更新角色失败")
			return
		}
	}

	menuIDs := parseIDList(req["menuIdList"])
	connection.DB.Self.Where("role_id = ?", id).Delete(&sysModel.RoleMenu{})
	for _, mid := range menuIDs {
		if mid == 0 {
			continue
		}
		connection.DB.Self.Create(&sysModel.RoleMenu{RoleID: id, MenuID: mid})
	}
	deptIDs := parseIDList(req["deptIdList"])
	compatMu.Lock()
	roleDeptBinding[id] = deptIDs
	compatMu.Unlock()
	resultResp.Response(c, nil, gin.H{"id": id}, "成功")
}

func RoleBatchDeleteCompatHandler(c *gin.Context) {
	ids := parseBodyIDs(c)
	if len(ids) == 0 {
		resultResp.Response(c, nil, nil, "请选择要删除的角色")
		return
	}
	connection.DB.Self.Where("id in ?", ids).Delete(&sysModel.Role{})
	connection.DB.Self.Where("role_id in ?", ids).Delete(&sysModel.RoleMenu{})
	connection.DB.Self.Where("role_id in ?", ids).Delete(&sysModel.UserRole{})
	compatMu.Lock()
	for _, id := range ids {
		delete(roleDeptBinding, id)
	}
	compatMu.Unlock()
	resultResp.Response(c, nil, nil, "成功")
}

func MenuListCompatHandler(c *gin.Context) {
	menus := loadMenuTree()
	resultResp.Response(c, nil, menus, "成功")
}

func MenuSelectCompatHandler(c *gin.Context) {
	menus := loadMenuTree()
	nodes := make([]simpleNode, 0, len(menus))
	for _, m := range menus {
		nodes = append(nodes, menuNodeToSimpleNode(m))
	}
	resultResp.Response(c, nil, nodes, "成功")
}

func menuNodeToSimpleNode(m menuNode) simpleNode {
	id, _ := strconv.ParseUint(m.ID, 10, 64)
	pid, _ := strconv.ParseUint(m.ParentID, 10, 64)
	n := simpleNode{ID: id, PID: pid, Name: m.Name}
	for _, c := range m.Children {
		n.Children = append(n.Children, menuNodeToSimpleNode(c))
	}
	return n
}

func MenuDetailCompatHandler(c *gin.Context) {
	id := parseUintFromAny(c.Param("id"))
	var menu sysModel.Menu
	if err := connection.DB.Self.Where("id = ?", id).First(&menu).Error; err != nil {
		resultResp.Response(c, nil, map[string]interface{}{}, "菜单不存在")
		return
	}
	parentName := "一级菜单"
	if menu.ParentID > 0 {
		var parent sysModel.Menu
		if err := connection.DB.Self.Where("id = ?", menu.ParentID).First(&parent).Error; err == nil {
			parentName = parent.Name
		}
	}
	resultResp.Response(c, nil, map[string]interface{}{
		"id":          menu.Id,
		"menuType":    menu.Type,
		"name":        menu.Name,
		"pid":         menu.ParentID,
		"parentName":  parentName,
		"url":         menu.Path,
		"permissions": menu.Perms,
		"sort":        menu.Sort,
		"icon":        menu.Icon,
		"openStyle":   boolToOpenStyle(menu.Path),
	}, "成功")
}

func MenuUpsertCompatHandler(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, nil, nil, "请求参数错误")
		return
	}
	id := parseUintFromAny(req["id"])
	if id == 0 {
		id = parseUintFromAny(c.Param("id"))
	}
	menu := sysModel.Menu{
		ParentID: parseUintFromAny(req["pid"]),
		Name:     parseStringFromAny(req["name"]),
		Path:     parseStringFromAny(req["url"]),
		Perms:    parseStringFromAny(req["permissions"]),
		Icon:     parseStringFromAny(req["icon"]),
		Type:     parseIntFromAny(req["menuType"]),
		Sort:     parseIntFromAny(req["sort"]),
		Visible:  1,
	}
	if menu.Name == "" {
		resultResp.Response(c, nil, nil, "菜单名称不能为空")
		return
	}
	if id == 0 {
		if err := connection.DB.Self.Create(&menu).Error; err != nil {
			resultResp.Response(c, nil, nil, "创建菜单失败")
			return
		}
	} else {
		if err := connection.DB.Self.Model(&sysModel.Menu{}).Where("id = ?", id).Updates(map[string]interface{}{
			"parent_id": menu.ParentID,
			"name":      menu.Name,
			"path":      menu.Path,
			"perms":     menu.Perms,
			"icon":      menu.Icon,
			"type":      menu.Type,
			"sort":      menu.Sort,
		}).Error; err != nil {
			resultResp.Response(c, nil, nil, "更新菜单失败")
			return
		}
	}
	resultResp.Response(c, nil, nil, "成功")
}

func MenuBatchDeleteCompatHandler(c *gin.Context) {
	ids := parseBodyIDs(c)
	if len(ids) == 0 {
		resultResp.Response(c, nil, nil, "请选择要删除的菜单")
		return
	}
	connection.DB.Self.Where("id in ?", ids).Delete(&sysModel.Menu{})
	connection.DB.Self.Where("role_id > 0 and menu_id in ?", ids).Delete(&sysModel.RoleMenu{})
	resultResp.Response(c, nil, nil, "成功")
}

func DeptListCompatHandler(c *gin.Context) {
	var depts []sysModel.Dept
	if err := connection.DB.Self.Order("sort asc, id asc").Find(&depts).Error; err != nil {
		resultResp.Response(c, nil, []interface{}{}, "成功")
		return
	}
	if len(depts) == 0 {
		resultResp.Response(c, nil, []interface{}{}, "成功")
		return
	}
	childrenByParent := make(map[uint64][]simpleNode)
	for _, d := range depts {
		node := simpleNode{ID: d.Id, PID: d.ParentID, Name: d.Name}
		childrenByParent[d.ParentID] = append(childrenByParent[d.ParentID], node)
	}
	for pid := range childrenByParent {
		sort.Slice(childrenByParent[pid], func(i, j int) bool {
			return childrenByParent[pid][i].ID < childrenByParent[pid][j].ID
		})
	}
	var attach func(uint64) []simpleNode
	attach = func(pid uint64) []simpleNode {
		nodes := childrenByParent[pid]
		for i := range nodes {
			nodes[i].Children = attach(nodes[i].ID)
		}
		return nodes
	}
	resultResp.Response(c, nil, attach(0), "成功")
}

func DeptDetailCompatHandler(c *gin.Context) {
	id := parseUintFromAny(c.Param("id"))
	var dept sysModel.Dept
	if err := connection.DB.Self.Where("id = ?", id).First(&dept).Error; err != nil {
		resultResp.Response(c, nil, map[string]interface{}{}, "部门不存在")
		return
	}
	parentName := "一级部门"
	if dept.ParentID > 0 {
		var p sysModel.Dept
		if err := connection.DB.Self.Where("id = ?", dept.ParentID).First(&p).Error; err == nil {
			parentName = p.Name
		}
	}
	resultResp.Response(c, nil, map[string]interface{}{
		"id":         dept.Id,
		"name":       dept.Name,
		"pid":        dept.ParentID,
		"parentName": parentName,
		"sort":       dept.Sort,
	}, "成功")
}

func DeptUpsertCompatHandler(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, nil, nil, "请求参数错误")
		return
	}
	id := parseUintFromAny(req["id"])
	if id == 0 {
		id = parseUintFromAny(c.Param("id"))
	}
	name := parseStringFromAny(req["name"])
	if name == "" {
		resultResp.Response(c, nil, nil, "部门名称不能为空")
		return
	}
	pid := parseUintFromAny(req["pid"])
	sortNo := parseIntFromAny(req["sort"])
	if id == 0 {
		dept := sysModel.Dept{ParentID: pid, Name: name, Sort: sortNo, Status: 1}
		if err := connection.DB.Self.Create(&dept).Error; err != nil {
			resultResp.Response(c, nil, nil, "创建部门失败")
			return
		}
	} else {
		if err := connection.DB.Self.Model(&sysModel.Dept{}).Where("id = ?", id).Updates(map[string]interface{}{"parent_id": pid, "name": name, "sort": sortNo}).Error; err != nil {
			resultResp.Response(c, nil, nil, "更新部门失败")
			return
		}
	}
	resultResp.Response(c, nil, nil, "成功")
}

func DeptBatchDeleteCompatHandler(c *gin.Context) {
	ids := parseBodyIDs(c)
	if len(ids) == 0 {
		resultResp.Response(c, nil, nil, "请选择要删除的部门")
		return
	}
	connection.DB.Self.Where("id in ?", ids).Delete(&sysModel.Dept{})
	resultResp.Response(c, nil, nil, "成功")
}

func RegionTreeCompatHandler(c *gin.Context) {
	resultResp.Response(c, nil, []interface{}{}, "成功")
}

func DictTypePageCompatHandler(c *gin.Context) {
	compatMu.RLock()
	copyItems := make([]map[string]interface{}, len(dictTypes))
	for i := range dictTypes {
		copyItems[i] = dictTypes[i]
	}
	compatMu.RUnlock()
	page, limit := pageLimit(c)
	part, total := paginateSlice(copyItems, page, limit)
	resultResp.Response(c, nil, pageResult{List: part, Total: total}, "成功")
}

func DictTypeDetailCompatHandler(c *gin.Context) {
	id := parseUintFromAny(c.Param("id"))
	compatMu.RLock()
	defer compatMu.RUnlock()
	for _, item := range dictTypes {
		if parseUintFromAny(item["id"]) == id {
			resultResp.Response(c, nil, item, "成功")
			return
		}
	}
	resultResp.Response(c, nil, map[string]interface{}{}, "成功")
}

func DictTypeUpsertCompatHandler(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, nil, nil, "请求参数错误")
		return
	}
	id := parseUintFromAny(req["id"])
	compatMu.Lock()
	defer compatMu.Unlock()
	if id == 0 {
		req["id"] = nextCompatID()
		req["createDate"] = nowText()
		dictTypes = append(dictTypes, req)
	} else {
		for i := range dictTypes {
			if parseUintFromAny(dictTypes[i]["id"]) == id {
				for k, v := range req {
					dictTypes[i][k] = v
				}
			}
		}
	}
	resultResp.Response(c, nil, nil, "成功")
}

func DictTypeDeleteCompatHandler(c *gin.Context) {
	ids := parseBodyIDs(c)
	compatMu.Lock()
	defer compatMu.Unlock()
	filtered := make([]map[string]interface{}, 0)
	for _, item := range dictTypes {
		keep := true
		for _, id := range ids {
			if parseUintFromAny(item["id"]) == id {
				keep = false
				break
			}
		}
		if keep {
			filtered = append(filtered, item)
		}
	}
	dictTypes = filtered
	resultResp.Response(c, nil, nil, "成功")
}

func DictDataPageCompatHandler(c *gin.Context) {
	dictTypeID := parseUintFromAny(c.Query("dictTypeId"))
	compatMu.RLock()
	items := make([]map[string]interface{}, 0)
	for _, item := range dictData {
		if dictTypeID == 0 || parseUintFromAny(item["dictTypeId"]) == dictTypeID {
			items = append(items, item)
		}
	}
	compatMu.RUnlock()
	page, limit := pageLimit(c)
	part, total := paginateSlice(items, page, limit)
	resultResp.Response(c, nil, pageResult{List: part, Total: total}, "成功")
}

func DictDataDetailCompatHandler(c *gin.Context) {
	id := parseUintFromAny(c.Param("id"))
	compatMu.RLock()
	defer compatMu.RUnlock()
	for _, item := range dictData {
		if parseUintFromAny(item["id"]) == id {
			resultResp.Response(c, nil, item, "成功")
			return
		}
	}
	resultResp.Response(c, nil, map[string]interface{}{}, "成功")
}

func DictDataUpsertCompatHandler(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, nil, nil, "请求参数错误")
		return
	}
	id := parseUintFromAny(req["id"])
	compatMu.Lock()
	defer compatMu.Unlock()
	if id == 0 {
		req["id"] = nextCompatID()
		req["createDate"] = nowText()
		dictData = append(dictData, req)
	} else {
		for i := range dictData {
			if parseUintFromAny(dictData[i]["id"]) == id {
				for k, v := range req {
					dictData[i][k] = v
				}
			}
		}
	}
	resultResp.Response(c, nil, nil, "成功")
}

func DictDataDeleteCompatHandler(c *gin.Context) {
	ids := parseBodyIDs(c)
	compatMu.Lock()
	defer compatMu.Unlock()
	filtered := make([]map[string]interface{}, 0)
	for _, item := range dictData {
		keep := true
		for _, id := range ids {
			if parseUintFromAny(item["id"]) == id {
				keep = false
				break
			}
		}
		if keep {
			filtered = append(filtered, item)
		}
	}
	dictData = filtered
	resultResp.Response(c, nil, nil, "成功")
}

func ParamsPageCompatHandler(c *gin.Context) {
	compatMu.RLock()
	items := make([]map[string]interface{}, len(paramsData))
	for i := range paramsData {
		items[i] = paramsData[i]
	}
	compatMu.RUnlock()
	page, limit := pageLimit(c)
	part, total := paginateSlice(items, page, limit)
	resultResp.Response(c, nil, pageResult{List: part, Total: total}, "成功")
}

func ParamsDetailCompatHandler(c *gin.Context) {
	id := parseUintFromAny(c.Param("id"))
	compatMu.RLock()
	defer compatMu.RUnlock()
	for _, item := range paramsData {
		if parseUintFromAny(item["id"]) == id {
			resultResp.Response(c, nil, item, "成功")
			return
		}
	}
	resultResp.Response(c, nil, map[string]interface{}{}, "成功")
}

func ParamsUpsertCompatHandler(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, nil, nil, "请求参数错误")
		return
	}
	id := parseUintFromAny(req["id"])
	compatMu.Lock()
	defer compatMu.Unlock()
	if id == 0 {
		req["id"] = nextCompatID()
		req["createDate"] = nowText()
		paramsData = append(paramsData, req)
	} else {
		for i := range paramsData {
			if parseUintFromAny(paramsData[i]["id"]) == id {
				for k, v := range req {
					paramsData[i][k] = v
				}
			}
		}
	}
	resultResp.Response(c, nil, nil, "成功")
}

func ParamsDeleteCompatHandler(c *gin.Context) {
	ids := parseBodyIDs(c)
	compatMu.Lock()
	defer compatMu.Unlock()
	filtered := make([]map[string]interface{}, 0)
	for _, item := range paramsData {
		keep := true
		for _, id := range ids {
			if parseUintFromAny(item["id"]) == id {
				keep = false
				break
			}
		}
		if keep {
			filtered = append(filtered, item)
		}
	}
	paramsData = filtered
	resultResp.Response(c, nil, nil, "成功")
}

func logPage(data []map[string]interface{}, c *gin.Context) {
	page, limit := pageLimit(c)
	part, total := paginateSlice(data, page, limit)
	resultResp.Response(c, nil, pageResult{List: part, Total: total}, "成功")
}

func LogLoginPageCompatHandler(c *gin.Context) {
	compatMu.RLock()
	items := make([]map[string]interface{}, len(loginLogs))
	copy(items, loginLogs)
	compatMu.RUnlock()
	logPage(items, c)
}

func LogOperationPageCompatHandler(c *gin.Context) {
	compatMu.RLock()
	items := make([]map[string]interface{}, len(operationLogs))
	copy(items, operationLogs)
	compatMu.RUnlock()
	logPage(items, c)
}

func LogErrorPageCompatHandler(c *gin.Context) {
	compatMu.RLock()
	items := make([]map[string]interface{}, len(errorLogs))
	copy(items, errorLogs)
	compatMu.RUnlock()
	logPage(items, c)
}

func LogLoginExportCompatHandler(c *gin.Context) {
	writeCSV(c, "log-login.csv", []string{"creatorName", "operation", "status", "ip", "createDate"}, [][]string{})
}

func LogOperationExportCompatHandler(c *gin.Context) {
	writeCSV(c, "log-operation.csv", []string{"creatorName", "operation", "requestUri", "status", "createDate"}, [][]string{})
}

func LogErrorExportCompatHandler(c *gin.Context) {
	writeCSV(c, "log-error.csv", []string{"requestUri", "requestMethod", "errorInfo", "createDate"}, [][]string{})
}

func SchedulePageCompatHandler(c *gin.Context) {
	compatMu.RLock()
	items := make([]map[string]interface{}, len(schedulesData))
	copy(items, schedulesData)
	compatMu.RUnlock()
	page, limit := pageLimit(c)
	part, total := paginateSlice(items, page, limit)
	resultResp.Response(c, nil, pageResult{List: part, Total: total}, "成功")
}

func ScheduleDetailCompatHandler(c *gin.Context) {
	id := parseUintFromAny(c.Param("id"))
	compatMu.RLock()
	defer compatMu.RUnlock()
	for _, item := range schedulesData {
		if parseUintFromAny(item["id"]) == id {
			resultResp.Response(c, nil, item, "成功")
			return
		}
	}
	resultResp.Response(c, nil, map[string]interface{}{}, "成功")
}

func ScheduleUpsertCompatHandler(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, nil, nil, "请求参数错误")
		return
	}
	id := parseUintFromAny(req["id"])
	compatMu.Lock()
	defer compatMu.Unlock()
	if id == 0 {
		req["id"] = nextCompatID()
		if req["status"] == nil {
			req["status"] = 1
		}
		req["createDate"] = nowText()
		schedulesData = append(schedulesData, req)
	} else {
		for i := range schedulesData {
			if parseUintFromAny(schedulesData[i]["id"]) == id {
				for k, v := range req {
					schedulesData[i][k] = v
				}
			}
		}
	}
	resultResp.Response(c, nil, nil, "成功")
}

func ScheduleDeleteCompatHandler(c *gin.Context) {
	ids := parseBodyIDs(c)
	compatMu.Lock()
	defer compatMu.Unlock()
	filtered := make([]map[string]interface{}, 0)
	for _, item := range schedulesData {
		keep := true
		for _, id := range ids {
			if parseUintFromAny(item["id"]) == id {
				keep = false
				break
			}
		}
		if keep {
			filtered = append(filtered, item)
		}
	}
	schedulesData = filtered
	resultResp.Response(c, nil, nil, "成功")
}

func ScheduleStatusCompatHandler(c *gin.Context, status int) {
	ids := parseBodyIDs(c)
	compatMu.Lock()
	defer compatMu.Unlock()
	for i := range schedulesData {
		id := parseUintFromAny(schedulesData[i]["id"])
		for _, target := range ids {
			if id == target {
				schedulesData[i]["status"] = status
			}
		}
	}
	resultResp.Response(c, nil, nil, "成功")
}

func SchedulePauseCompatHandler(c *gin.Context)  { ScheduleStatusCompatHandler(c, 0) }
func ScheduleResumeCompatHandler(c *gin.Context) { ScheduleStatusCompatHandler(c, 1) }

func ScheduleRunCompatHandler(c *gin.Context) {
	ids := parseBodyIDs(c)
	compatMu.Lock()
	defer compatMu.Unlock()
	for _, id := range ids {
		scheduleLogs = append(scheduleLogs, map[string]interface{}{
			"id":         nextCompatID(),
			"jobId":      id,
			"beanName":   fmt.Sprintf("job-%d", id),
			"params":     "",
			"status":     1,
			"times":      1,
			"error":      "",
			"createDate": nowText(),
		})
	}
	resultResp.Response(c, nil, nil, "成功")
}

func ScheduleLogPageCompatHandler(c *gin.Context) {
	compatMu.RLock()
	items := make([]map[string]interface{}, len(scheduleLogs))
	copy(items, scheduleLogs)
	compatMu.RUnlock()
	page, limit := pageLimit(c)
	part, total := paginateSlice(items, page, limit)
	resultResp.Response(c, nil, pageResult{List: part, Total: total}, "成功")
}

func ScheduleLogDetailCompatHandler(c *gin.Context) {
	id := parseUintFromAny(c.Param("id"))
	compatMu.RLock()
	defer compatMu.RUnlock()
	for _, item := range scheduleLogs {
		if parseUintFromAny(item["id"]) == id {
			resultResp.Response(c, nil, item, "成功")
			return
		}
	}
	resultResp.Response(c, nil, map[string]interface{}{"error": ""}, "成功")
}

func OssPageCompatHandler(c *gin.Context) {
	compatMu.RLock()
	items := make([]map[string]interface{}, len(ossFiles))
	copy(items, ossFiles)
	compatMu.RUnlock()
	page, limit := pageLimit(c)
	part, total := paginateSlice(items, page, limit)
	resultResp.Response(c, nil, pageResult{List: part, Total: total}, "成功")
}

func OssInfoCompatHandler(c *gin.Context) {
	compatMu.RLock()
	defer compatMu.RUnlock()
	resultResp.Response(c, nil, ossConfig, "成功")
}

func OssConfigCompatHandler(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		resultResp.Response(c, nil, nil, "请求参数错误")
		return
	}
	compatMu.Lock()
	for k, v := range req {
		ossConfig[k] = v
	}
	compatMu.Unlock()
	resultResp.Response(c, nil, nil, "成功")
}

func OssDeleteCompatHandler(c *gin.Context) {
	ids := parseBodyIDs(c)
	compatMu.Lock()
	defer compatMu.Unlock()
	filtered := make([]map[string]interface{}, 0)
	for _, item := range ossFiles {
		keep := true
		for _, id := range ids {
			if parseUintFromAny(item["id"]) == id {
				keep = false
				break
			}
		}
		if keep {
			filtered = append(filtered, item)
		}
	}
	ossFiles = filtered
	resultResp.Response(c, nil, nil, "成功")
}

func OssUploadCompatHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		resultResp.Response(c, nil, nil, "上传失败")
		return
	}
	url := "/uploads/" + file.Filename
	compatMu.Lock()
	ossFiles = append(ossFiles, map[string]interface{}{
		"id":         nextCompatID(),
		"url":        url,
		"createDate": nowText(),
	})
	compatMu.Unlock()
	resultResp.Response(c, nil, map[string]interface{}{"src": url, "url": url}, "成功")
}
