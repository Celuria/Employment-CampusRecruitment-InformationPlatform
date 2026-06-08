package service

import (
	"context"
	"errors"
	"time"

	"github.com/employment-center/campus-recruitment/internal/dto/request"
	dtoresp "github.com/employment-center/campus-recruitment/internal/dto/response"
	"github.com/employment-center/campus-recruitment/internal/model"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type calendarService struct {
	repo         repository.CalendarRepository
	careerTalk   repository.CareerTalkRepository
	jobFair      repository.JobFairRepository
	preference   repository.PreferenceRepository
}

func NewCalendarService(
	repo repository.CalendarRepository,
	careerTalk repository.CareerTalkRepository,
	jobFair repository.JobFairRepository,
	preference repository.PreferenceRepository,
) CalendarService {
	return &calendarService{
		repo:       repo,
		careerTalk: careerTalk,
		jobFair:    jobFair,
		preference: preference,
	}
}

func (s *calendarService) List(ctx context.Context, userID uint64, c *gin.Context) (interface{}, error) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	eventType := c.Query("eventType")

	events, err := s.repo.ListByUser(ctx, userID, startDate, endDate, eventType)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	vos := make([]dtoresp.CalendarEventVO, 0, len(events))
	for i := range events {
		vos = append(vos, s.toVO(ctx, &events[i]))
	}
	return vos, nil
}

func (s *calendarService) Create(ctx context.Context, userID uint64, req *request.CreateCalendarEventRequest) (interface{}, error) {
	exists, err := s.repo.ExistsByUserEvent(ctx, userID, req.RefID, req.EventType)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	if exists {
		return nil, apperrors.ErrCalendarDuplicate
	}

	event, err := s.buildEventFromSource(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return nil, apperrors.ErrInternalServer
	}

	vo := s.toVO(ctx, event)
	return vo, nil
}

func (s *calendarService) Update(ctx context.Context, userID, id uint64, req *request.UpdateCalendarEventRequest) (interface{}, error) {
	event, err := s.repo.FindByID(ctx, userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCalendarNotFound
		}
		return nil, apperrors.ErrInternalServer
	}

	remindBefore := event.RemindBefore
	if req.RemindBefore != nil {
		remindBefore = model.JSONStrings(req.RemindBefore)
	}

	if err := s.repo.Update(ctx, userID, id, req.CustomNote, remindBefore); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCalendarNotFound
		}
		return nil, apperrors.ErrInternalServer
	}

	event.CustomNote = req.CustomNote
	event.RemindBefore = remindBefore
	return s.toVO(ctx, event), nil
}

func (s *calendarService) Delete(ctx context.Context, userID, id uint64) error {
	if err := s.repo.Delete(ctx, userID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrCalendarNotFound
		}
		return apperrors.ErrInternalServer
	}
	return nil
}

func (s *calendarService) buildEventFromSource(ctx context.Context, userID uint64, req *request.CreateCalendarEventRequest) (*model.CalendarEvent, error) {
	remindBefore := s.resolveRemindBefore(ctx, userID, req.RemindBefore)

	switch req.EventType {
	case string(model.EventTypeCareerTalk):
		talk, err := s.careerTalk.FindPublishedByID(ctx, req.RefID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperrors.ErrEventNotFound
			}
			return nil, apperrors.ErrInternalServer
		}
		event := &model.CalendarEvent{
			UserID:       userID,
			EventType:    model.EventTypeCareerTalk,
			RefID:        talk.ID,
			Title:        talk.Title,
			StartTime:    talk.StartTime,
			Location:     talk.Location,
			CustomNote:   req.CustomNote,
			RemindBefore: remindBefore,
			Status:       "active",
		}
		if talk.EndTime != nil {
			event.EndTime = talk.EndTime
		}
		return event, nil

	case string(model.EventTypeJobFair):
		fair, err := s.jobFair.FindPublishedByID(ctx, req.RefID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperrors.ErrEventNotFound
			}
			return nil, apperrors.ErrInternalServer
		}
		startTime := fair.StartDate
		if fair.StartTime != nil {
			startTime = *fair.StartTime
		}
		event := &model.CalendarEvent{
			UserID:       userID,
			EventType:    model.EventTypeJobFair,
			RefID:        fair.ID,
			Title:        fair.Title,
			StartTime:    startTime,
			Location:     fair.Location,
			CustomNote:   req.CustomNote,
			RemindBefore: remindBefore,
			Status:       "active",
		}
		if fair.EndDate != nil {
			end := time.Date(fair.EndDate.Year(), fair.EndDate.Month(), fair.EndDate.Day(), 23, 59, 59, 0, fair.EndDate.Location())
			event.EndTime = &end
		}
		return event, nil

	default:
		return nil, apperrors.ErrBadRequest
	}
}

func (s *calendarService) resolveRemindBefore(ctx context.Context, userID uint64, requested []string) model.JSONStrings {
	if len(requested) > 0 {
		return model.JSONStrings(requested)
	}
	pref, err := s.preference.FindByUserID(ctx, userID)
	if err == nil && len(pref.RemindBefore) > 0 {
		return pref.RemindBefore
	}
	return model.JSONStrings{"1d"}
}

func (s *calendarService) toVO(ctx context.Context, event *model.CalendarEvent) dtoresp.CalendarEventVO {
	vo := dtoresp.CalendarEventVO{
		ID:             event.ID,
		EventType:      string(event.EventType),
		RefID:          event.RefID,
		Title:          event.Title,
		StartTime:      event.StartTime.Format(time.RFC3339),
		Location:       event.Location,
		CustomNote:     event.CustomNote,
		RemindBefore:   []string(event.RemindBefore),
		ReminderStatus: "pending",
		CreatedAt:      event.CreatedAt.Format(time.RFC3339),
	}
	if vo.RemindBefore == nil {
		vo.RemindBefore = []string{}
	}
	if event.EndTime != nil {
		vo.EndTime = event.EndTime.Format(time.RFC3339)
	}
	vo.SourceUpdated = s.checkSourceUpdated(ctx, event)
	return vo
}

func (s *calendarService) checkSourceUpdated(ctx context.Context, event *model.CalendarEvent) bool {
	switch event.EventType {
	case model.EventTypeCareerTalk:
		talk, err := s.careerTalk.FindPublishedByID(ctx, event.RefID)
		if err != nil {
			return false
		}
		return talk.Title != event.Title ||
			!talk.StartTime.Equal(event.StartTime) ||
			talk.Location != event.Location
	case model.EventTypeJobFair:
		fair, err := s.jobFair.FindPublishedByID(ctx, event.RefID)
		if err != nil {
			return false
		}
		startTime := fair.StartDate
		if fair.StartTime != nil {
			startTime = *fair.StartTime
		}
		return fair.Title != event.Title ||
			!startTime.Equal(event.StartTime) ||
			fair.Location != event.Location
	}
	return false
}
