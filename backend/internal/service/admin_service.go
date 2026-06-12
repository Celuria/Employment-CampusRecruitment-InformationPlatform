package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/employment-center/campus-recruitment/internal/dto/request"
	dtoresp "github.com/employment-center/campus-recruitment/internal/dto/response"
	"github.com/employment-center/campus-recruitment/internal/model"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/employment-center/campus-recruitment/pkg/pagination"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type adminService struct {
	careerTalk repository.CareerTalkRepository
	jobFair    repository.JobFairRepository
	user       repository.UserRepository
	auditLog   repository.AuditLogRepository
	syncLog    repository.SyncLogRepository
}

func NewAdminService(
	careerTalk repository.CareerTalkRepository,
	jobFair repository.JobFairRepository,
	user repository.UserRepository,
	auditLog repository.AuditLogRepository,
	syncLog repository.SyncLogRepository,
) AdminService {
	return &adminService{
		careerTalk: careerTalk,
		jobFair:    jobFair,
		user:       user,
		auditLog:   auditLog,
		syncLog:    syncLog,
	}
}

func (s *adminService) ListCareerTalks(ctx context.Context, q *request.AdminCareerTalkQuery) (interface{}, int64, int, int, error) {
	pq := &pagination.Query{Page: q.Page, PageSize: q.PageSize}
	page, pageSize := pq.Normalize()
	list, total, err := s.careerTalk.ListAdmin(ctx, q, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	vos := make([]dtoresp.AdminCareerTalkVO, len(list))
	for i := range list {
		vos[i] = dtoresp.ToAdminCareerTalkVO(&list[i])
	}
	return vos, total, page, pageSize, nil
}

func (s *adminService) CreateCareerTalk(ctx context.Context, operatorID uint64, req *request.AdminCareerTalkCreateRequest, ip string) (interface{}, error) {
	startTime, err := parseDateTime(req.StartTime)
	if err != nil {
		return nil, apperrors.ErrBadRequest
	}
	talk := &model.CareerTalk{
		Title:           req.Title,
		Company:         req.Company,
		IndustryCode:    req.IndustryCode,
		CompanySize:     req.CompanySize,
		StartTime:       startTime,
		Campus:          req.Campus,
		Venue:           req.Venue,
		Location:        model.FormatEventLocation(req.Campus, req.Venue),
		Format:          model.EventFormat(req.Format),
		Positions:       model.JSONStrings(req.Positions),
		TargetMajors:    model.JSONStrings(req.TargetMajors),
		RegistrationURL: req.RegistrationURL,
		SourceURL:       req.SourceURL,
		LogoURL:         req.LogoURL,
		Description:     req.Description,
		PublishStatus:   model.PublishDraft,
		SourceType:      "manual",
		CreatedBy:       &operatorID,
		UpdatedBy:       &operatorID,
	}
	if req.EndTime != "" {
		endTime, err := parseDateTime(req.EndTime)
		if err != nil {
			return nil, apperrors.ErrBadRequest
		}
		talk.EndTime = &endTime
	}
	if req.PublishStatus != "" {
		talk.PublishStatus = model.PublishStatus(req.PublishStatus)
	}
	if err := s.careerTalk.Create(ctx, talk); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "CREATE", "career_talk", talk.ID, req, ip)
	vo := dtoresp.ToAdminCareerTalkVO(talk)
	return vo, nil
}

