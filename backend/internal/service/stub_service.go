package service

import (
	"context"
	"errors"

	"github.com/employment-center/campus-recruitment/internal/dto/request"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/employment-center/campus-recruitment/pkg/pagination"
	"gorm.io/gorm"
)

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

