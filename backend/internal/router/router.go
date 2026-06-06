package router

import (
	"github.com/employment-center/campus-recruitment/config"
	"github.com/employment-center/campus-recruitment/internal/handler"
	"github.com/employment-center/campus-recruitment/internal/middleware"
	"github.com/employment-center/campus-recruitment/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// Setup 注册所有路由
func Setup(cfg *config.Config, h *handler.Handler, jwtManager *jwt.Manager) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS(cfg.CORS))
	r.Use(middleware.AuthErrorHandler())
	r.Use(middleware.ErrorHandler())

	r.GET("/health", h.Health.Check)

	v1 := r.Group("/api/v1")
	registerPublicRoutes(v1, h)
	registerOptionalAuthRoutes(v1, h, jwtManager)

	auth := v1.Group("", middleware.Auth(jwtManager), middleware.RequireAuth())
	registerStudentRoutes(auth, h)

	admin := v1.Group("/admin", middleware.Auth(jwtManager), middleware.RequireAuth(), middleware.RequireAdmin())
	registerAdminRoutes(admin, h)

	return r
}

func registerPublicRoutes(v1 *gin.RouterGroup, h *handler.Handler) {
	auth := v1.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
		auth.GET("/captcha", h.Auth.Captcha)
	}
}

func registerOptionalAuthRoutes(v1 *gin.RouterGroup, h *handler.Handler, jwtManager *jwt.Manager) {
	optional := v1.Group("", middleware.OptionalAuth(jwtManager))
	{
		optional.GET("/career-talks", h.CareerTalk.List)
		optional.GET("/career-talks/:id", h.CareerTalk.Detail)
		optional.GET("/job-fairs", h.JobFair.List)
		optional.GET("/job-fairs/:id", h.JobFair.Detail)
	}
}

func registerStudentRoutes(auth *gin.RouterGroup, h *handler.Handler) {
	auth.POST("/auth/logout", h.Auth.Logout)

	auth.GET("/users/me", h.User.GetProfile)
	auth.PUT("/users/me", h.User.UpdateProfile)
	auth.GET("/users/me/preferences", h.User.GetPreferences)
	auth.PUT("/users/me/preferences", h.User.UpdatePreferences)

	auth.GET("/recommendations", h.Recommendation.List)
	auth.POST("/recommendations/:refId/dismiss", h.Recommendation.Dismiss)

	auth.GET("/calendar/events", h.Calendar.List)
	auth.POST("/calendar/events", h.Calendar.Create)
	auth.PUT("/calendar/events/:id", h.Calendar.Update)
	auth.DELETE("/calendar/events/:id", h.Calendar.Delete)

	auth.GET("/reminders/logs", h.Reminder.ListLogs)
}

func registerAdminRoutes(admin *gin.RouterGroup, h *handler.Handler) {
	admin.POST("/sync", h.Admin.Sync)

	admin.GET("/career-talks", h.Admin.Placeholder)
	admin.POST("/career-talks", h.Admin.Placeholder)
	admin.PUT("/career-talks/:id", h.Admin.Placeholder)
	admin.DELETE("/career-talks/:id", h.Admin.Placeholder)

	admin.GET("/job-fairs", h.Admin.Placeholder)
	admin.POST("/job-fairs", h.Admin.Placeholder)
	admin.PUT("/job-fairs/:id", h.Admin.Placeholder)
	admin.DELETE("/job-fairs/:id", h.Admin.Placeholder)

	admin.GET("/users", h.Admin.Placeholder)
	admin.GET("/audit-logs", h.Admin.Placeholder)
}
