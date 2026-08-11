package sys

import (
	"sort"
	"strconv"
	"strings"

	sysModel "golang-common-base/app/models/sys"
	"golang-common-base/pkg/connection"
	resultResp "golang-common-base/pkg/response/response"

	"github.com/gin-gonic/gin"
)

type menuNode struct {
	ID        string     `json:"id"`
	ParentID  string     `json:"parent_id"`
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	Icon      string     `json:"icon"`
	OpenStyle int        `json:"openStyle"`
	Type      int        `json:"type"`
	Sort      int        `json:"sort"`
	Visible   int        `json:"visible"`
	Children  []menuNode `json:"children,omitempty"`
}

func UserInfoHandler(c *gin.Context) {
	username := strings.TrimSpace(c.GetString("user"))
	if username == "" {
		username = "admin"
	}

	var user sysModel.User
	if err := connection.DB.Self.Where("username = ?", username).First(&user).Error; err != nil {
		resultResp.Response(c, nil, gin.H{
			"createDate": "",
			"deptId":     "0",
			"deptName":   "",
			"email":      "",
			"gender":     0,
			"headUrl":    "",
			"id":         "0",
			"mobile":     "",
			"postIdList": "",
			"realName":   username,
			"roleIdList": "1",
			"status":     1,
			"superAdmin": 1,
			"username":   username,
		}, "成功获取用户信息")
		return
	}

	roleIDs, roleCodes := userRoleBindings(user.Id)
	deptName := ""
	if user.DeptID != 0 {
		var dept sysModel.Dept
		if err := connection.DB.Self.Where("id = ?", user.DeptID).First(&dept).Error; err == nil {
			deptName = dept.Name
		}
	}

	resultResp.Response(c, nil, gin.H{
		"createDate": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"deptId":     strconv.FormatUint(user.DeptID, 10),
		"deptName":   deptName,
		"email":      user.Email,
		"gender":     0,
		"headUrl":    "",
		"id":         strconv.FormatUint(user.Id, 10),
		"mobile":     user.Mobile,
		"postIdList": "",
		"realName":   user.Nickname,
		"roleIdList": joinUint64(roleIDs),
		"roleCodes":  strings.Join(roleCodes, ","),
		"status":     user.Status,
		"superAdmin": boolToInt(username == "admin" || containsString(roleCodes, "admin")),
		"username":   user.Username,
	}, "成功获取用户信息")
}

func MenuNavHandler(c *gin.Context) {
	menus := loadMenuTree()
	resultResp.Response(c, nil, menus, "成功获取菜单导航")
}

func PermissionListHandler(c *gin.Context) {
	perms := loadPermissions()
	resultResp.Response(c, nil, perms, "成功获取权限列表")
}

func DictTypeAllHandler(c *gin.Context) {
	dicts := []gin.H{
		{
			"dictType": "sys_normal_disable",
			"dataList": []gin.H{
				{"dictLabel": "正常", "dictValue": "0", "listClass": "success", "cssClass": ""},
				{"dictLabel": "停用", "dictValue": "1", "listClass": "danger", "cssClass": ""},
			},
		},
		{
			"dictType": "sys_yes_no",
			"dataList": []gin.H{
				{"dictLabel": "是", "dictValue": "1", "listClass": "success", "cssClass": ""},
				{"dictLabel": "否", "dictValue": "0", "listClass": "info", "cssClass": ""},
			},
		},
	}
	resultResp.Response(c, nil, dicts, "成功获取字典列表")
}

func loadMenuTree() []menuNode {
	var rows []sysModel.Menu
	if err := connection.DB.Self.Order("sort asc, id asc").Find(&rows).Error; err == nil && len(rows) > 0 {
		return normalizeMenuTree(buildMenuTree(rows))
	}
	return normalizeMenuTree(staticMenuTree())
}

