package service

import (
	"context"
	"fmt"
	"time"

	"github.com/employment-center/campus-recruitment/internal/model"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/employment-center/campus-recruitment/pkg/email"
	"github.com/employment-center/campus-recruitment/pkg/logger"
	"github.com/employment-center/campus-recruitment/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type reminderServiceImpl struct {
	repo        repository.ReminderRepository
	userRepo    repository.UserRepository
	emailSender *email.Sender
}

func NewReminderService(
	repo repository.ReminderRepository,
	userRepo repository.UserRepository,
	emailSender *email.Sender,
) ReminderService {
	return &reminderServiceImpl{
		repo:        repo,
		userRepo:    userRepo,
		emailSender: emailSender,
	}
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
// 由定时调度器调用，标记已发送并发邮件通知用户
func (s *reminderServiceImpl) ProcessPending(ctx context.Context) (int, error) {
	now := time.Now()
	list, err := s.repo.FindPendingByScheduledTime(ctx, now, 100)
	if err != nil {
		return 0, err
	}

	processed := 0
	for _, rm := range list {
		if err := s.repo.MarkSent(ctx, rm.ID, now); err != nil {
			_ = s.repo.MarkFailed(ctx, rm.ID, err.Error())
			continue
		}
		processed++

		// 发送邮件提醒
		go s.sendReminderEmail(rm)
	}
	return processed, nil
}

// CancelByCalendarEvent 取消指定日历事件的所有待发送提醒
func (s *reminderServiceImpl) CancelByCalendarEvent(ctx context.Context, calendarEventID uint64) error {
	return s.repo.DeleteByCalendarEvent(ctx, calendarEventID)
}

// sendReminderEmail 向用户发送活动提醒邮件
func (s *reminderServiceImpl) sendReminderEmail(rm model.ReminderLog) {
	sugar := logger.Log.Sugar()

	if s.emailSender == nil || !s.emailSender.Enabled() {
		sugar.Warnf("reminder email: skipped (email disabled), event=%s user=%d", rm.EventTitle, rm.UserID)
		return
	}

	user, err := s.userRepo.FindByID(context.Background(), rm.UserID)
	if err != nil {
		sugar.Errorf("reminder email: find user %d failed: %v", rm.UserID, err)
		return
	}
	if user.Email == "" {
		sugar.Warnf("reminder email: user %d has no email, skipped", rm.UserID)
		return
	}

	eventLabel := "宣讲会"
	if rm.EventType == model.EventTypeJobFair {
		eventLabel = "双选会"
	}

	subject := fmt.Sprintf("【活动提醒】%s - %s", rm.EventTitle, eventLabel)
	body := fmt.Sprintf(`您好 %s，

您关注的%s「%s」即将开始。

活动时间：%s
活动地点：%s

请留意活动时间，提前做好准备。

此邮件由系统自动发送，请勿回复。
`,
		user.Name,
		eventLabel,
		rm.EventTitle,
		rm.ScheduledTime.Format("2006-01-02 15:04"),
		"-",
	)

	if err := s.emailSender.Send([]string{user.Email}, subject, body); err != nil {
		sugar.Errorf("reminder email: send to %s failed: %v", user.Email, err)
	} else {
		sugar.Infof("reminder email: sent to %s, event=%s", user.Email, rm.EventTitle)
	}
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
