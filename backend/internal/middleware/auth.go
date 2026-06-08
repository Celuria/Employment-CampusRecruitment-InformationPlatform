package middleware

import (
	"net/http"
	"strings"

	"github.com/employment-center/campus-recruitment/config"
	"github.com/employment-center/campus-recruitment/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey   = "userID"
	ContextUsernameKey = "username"
	ContextUserRoleKey = "userRole"
)

// CORS 跨域中间件
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		for _, allowed := range cfg.AllowOrigins {
			if allowed == "*" || allowed == origin {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}
		c.Header("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ", "))
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Auth JWT 鉴权中间件（必须携带有效 Token）
func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !parseToken(c, jwtManager) {
			c.Abort()
			c.Set("auth_error", true)
			return
		}
		c.Next()
	}
}

// OptionalAuth 可选鉴权：有 Token 则解析用户信息，无 Token 也放行
func OptionalAuth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = parseToken(c, jwtManager)
		c.Next()
	}
}

func parseToken(c *gin.Context, jwtManager *jwt.Manager) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwtManager.Parse(tokenStr)
	if err != nil {
		return false
	}
	c.Set(ContextUserIDKey, claims.UserID)
	c.Set(ContextUsernameKey, claims.Username)
	c.Set(ContextUserRoleKey, claims.Role)
	return true
}

// RequireAuth 必须登录
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get(ContextUserIDKey); !exists {
			c.Abort()
			c.Set("auth_error", true)
			return
		}
		c.Next()
	}
}

// RequireAdmin 必须管理员
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ContextUserRoleKey)
		if role != string(modelRoleAdmin) {
			c.Abort()
			c.Set("admin_forbidden", true)
			return
		}
		c.Next()
	}
}

const modelRoleAdmin = "admin"

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) (uint64, bool) {
	v, ok := c.Get(ContextUserIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint64)
	return id, ok
}

// GetUserRole 从上下文获取用户角色
func GetUserRole(c *gin.Context) string {
	v, _ := c.Get(ContextUserRoleKey)
	role, _ := v.(string)
	return role
}