func (s *adminService) UpdateCareerTalk(ctx context.Context, operatorID, id uint64, req *request.AdminCareerTalkUpdateRequest, ip string) (interface{}, error) {
	talk, err := s.careerTalk.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCareerTalkNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	if req.Title != nil {
		talk.Title = *req.Title
	}
	if req.Company != nil {
		talk.Company = *req.Company
	}
	if req.IndustryCode != nil {
		talk.IndustryCode = *req.IndustryCode
	}
	if req.CompanySize != nil {
		talk.CompanySize = *req.CompanySize
	}
	if req.StartTime != nil {
		startTime, err := parseDateTime(*req.StartTime)
		if err != nil {
			return nil, apperrors.ErrBadRequest
		}
		talk.StartTime = startTime
	}
	if req.EndTime != nil {
		if *req.EndTime == "" {
			talk.EndTime = nil
		} else {
			endTime, err := parseDateTime(*req.EndTime)
			if err != nil {
				return nil, apperrors.ErrBadRequest
			}
			talk.EndTime = &endTime
		}
	}
	if req.Campus != nil {
		talk.Campus = *req.Campus
	}
	if req.Venue != nil {
		talk.Venue = *req.Venue
	}
	if req.Campus != nil || req.Venue != nil {
		talk.Location = model.FormatEventLocation(talk.Campus, talk.Venue)
	}
	if req.Format != nil {
		talk.Format = model.EventFormat(*req.Format)
	}
	if req.Positions != nil {
		talk.Positions = model.JSONStrings(req.Positions)
	}
	if req.TargetMajors != nil {
		talk.TargetMajors = model.JSONStrings(req.TargetMajors)
	}
	if req.RegistrationURL != nil {
		talk.RegistrationURL = *req.RegistrationURL
	}
	if req.SourceURL != nil {
		talk.SourceURL = *req.SourceURL
	}
	if req.LogoURL != nil {
		talk.LogoURL = *req.LogoURL
	}
	if req.Description != nil {
		talk.Description = *req.Description
	}
	if req.PublishStatus != nil {
		talk.PublishStatus = model.PublishStatus(*req.PublishStatus)
	}
	talk.UpdatedBy = &operatorID
	if err := s.careerTalk.Update(ctx, talk); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "UPDATE", "career_talk", id, req, ip)
	vo := dtoresp.ToAdminCareerTalkVO(talk)
	return vo, nil
}

func (s *adminService) DeleteCareerTalk(ctx context.Context, operatorID, id uint64, ip string) error {
	if err := s.careerTalk.SoftDelete(ctx, id, operatorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrCareerTalkNotFound
		}
		return apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "DELETE", "career_talk", id, nil, ip)
	return nil
}

func (s *adminService) BatchCareerTalkStatus(ctx context.Context, operatorID uint64, req *request.BatchPublishStatusRequest, ip string) error {
	if err := s.careerTalk.BatchUpdateStatus(ctx, req.IDs, model.PublishStatus(req.PublishStatus), operatorID); err != nil {
		return apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "UPDATE", "career_talk", 0, req, ip)
	return nil
}

func (s *adminService) ListJobFairs(ctx context.Context, q *request.AdminJobFairQuery) (interface{}, int64, int, int, error) {
	pq := &pagination.Query{Page: q.Page, PageSize: q.PageSize}
	page, pageSize := pq.Normalize()
	list, total, err := s.jobFair.ListAdmin(ctx, q, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	vos := make([]dtoresp.AdminJobFairVO, len(list))
	for i := range list {
		vos[i] = dtoresp.ToAdminJobFairVO(&list[i])
	}
	return vos, total, page, pageSize, nil
}

func (s *adminService) CreateJobFair(ctx context.Context, operatorID uint64, req *request.AdminJobFairCreateRequest, ip string) (interface{}, error) {
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, apperrors.ErrBadRequest
	}
	fair := &model.JobFair{
		Title:          req.Title,
		StartDate:      startDate,
		Campus:         req.Campus,
		Venue:          req.Venue,
		Location:       model.FormatEventLocation(req.Campus, req.Venue),
		CompanyCount:   req.CompanyCount,
		TargetAudience: req.TargetAudience,
		TargetMajors:   model.JSONStrings(req.TargetMajors),
		DetailURL:      req.DetailURL,
		SourceURL:      req.SourceURL,
		Description:    req.Description,
		PublishStatus:  model.PublishDraft,
		SourceType:     "manual",
		CreatedBy:      &operatorID,
		UpdatedBy:      &operatorID,
	}
	if req.EndDate != "" {
		endDate, err := parseDate(req.EndDate)
		if err != nil {
			return nil, apperrors.ErrBadRequest
		}
		fair.EndDate = &endDate
	}
	if req.StartTime != "" {
		startTime, err := parseDateTime(req.StartTime)
		if err != nil {
			return nil, apperrors.ErrBadRequest
		}
		fair.StartTime = &startTime
	}
	if req.Deadline != "" {
		deadline, err := parseDateTime(req.Deadline)
		if err != nil {
			return nil, apperrors.ErrBadRequest
		}
		fair.Deadline = &deadline
	}
	if req.PublishStatus != "" {
		fair.PublishStatus = model.PublishStatus(req.PublishStatus)
	}
	if err := s.jobFair.Create(ctx, fair); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "CREATE", "job_fair", fair.ID, req, ip)
	vo := dtoresp.ToAdminJobFairVO(fair)
	return vo, nil
}