func loadPermissions() []string {
	var rows []sysModel.Menu
	if err := connection.DB.Self.Where("perms <> ''").Order("sort asc, id asc").Find(&rows).Error; err == nil && len(rows) > 0 {
		perms := make([]string, 0, len(rows))
		seen := make(map[string]struct{}, len(rows))
		for _, row := range rows {
			for _, item := range strings.Split(row.Perms, ",") {
				item = strings.TrimSpace(item)
				if item == "" {
					continue
				}
				if _, ok := seen[item]; ok {
					continue
				}
				seen[item] = struct{}{}
				perms = append(perms, item)
			}
		}
		if len(perms) > 0 {
			return perms
		}
	}
	return []string{
		"sys:user:save",
		"sys:user:update",
		"sys:user:delete",
		"sys:user:export",
		"sys:role:save",
		"sys:role:update",
		"sys:role:delete",
		"sys:menu:save",
		"sys:menu:update",
		"sys:menu:delete",
		"sys:dept:save",
		"sys:dept:update",
		"sys:dept:delete",
		"sys:dict:save",
		"sys:dict:update",
		"sys:dict:delete",
		"sys:params:save",
		"sys:params:update",
		"sys:params:delete",
		"sys:schedule:save",
		"sys:schedule:update",
		"sys:schedule:delete",
		"sys:schedule:pause",
		"sys:schedule:resume",
		"sys:schedule:run",
		"sys:schedule:log",
	}
}

func userRoleBindings(userID uint64) ([]uint64, []string) {
	if userID == 0 {
		return nil, nil
	}
	var rows []sysModel.UserRole
	if err := connection.DB.Self.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, nil
	}
	roleIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		roleIDs = append(roleIDs, row.RoleID)
	}
	if len(roleIDs) == 0 {
		return roleIDs, nil
	}
	var roles []sysModel.Role
	if err := connection.DB.Self.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return roleIDs, nil
	}
	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		if role.Code != "" {
			roleCodes = append(roleCodes, role.Code)
		}
	}
	return roleIDs, roleCodes
}

func buildMenuTree(rows []sysModel.Menu) []menuNode {
	childrenByParent := make(map[uint64][]menuNode)
	for _, row := range rows {
		node := toMenuNode(row)
		childrenByParent[row.ParentID] = append(childrenByParent[row.ParentID], node)
	}
	for parentID := range childrenByParent {
		sort.Slice(childrenByParent[parentID], func(i, j int) bool {
			if childrenByParent[parentID][i].Sort == childrenByParent[parentID][j].Sort {
				return childrenByParent[parentID][i].ID < childrenByParent[parentID][j].ID
			}
			return childrenByParent[parentID][i].Sort < childrenByParent[parentID][j].Sort
		})
	}
	return attachChildren(0, childrenByParent)
}

func attachChildren(parentID uint64, childrenByParent map[uint64][]menuNode) []menuNode {
	children := childrenByParent[parentID]
	for idx := range children {
		childID, _ := strconv.ParseUint(children[idx].ID, 10, 64)
		children[idx].Children = attachChildren(childID, childrenByParent)
	}
	return children
}

func toMenuNode(row sysModel.Menu) menuNode {
	return menuNode{
		ID:        strconv.FormatUint(row.Id, 10),
		ParentID:  strconv.FormatUint(row.ParentID, 10),
		Name:      normalizeMenuName(row.Path, row.Name),
		URL:       row.Path,
		Icon:      row.Icon,
		OpenStyle: boolToOpenStyle(row.Path),
		Type:      row.Type,
		Sort:      row.Sort,
		Visible:   row.Visible,
	}
}

func normalizeMenuTree(nodes []menuNode) []menuNode {
	normalized := make([]menuNode, 0, len(nodes))
	for _, node := range nodes {
		node.Children = normalizeMenuTree(node.Children)
		node.Name = normalizeMenuName(node.URL, node.Name)

		if len(node.Children) == 1 {
			onlyChild := node.Children[0]
			sameName := strings.EqualFold(strings.TrimSpace(node.Name), strings.TrimSpace(onlyChild.Name))
			sameURL := strings.EqualFold(strings.TrimSpace(node.URL), strings.TrimSpace(onlyChild.URL))
			if sameName || sameURL {
				onlyChild.ParentID = node.ParentID
				if onlyChild.Sort == 0 {
					onlyChild.Sort = node.Sort
				}
				node = onlyChild
			}
		}

		normalized = append(normalized, node)
	}
	return normalized
}

func normalizeMenuName(path string, name string) string {
	path = strings.TrimSpace(path)
	name = strings.TrimSpace(name)
	if path == "" {
		return name
	}

	nameByPath := map[string]string{
		"/home":             "工作台",
		"/job/schedule":     "任务调度",
		"/job/schedule-log": "调度日志",
		"/oss/oss":          "文件管理",
		"/oss/oss-config":   "存储配置",
		"/sys/dict-type":    "字典类型",
		"/sys/dict-data":    "字典数据",
	}
	if v, ok := nameByPath[path]; ok {
		return v
	}
	if name != "" {
		return name
	}
	return path
}

