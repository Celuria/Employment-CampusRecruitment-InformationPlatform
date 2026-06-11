package service

import (
	"context"
	"errors"

	"github.com/employment-center/campus-recruitment/internal/dto/request"
	dtoresp "github.com/employment-center/campus-recruitment/internal/dto/response"
	"github.com/employment-center/campus-recruitment/internal/model"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/employment-center/campus-recruitment/pkg/pagination"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// CareerTalkService
// ---------------------------------------------------------------------------

type careerTalkService struct {
	repo        repository.CareerTalkRepository
	calendarRepo repository.CalendarRepository
}

func NewCareerTalkService(repo repository.CareerTalkRepository, calendarRepo repository.CalendarRepository) CareerTalkService {
	return &careerTalkService{repo: repo, calendarRepo: calendarRepo}
}

func (s *careerTalkService) List(ctx context.Context, q *request.CareerTalkQuery, userID uint64) (interface{}, int64, int, int, error) {
	pq := &pagination.Query{Page: q.Page, PageSize: q.PageSize}
	page, pageSize := pq.Normalize()
	list, total, err := s.repo.ListPublished(ctx, q, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	calSet := s.getCalendarSet(ctx, userID, string(model.EventTypeCareerTalk))
	vos := make([]dtoresp.CareerTalkVO, len(list))
	for i := range list {
		_, inCal := calSet[list[i].ID]
		vos[i] = dtoresp.ToCareerTalkVO(&list[i], inCal)
	}
	return vos, total, page, pageSize, nil
}

func (s *careerTalkService) GetByID(ctx context.Context, id, userID uint64) (interface{}, error) {
	talk, err := s.repo.FindPublishedByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCareerTalkNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	inCal := false
	if userID > 0 {
		inCal, _ = s.calendarRepo.ExistsByUserEvent(ctx, userID, id, string(model.EventTypeCareerTalk))
	}
	return dtoresp.ToCareerTalkVO(talk, inCal), nil
}

func (s *careerTalkService) ListUpcomingWithin24h(ctx context.Context) (interface{}, error) {
	list, err := s.repo.ListUpcomingWithin24h(ctx, 20)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	return list, nil
}

func (s *careerTalkService) ListHotCompanies(ctx context.Context, limit int) (interface{}, error) {
	if limit <= 0 || limit > 20 {
		limit = 6
	}
	list, err := s.repo.ListHotCompanies(ctx, limit)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	return list, nil
}

// getCalendarSet 返回用户已加入日历的 refId 集合（指定 eventType）
func (s *careerTalkService) getCalendarSet(ctx context.Context, userID uint64, eventType string) map[uint64]struct{} {
	if userID == 0 {
		return nil
	}
	ids, err := s.calendarRepo.ListRefIDsByUser(ctx, userID, eventType)
	if err != nil {
		return nil
	}
	set := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return set
}

// ---------------------------------------------------------------------------
// JobFairService
// ---------------------------------------------------------------------------

type jobFairService struct {
	repo        repository.JobFairRepository
	calendarRepo repository.CalendarRepository
}

func NewJobFairService(repo repository.JobFairRepository, calendarRepo repository.CalendarRepository) JobFairService {
	return &jobFairService{repo: repo, calendarRepo: calendarRepo}
}

func (s *jobFairService) List(ctx context.Context, q *request.JobFairQuery, userID uint64) (interface{}, int64, int, int, error) {
	pq := &pagination.Query{Page: q.Page, PageSize: q.PageSize}
	page, pageSize := pq.Normalize()
	list, total, err := s.repo.ListPublished(ctx, q, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	calSet := s.getCalendarSet(ctx, userID, string(model.EventTypeJobFair))
	vos := make([]dtoresp.JobFairVO, len(list))
	for i := range list {
		_, inCal := calSet[list[i].ID]
		vos[i] = dtoresp.ToJobFairVO(&list[i], inCal)
	}
	return vos, total, page, pageSize, nil
}

func (s *jobFairService) GetByID(ctx context.Context, id, userID uint64) (interface{}, error) {
	fair, err := s.repo.FindPublishedByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrJobFairNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	inCal := false
	if userID > 0 {
		inCal, _ = s.calendarRepo.ExistsByUserEvent(ctx, userID, id, string(model.EventTypeJobFair))
	}
	return dtoresp.ToJobFairVO(fair, inCal), nil
}

func (s *jobFairService) getCalendarSet(ctx context.Context, userID uint64, eventType string) map[uint64]struct{} {
	if userID == 0 {
		return nil
	}
	ids, err := s.calendarRepo.ListRefIDsByUser(ctx, userID, eventType)
	if err != nil {
		return nil
	}
	set := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return set
}

// 确保 fmt 被引用（如果不用可删）