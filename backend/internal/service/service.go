package service

import (
	"context"

	"github.com/employment-center/campus-recruitment/config"
	"github.com/employment-center/campus-recruitment/internal/dto/request"
	dtoresp "github.com/employment-center/campus-recruitment/internal/dto/response"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/jwt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Services 聚合所有业务服务，便于依赖注入
type Services struct {
	Auth           AuthService
	User           UserService
	CareerTalk     CareerTalkService
	JobFair        JobFairService
	Recommendation RecommendationService
	Calendar       CalendarService
	Reminder       ReminderService
	Admin          AdminService
}

func NewServices(db *gorm.DB, cfg *config.Config, jwtManager *jwt.Manager) *Services {
	repos := repository.NewRepositories(db)

	return &Services{
		Auth:           NewAuthService(repos.User, jwtManager, cfg.Auth),
		User:           NewUserService(repos.User, repos.Preference),
		CareerTalk:     NewCareerTalkService(repos.CareerTalk),
		JobFair:        NewJobFairService(repos.JobFair),
		Recommendation: NewRecommendationService(repos.User, repos.Preference, repos.CareerTalk, repos.JobFair),
		Calendar:       NewCalendarService(repos.Calendar, repos.CareerTalk, repos.JobFair, repos.Preference),
		Reminder:       NewReminderService(repos.Reminder),
		Admin:          NewAdminService(repos.Admin),
	}
}

// AuthService 认证服务
type AuthService interface {
	Register(ctx context.Context, req *request.RegisterRequest) error
	Login(ctx context.Context, req *request.LoginRequest) (*dtoresp.LoginResponse, error)
}

// UserService 用户服务
type UserService interface {
	GetProfile(ctx context.Context, userID uint64) (*dtoresp.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint64, req *request.UpdateProfileRequest) (*dtoresp.UserProfileResponse, error)
	GetPreferences(ctx context.Context, userID uint64) (interface{}, error)
	UpdatePreferences(ctx context.Context, userID uint64, req *request.UpdatePreferenceRequest) (interface{}, error)
}

// CareerTalkService 宣讲会服务
type CareerTalkService interface {
	List(ctx context.Context, q *request.CareerTalkQuery, userID uint64) (list interface{}, total int64, page, pageSize int, err error)
	GetByID(ctx context.Context, id, userID uint64) (interface{}, error)
	ListUpcomingWithin24h(ctx context.Context) (interface{}, error)
	ListHotCompanies(ctx context.Context, limit int) (interface{}, error)
}

// JobFairService 双选会服务
type JobFairService interface {
	List(ctx context.Context, q *request.JobFairQuery, userID uint64) (list interface{}, total int64, page, pageSize int, err error)
	GetByID(ctx context.Context, id, userID uint64) (interface{}, error)
}

// RecommendationService 推荐服务
type RecommendationService interface {
	List(ctx context.Context, userID uint64, c *gin.Context) (*dtoresp.RecommendationListResult, error)
	Dismiss(ctx context.Context, userID, refID uint64, eventType string) error
}

// CalendarService 日历服务
type CalendarService interface {
	List(ctx context.Context, userID uint64, c *gin.Context) (interface{}, error)
	Create(ctx context.Context, userID uint64, req *request.CreateCalendarEventRequest) (interface{}, error)
	Update(ctx context.Context, userID, id uint64, req *request.UpdateCalendarEventRequest) (interface{}, error)
	Delete(ctx context.Context, userID, id uint64) error
}

// ReminderService 提醒服务
type ReminderService interface {
	ListLogs(ctx context.Context, userID uint64, c *gin.Context) (list interface{}, total int64, page, pageSize int, err error)
}

// AdminService 管理端服务
type AdminService interface {
	TriggerSync(ctx context.Context, req *request.SyncTriggerRequest) (interface{}, error)
}
