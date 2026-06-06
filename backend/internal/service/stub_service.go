package service

import (
	"context"
	"errors"

	"github.com/employment-center/campus-recruitment/internal/dto/request"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/employment-center/campus-recruitment/pkg/pagination"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type stubListService struct{}

func (stubListService) emptyPage() (interface{}, int64, int, int) {
	return []interface{}{}, 0, 1, pagination.DefaultPageSize
}

type careerTalkService struct{ repo repository.CareerTalkRepository }

func NewCareerTalkService(repo repository.CareerTalkRepository) CareerTalkService {
	return &careerTalkService{repo: repo}
}

func (s *careerTalkService) List(ctx context.Context, q *request.CareerTalkQuery, _ uint64) (interface{}, int64, int, int, error) {
	pq := &pagination.Query{Page: q.Page, PageSize: q.PageSize}
	page, pageSize := pq.Normalize()
	list, total, err := s.repo.ListPublished(ctx, q, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return list, total, page, pageSize, nil
}

func (s *careerTalkService) GetByID(ctx context.Context, id, _ uint64) (interface{}, error) {
	talk, err := s.repo.FindPublishedByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCareerTalkNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	return talk, nil
}

type jobFairService struct{ repo repository.JobFairRepository }

func NewJobFairService(repo repository.JobFairRepository) JobFairService {
	return &jobFairService{repo: repo}
}

func (s *jobFairService) List(ctx context.Context, q *request.JobFairQuery, _ uint64) (interface{}, int64, int, int, error) {
	pq := &pagination.Query{Page: q.Page, PageSize: q.PageSize}
	page, pageSize := pq.Normalize()
	list, total, err := s.repo.ListPublished(ctx, q, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return list, total, page, pageSize, nil
}

func (s *jobFairService) GetByID(ctx context.Context, id, _ uint64) (interface{}, error) {
	fair, err := s.repo.FindPublishedByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrJobFairNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	return fair, nil
}

type recommendationService struct{}

func NewRecommendationService(_ repository.RecommendationRepository) RecommendationService {
	return &recommendationService{}
}

func (s *recommendationService) List(_ context.Context, _ uint64, _ *gin.Context) (interface{}, int64, int, int, error) {
	list, total, page, pageSize := stubListService{}.emptyPage()
	return list, total, page, pageSize, nil
}

func (s *recommendationService) Dismiss(_ context.Context, _, _ uint64, _ string) error {
	return nil
}

type calendarService struct{ repo repository.CalendarRepository }

func NewCalendarService(repo repository.CalendarRepository) CalendarService {
	return &calendarService{repo: repo}
}

func (s *calendarService) List(_ context.Context, _ uint64, _ *gin.Context) (interface{}, error) {
	return []interface{}{}, nil
}

func (s *calendarService) Create(_ context.Context, _ uint64, _ *request.CreateCalendarEventRequest) (interface{}, error) {
	return map[string]interface{}{"message": "待实现"}, nil
}

func (s *calendarService) Update(_ context.Context, _, _ uint64, _ *request.UpdateCalendarEventRequest) (interface{}, error) {
	return map[string]interface{}{"message": "待实现"}, nil
}

func (s *calendarService) Delete(_ context.Context, _, _ uint64) error {
	return nil
}

type reminderService struct{}

func NewReminderService(_ repository.ReminderRepository) ReminderService {
	return &reminderService{}
}

func (s *reminderService) ListLogs(_ context.Context, _ uint64, _ *gin.Context) (interface{}, int64, int, int, error) {
	list, total, page, pageSize := stubListService{}.emptyPage()
	return list, total, page, pageSize, nil
}

type adminService struct{}

func NewAdminService(_ repository.AdminRepository) AdminService {
	return &adminService{}
}

func (s *adminService) TriggerSync(_ context.Context, _ *request.SyncTriggerRequest) (interface{}, error) {
	return map[string]interface{}{
		"taskId":    "sync-task-placeholder",
		"status":    "pending",
		"message":   "同步任务已创建，业务逻辑待实现",
	}, nil
}
