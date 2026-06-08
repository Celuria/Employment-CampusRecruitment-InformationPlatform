package service

import (
	"context"
	"time"

	"github.com/employment-center/campus-recruitment/internal/model"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/employment-center/campus-recruitment/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type reminderServiceImpl struct {
	repo repository.ReminderRepository
}

func NewReminderService(repo repository.ReminderRepository) ReminderService {
	return &reminderServiceImpl{repo: repo}
}

// ListLogs 查询用户的提醒记录（分页）
func (s *reminderServiceImpl) ListLogs(ctx context.Context, userID uint64, c *gin.Context) (interface{}, int64, int, int, error) {
	page := pagination.ParsePage(c.Query("page"))
	pageSize := pagination.ParsePageSize(c.Query("pageSize"))
	list, total, err := s.repo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, apperrors.ErrInternalServer
	}
	return list, total, page, pageSize, nil
}

// GenerateReminders 为日历事件生成提醒记录
// 根据 remindBefore 数组，为每个提醒时间创建一条 pending 记录
func (s *reminderServiceImpl) GenerateReminders(ctx context.Context, event *model.CalendarEvent) error {
	// 先清理该事件已有的 pending 提醒，避免重复
	if err := s.repo.DeleteByCalendarEvent(ctx, event.ID); err != nil {
		return err
	}

	if len(event.RemindBefore) == 0 {
		return nil
	}

	logs := make([]model.ReminderLog, 0, len(event.RemindBefore))
	for _, rb := range event.RemindBefore {
		scheduledTime := calculateScheduledTime(event.StartTime, rb)
		// 只生成未来的提醒
		if scheduledTime.Before(time.Now()) {
			continue
		}
		logs = append(logs, model.ReminderLog{
			CalendarEventID: event.ID,
			UserID:          event.UserID,
			EventTitle:      event.Title,
			EventType:       event.EventType,
			RemindBefore:    rb,
			ScheduledTime:   scheduledTime,
			Status:          model.ReminderPending,
		})
	}

	return s.repo.BatchCreate(ctx, logs)
}

// ProcessPending 处理所有到时间的待发送提醒
// 由定时调度器调用，返回本次处理数量
func (s *reminderServiceImpl) ProcessPending(ctx context.Context) (int, error) {
	now := time.Now()
	list, err := s.repo.FindPendingByScheduledTime(ctx, now, 100)
	if err != nil {
		return 0, err
	}

	processed := 0
	for _, log := range list {
		if err := s.repo.MarkSent(ctx, log.ID, now); err != nil {
			// 标记失败，写入原因并继续
			_ = s.repo.MarkFailed(ctx, log.ID, err.Error())
			continue
		}
		processed++
	}
	return processed, nil
}

// CancelByCalendarEvent 取消指定日历事件的所有待发送提醒
func (s *reminderServiceImpl) CancelByCalendarEvent(ctx context.Context, calendarEventID uint64) error {
	return s.repo.DeleteByCalendarEvent(ctx, calendarEventID)
}

// calculateScheduledTime 根据 remindBefore 值计算提醒目标时间
func calculateScheduledTime(startTime time.Time, remindBefore string) time.Time {
	switch remindBefore {
	case "1h":
		return startTime.Add(-1 * time.Hour)
	case "1d":
		return startTime.Add(-24 * time.Hour)
	case "3d":
		return startTime.Add(-72 * time.Hour)
	default:
		return startTime.Add(-24 * time.Hour) // 默认提前一天
	}
}
