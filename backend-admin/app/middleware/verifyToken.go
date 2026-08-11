package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"golang-common-base/app/models/sys"
	"golang-common-base/pkg/connection"
	"golang-common-base/pkg/logger"
	"golang-common-base/pkg/utils"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func CheckToken() func(context *gin.Context) {
	return func(context *gin.Context) {
		if !isMustApi(context) {
			context.Next()
			return
		}

		tokenString := normalizeBearer(context.GetHeader(utils.TokenNameInHeader))
		if tokenString == "" {
			context.JSON(http.StatusUnauthorized, map[string]interface{}{"errno": 100001, "errmsg": "缺少 access token"})
			context.Abort()
			return
		}

		if !authenticate(context, tokenString) {
			context.JSON(http.StatusUnauthorized, map[string]interface{}{"errno": 100001, "errmsg": "access token无效"})
			context.Abort()
			return
		}

		if !authorize(context) {
			context.JSON(http.StatusForbidden, map[string]interface{}{"errno": 100403, "errmsg": "无权限访问"})
			context.Abort()
			return
		}

		context.Next()
	}
}

func authenticate(context *gin.Context, tokenString string) bool {
	if connection.OIDCVerifier == nil {
		fallbackUser := "anonymous"
		if !viper.GetBool("auth.casdoor.enabled") {
			bootstrap := strings.TrimSpace(viper.GetString("authz.casbin.bootstrap.subject"))
			if bootstrap != "" {
				fallbackUser = bootstrap
			}
		}
		logger.Warnf("OIDC verifier 未启用，跳过 token 验证，使用用户: %s", fallbackUser)
		context.Set("user", fallbackUser)
		return true
	}

	idToken, err := connection.OIDCVerifier.Verify(context.Request.Context(), tokenString)
	if err != nil {
		logger.Errorf("OIDC token 验证失败: %v", err)
		return false
	}

	var claims struct {
		Sub               string   `json:"sub"`
		PreferredUsername string   `json:"preferred_username"`
		Name              string   `json:"name"`
		Email             string   `json:"email"`
		Roles             []string `json:"roles"`
	}

	if err = idToken.Claims(&claims); err != nil {
		logger.Errorf("读取 token claims 失败: %v", err)
		return false
	}

	username := claims.PreferredUsername
	if username == "" {
		username = claims.Name
	}
	if username == "" {
		username = claims.Sub
	}

	context.Set("user", username)
	context.Set("sub", claims.Sub)
	context.Set("email", claims.Email)
	context.Set("roles", claims.Roles)
	return true
}

func authorize(context *gin.Context) bool {
	if !viper.GetBool("authz.casbin.enabled") {
		return true
	}

	if connection.CasbinEnforcer == nil {
		logger.Warn("Casbin enforcer 未初始化，拒绝请求")
		return false
	}

	subject := context.GetString("user")
	object := context.Request.URL.Path
	action := strings.ToLower(context.Request.Method)

	for _, candidate := range authorizeSubjects(context, subject) {
		allowed, err := connection.CasbinEnforcer.Enforce(candidate, object, action)
		if err != nil {
			logger.Errorf("Casbin 鉴权失败: %v", err)
			return false
		}
		if allowed {
			return true
		}
	}

	return false
}

func authorizeSubjects(context *gin.Context, subject string) []string {
	subjects := make([]string, 0, 4)
	if subject != "" {
		subjects = append(subjects, subject)
	}

	for _, role := range context.GetStringSlice("roles") {
		role = strings.TrimSpace(role)
		if role != "" {
			subjects = append(subjects, role)
		}
	}

	if subject == "" {
		subject = "anonymous"
	}

	var user sys.User
	if err := connection.DB.Self.Where("username = ?", subject).First(&user).Error; err == nil && user.Id != 0 {
		var rows []sys.UserRole
		if err := connection.DB.Self.Where("user_id = ?", user.Id).Find(&rows).Error; err == nil {
			roleIDs := make([]uint64, 0, len(rows))
			for _, row := range rows {
				roleIDs = append(roleIDs, row.RoleID)
			}
			if len(roleIDs) > 0 {
				var roles []sys.Role
				if err := connection.DB.Self.Where("id IN ?", roleIDs).Find(&roles).Error; err == nil {
					for _, role := range roles {
						if role.Code != "" {
							subjects = append(subjects, role.Code)
						}
					}
				}
			}
		}
	}

	return uniqueStrings(subjects)
}

func uniqueStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

// 定义无需登陆检测的接口
func isMustApi(context *gin.Context) bool {
	return context.Request.URL.Path != "/api/v1/login" &&
		context.Request.URL.Path != "/api/v1/logout" &&
		context.Request.URL.Path != "/api/v1/health" &&
		!ignoreMatchErr(`/api/v1/reservations/([0-9]+)/checkin`, context.Request.URL.Path) &&
		!ignoreMatchErr(`/api/v1/upload`, context.Request.URL.Path)
}

func ignoreMatchErr(pattern, str string) bool {
	match, _ := regexp.MatchString(pattern, str)
	return match
}

func normalizeBearer(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	const prefix = "Bearer "
	if strings.HasPrefix(raw, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(raw, prefix))
	}
	return raw
}