func (s *adminService) UpdateJobFair(ctx context.Context, operatorID, id uint64, req *request.AdminJobFairUpdateRequest, ip string) (interface{}, error) {
	fair, err := s.jobFair.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrJobFairNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	if req.Title != nil {
		fair.Title = *req.Title
	}
	if req.StartDate != nil {
		startDate, err := parseDate(*req.StartDate)
		if err != nil {
			return nil, apperrors.ErrBadRequest
		}
		fair.StartDate = startDate
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			fair.EndDate = nil
		} else {
			endDate, err := parseDate(*req.EndDate)
			if err != nil {
				return nil, apperrors.ErrBadRequest
			}
			fair.EndDate = &endDate
		}
	}
	if req.StartTime != nil {
		if *req.StartTime == "" {
			fair.StartTime = nil
		} else {
			startTime, err := parseDateTime(*req.StartTime)
			if err != nil {
				return nil, apperrors.ErrBadRequest
			}
			fair.StartTime = &startTime
		}
	}
	if req.Campus != nil {
		fair.Campus = *req.Campus
	}
	if req.Venue != nil {
		fair.Venue = *req.Venue
	}
	if req.Campus != nil || req.Venue != nil {
		fair.Location = model.FormatEventLocation(fair.Campus, fair.Venue)
	}
	if req.CompanyCount != nil {
		fair.CompanyCount = req.CompanyCount
	}
	if req.TargetAudience != nil {
		fair.TargetAudience = *req.TargetAudience
	}
	if req.TargetMajors != nil {
		fair.TargetMajors = model.JSONStrings(req.TargetMajors)
	}
	if req.Deadline != nil {
		if *req.Deadline == "" {
			fair.Deadline = nil
		} else {
			deadline, err := parseDateTime(*req.Deadline)
			if err != nil {
				return nil, apperrors.ErrBadRequest
			}
			fair.Deadline = &deadline
		}
	}
	if req.DetailURL != nil {
		fair.DetailURL = *req.DetailURL
	}
	if req.SourceURL != nil {
		fair.SourceURL = *req.SourceURL
	}
	if req.Description != nil {
		fair.Description = *req.Description
	}
	if req.PublishStatus != nil {
		fair.PublishStatus = model.PublishStatus(*req.PublishStatus)
	}
	fair.UpdatedBy = &operatorID
	if err := s.jobFair.Update(ctx, fair); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "UPDATE", "job_fair", id, req, ip)
	vo := dtoresp.ToAdminJobFairVO(fair)
	return vo, nil
}

func (s *adminService) DeleteJobFair(ctx context.Context, operatorID, id uint64, ip string) error {
	if err := s.jobFair.SoftDelete(ctx, id, operatorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrJobFairNotFound
		}
		return apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "DELETE", "job_fair", id, nil, ip)
	return nil
}