func staticMenuTree() []menuNode {
	return []menuNode{
		{ID: "1", ParentID: "0", Name: "工作台", URL: "/home", Icon: "icon-dashboard-fill", OpenStyle: 0, Type: 0, Sort: 0, Visible: 1},
		{
			ID:        "2",
			ParentID:  "0",
			Name:      "系统管理",
			URL:       "/sys",
			Icon:      "icon-setting",
			OpenStyle: 0,
			Type:      0,
			Sort:      1,
			Visible:   1,
			Children: []menuNode{
				{ID: "21", ParentID: "2", Name: "用户管理", URL: "/sys/user", Icon: "icon-user", OpenStyle: 0, Type: 0, Sort: 1, Visible: 1},
				{ID: "22", ParentID: "2", Name: "角色管理", URL: "/sys/role", Icon: "icon-roles", OpenStyle: 0, Type: 0, Sort: 2, Visible: 1},
				{ID: "23", ParentID: "2", Name: "菜单管理", URL: "/sys/menu", Icon: "icon-menu", OpenStyle: 0, Type: 0, Sort: 3, Visible: 1},
				{ID: "24", ParentID: "2", Name: "部门管理", URL: "/sys/dept", Icon: "icon-dept", OpenStyle: 0, Type: 0, Sort: 4, Visible: 1},
				{ID: "25", ParentID: "2", Name: "字典类型", URL: "/sys/dict-type", Icon: "icon-dictionary", OpenStyle: 0, Type: 0, Sort: 5, Visible: 1},
				{ID: "26", ParentID: "2", Name: "字典数据", URL: "/sys/dict-data", Icon: "icon-dictionary", OpenStyle: 0, Type: 0, Sort: 6, Visible: 1},
				{ID: "27", ParentID: "2", Name: "参数管理", URL: "/sys/params", Icon: "icon-params", OpenStyle: 0, Type: 0, Sort: 7, Visible: 1},
				{ID: "28", ParentID: "2", Name: "登录日志", URL: "/sys/log-login", Icon: "icon-log", OpenStyle: 0, Type: 0, Sort: 8, Visible: 1},
				{ID: "29", ParentID: "2", Name: "操作日志", URL: "/sys/log-operation", Icon: "icon-log", OpenStyle: 0, Type: 0, Sort: 9, Visible: 1},
				{ID: "30", ParentID: "2", Name: "异常日志", URL: "/sys/log-error", Icon: "icon-log", OpenStyle: 0, Type: 0, Sort: 10, Visible: 1},
			},
		},
		{
			ID:        "3",
			ParentID:  "0",
			Name:      "任务调度",
			URL:       "/job",
			Icon:      "icon-calculator",
			OpenStyle: 0,
			Type:      0,
			Sort:      2,
			Visible:   1,
			Children: []menuNode{
				{ID: "31", ParentID: "3", Name: "任务列表", URL: "/job/schedule", Icon: "icon-schedule", OpenStyle: 0, Type: 0, Sort: 1, Visible: 1},
				{ID: "32", ParentID: "3", Name: "执行日志", URL: "/job/schedule-log", Icon: "icon-log", OpenStyle: 0, Type: 0, Sort: 2, Visible: 1},
			},
		},
		{
			ID:        "4",
			ParentID:  "0",
			Name:      "对象存储",
			URL:       "/oss",
			Icon:      "icon-wallet",
			OpenStyle: 0,
			Type:      0,
			Sort:      3,
			Visible:   1,
			Children: []menuNode{
				{ID: "41", ParentID: "4", Name: "文件列表", URL: "/oss/oss", Icon: "icon-oss", OpenStyle: 0, Type: 0, Sort: 1, Visible: 1},
				{ID: "42", ParentID: "4", Name: "存储配置", URL: "/oss/oss-config", Icon: "icon-setting", OpenStyle: 0, Type: 0, Sort: 2, Visible: 1},
			},
		},
	}
}

func boolToOpenStyle(path string) int {
	if strings.HasPrefix(strings.TrimSpace(path), "http://") || strings.HasPrefix(strings.TrimSpace(path), "https://") {
		return 1
	}
	return 0
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func joinUint64(values []uint64) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, strconv.FormatUint(value, 10))
	}
	return strings.Join(parts, ",")
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), needle) {
			return true
		}
	}
	return false
}
