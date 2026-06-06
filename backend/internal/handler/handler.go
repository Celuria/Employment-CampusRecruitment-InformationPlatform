package handler

import (
	"net/http"

	dtoresp "github.com/employment-center/campus-recruitment/internal/dto/response"
	"github.com/employment-center/campus-recruitment/internal/service"
	"github.com/employment-center/campus-recruitment/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handler 聚合所有 HTTP 处理器
type Handler struct {
	Health         *HealthHandler
	Auth           *AuthHandler
	User           *UserHandler
	CareerTalk     *CareerTalkHandler
	JobFair        *JobFairHandler
	Recommendation *RecommendationHandler
	Calendar       *CalendarHandler
	Reminder       *ReminderHandler
	Admin          *AdminHandler
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		Health:         NewHealthHandler(),
		Auth:           NewAuthHandler(services.Auth),
		User:           NewUserHandler(services.User),
		CareerTalk:     NewCareerTalkHandler(services.CareerTalk),
		JobFair:        NewJobFairHandler(services.JobFair),
		Recommendation: NewRecommendationHandler(services.Recommendation),
		Calendar:       NewCalendarHandler(services.Calendar),
		Reminder:       NewReminderHandler(services.Reminder),
		Admin:          NewAdminHandler(services.Admin),
	}
}

// HealthHandler 健康检查
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) Check(c *gin.Context) {
	response.OK(c, dtoresp.HealthResponse{Status: "ok", Version: "0.1.0"})
}

func (h *HealthHandler) NotImplemented(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, 50100, "接口框架已就绪，业务逻辑待实现")
}