func (s *adminService) BatchJobFairStatus(ctx context.Context, operatorID uint64, req *request.BatchPublishStatusRequest, ip string) error {
	if err := s.jobFair.BatchUpdateStatus(ctx, req.IDs, model.PublishStatus(req.PublishStatus), operatorID); err != nil {
		return apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "UPDATE", "job_fair", 0, req, ip)
	return nil
}

func (s *adminService) ListUsers(ctx context.Context, q *request.AdminUserQuery) (interface{}, int64, int, int, error) {
	pq := &pagination.Query{Page: q.Page, PageSize: q.PageSize}
	page, pageSize := pq.Normalize()
	list, total, err := s.user.ListAdmin(ctx, q, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	vos := make([]dtoresp.AdminUserVO, len(list))
	for i := range list {
		vos[i] = dtoresp.ToAdminUserVO(&list[i])
	}
	return vos, total, page, pageSize, nil
}

func (s *adminService) CreateUser(ctx context.Context, operatorID uint64, req *request.AdminUserCreateRequest, ip string) (interface{}, error) {
	exists, err := s.user.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	if exists {
		return nil, apperrors.ErrUsernameExists
	}
	exists, err = s.user.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	if exists {
		return nil, apperrors.ErrEmailExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	u := &model.User{
		Username:        req.Username,
		PasswordHash:    string(hash),
		Email:           req.Email,
		Name:            req.Name,
		TargetPositions: model.JSONStrings{},
		Role:            model.UserRole(req.Role),
		Status:          model.StatusActive,
	}
	if err := s.user.Create(ctx, u); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "CREATE", "user", u.ID, map[string]string{"username": req.Username, "role": req.Role}, ip)
	vo := dtoresp.ToAdminUserVO(u)
	return vo, nil
}

func (s *adminService) UpdateUserStatus(ctx context.Context, operatorID, id uint64, req *request.AdminUserStatusRequest, ip string) error {
	if operatorID == id && req.Status == string(model.StatusDisabled) {
		return apperrors.ErrBadRequest
	}
	u, err := s.user.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrUserNotFound
		}
		return apperrors.ErrInternalServer
	}
	u.Status = model.UserStatus(req.Status)
	if req.Status == string(model.StatusActive) {
		u.LoginAttempts = 0
		u.LockedUntil = nil
	}
	if err := s.user.Update(ctx, u); err != nil {
		return apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "UPDATE", "user", id, req, ip)
	return nil
}

func (s *adminService) UpdateUser(ctx context.Context, operatorID, id uint64, req *request.AdminUserUpdateRequest, ip string) (interface{}, error) {
	if operatorID == id && req.Role != nil && *req.Role != string(model.RoleAdmin) {
		return nil, apperrors.ErrBadRequest
	}
	u, err := s.user.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Email != nil {
		exists, err := s.user.ExistsByEmailExcludeID(ctx, *req.Email, id)
		if err != nil {
			return nil, apperrors.ErrInternalServer
		}
		if exists {
			return nil, apperrors.ErrEmailExists
		}
		u.Email = *req.Email
	}
	if req.College != nil {
		u.College = *req.College
	}
	if req.Major != nil {
		u.Major = *req.Major
	}
	if req.Role != nil {
		u.Role = model.UserRole(*req.Role)
	}
	if err := s.user.Update(ctx, u); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "UPDATE", "user", id, req, ip)
	vo := dtoresp.ToAdminUserVO(u)
	return vo, nil
}

func (s *adminService) ResetUserPassword(ctx context.Context, operatorID, id uint64, req *request.AdminResetPasswordRequest, ip string) error {
	u, err := s.user.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrUserNotFound
		}
		return apperrors.ErrInternalServer
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	u.PasswordHash = string(hash)
	if err := s.user.Update(ctx, u); err != nil {
		return apperrors.ErrInternalServer
	}
	s.recordAudit(ctx, operatorID, "UPDATE", "user", id, map[string]string{"action": "reset_password"}, ip)
	return nil
}

func (s *adminService) TriggerSync(ctx context.Context, operatorID uint64, req *request.SyncTriggerRequest, ip string) (interface{}, error) {
	sourceType := req.SourceType
	if sourceType == "" {
		sourceType = "all"
	}
	taskID := fmt.Sprintf("sync-%d", time.Now().UnixNano())
	now := time.Now()
	log := &model.SyncLog{
		TaskID:     taskID,
		SourceType: sourceType,
		Status:     "running",
		StartedAt:  now,
		OperatorID: operatorID,
	}
	if err := s.syncLog.Create(ctx, log); err != nil {
		return nil, apperrors.ErrInternalServer
	}

	var updated int64
	var syncErr error
	if sourceType == "all" || sourceType == "career_talk" {
		n, err := s.careerTalk.RefreshSyncedAt(ctx, sourceType)
		updated += n
		if err != nil {
			syncErr = err
		}
	}
	if sourceType == "all" || sourceType == "job_fair" {
		n, err := s.jobFair.RefreshSyncedAt(ctx, sourceType)
		updated += n
		if err != nil {
			syncErr = err
		}
	}

	finished := time.Now()
	log.FinishedAt = &finished
	log.UpdatedCount = int(updated)
	if syncErr != nil {
		log.Status = "failed"
		log.ErrorMessage = syncErr.Error()
		_ = s.syncLog.Update(ctx, log)
		return nil, apperrors.ErrInternalServer
	}
	log.Status = "success"
	_ = s.syncLog.Update(ctx, log)
	s.recordAudit(ctx, operatorID, "SYNC", "sync", log.ID, req, ip)

	return map[string]interface{}{
		"taskId":    taskID,
		"status":    log.Status,
		"startedAt": now.Format(time.RFC3339),
		"message":   fmt.Sprintf("同步完成，更新 %d 条记录", updated),
	}, nil
}

func (s *adminService) ListSyncLogs(ctx context.Context, q *request.SyncLogQuery) (interface{}, int64, int, int, error) {
	pq := &pagination.Query{Page: q.Page, PageSize: q.PageSize}
	page, pageSize := pq.Normalize()
	list, total, err := s.syncLog.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	vos := make([]dtoresp.SyncLogVO, len(list))
	for i := range list {
		vos[i] = dtoresp.ToSyncLogVO(&list[i])
	}
	return vos, total, page, pageSize, nil
}

func (s *adminService) ListAuditLogs(ctx context.Context, q *request.AuditLogQuery) (interface{}, int64, int, int, error) {
	pq := &pagination.Query{Page: q.Page, PageSize: q.PageSize}
	page, pageSize := pq.Normalize()
	list, total, err := s.auditLog.List(ctx, q, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	vos := make([]dtoresp.AuditLogVO, len(list))
	for i := range list {
		operatorName := ""
		if op, err := s.user.FindByID(ctx, list[i].OperatorID); err == nil {
			operatorName = op.Name
			if operatorName == "" {
				operatorName = op.Username
			}
		}
		vos[i] = dtoresp.ToAuditLogVO(&list[i], operatorName)
	}
	return vos, total, page, pageSize, nil
}

func (s *adminService) recordAudit(ctx context.Context, operatorID uint64, action, resourceType string, resourceID uint64, detail interface{}, ip string) {
	var detailStr string
	if detail != nil {
		if b, err := json.Marshal(detail); err == nil {
			detailStr = string(b)
		}
	}
	_ = s.auditLog.Create(ctx, &model.AuditLog{
		OperatorID:   operatorID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Detail:       detailStr,
		IP:           ip,
	})
}

func parseDateTime(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid datetime: %s", s)
}

func parseDate(s string) (time.Time, error) {
	layouts := []string{"2006-01-02", "2006-01-02T15:04:05", time.RFC3339}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local), nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date: %s", s)
}
