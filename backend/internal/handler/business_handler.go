package handler

import (
	"errors"

	"github.com/employment-center/campus-recruitment/internal/dto/request"
	"github.com/employment-center/campus-recruitment/internal/middleware"
	"github.com/employment-center/campus-recruitment/internal/service"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/employment-center/campus-recruitment/pkg/response"
	"github.com/gin-gonic/gin"
)

type CareerTalkHandler struct {
	svc service.CareerTalkService
}

func NewCareerTalkHandler(svc service.CareerTalkService) *CareerTalkHandler {
	return &CareerTalkHandler{svc: svc}
}

func (h *CareerTalkHandler) List(c *gin.Context) {
	var q request.CareerTalkQuery
	if !bindQuery(c, &q) {
		return
	}
	userID, _ := middleware.GetUserID(c)
	list, total, page, pageSize, err := h.svc.List(c.Request.Context(), &q, userID)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.Page(c, list, total, page, pageSize)
}

func (h *CareerTalkHandler) Detail(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	userID, _ := middleware.GetUserID(c)
	data, svcErr := h.svc.GetByID(c.Request.Context(), id, userID)
	if svcErr != nil {
		abortError(c, toAppError(svcErr))
		return
	}
	response.OK(c, data)
}

func (h *CareerTalkHandler) Upcoming(c *gin.Context) {
	data, err := h.svc.ListUpcomingWithin24h(c.Request.Context())
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *CareerTalkHandler) HotCompanies(c *gin.Context) {
	limit, err := parseUint(c.DefaultQuery("limit", "6"))
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	data, svcErr := h.svc.ListHotCompanies(c.Request.Context(), int(limit))
	if svcErr != nil {
		abortError(c, toAppError(svcErr))
		return
	}
	response.OK(c, data)
}

type JobFairHandler struct {
	svc service.JobFairService
}

func NewJobFairHandler(svc service.JobFairService) *JobFairHandler {
	return &JobFairHandler{svc: svc}
}

func (h *JobFairHandler) List(c *gin.Context) {
	var q request.JobFairQuery
	if !bindQuery(c, &q) {
		return
	}
	userID, _ := middleware.GetUserID(c)
	list, total, page, pageSize, err := h.svc.List(c.Request.Context(), &q, userID)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.Page(c, list, total, page, pageSize)
}

func (h *JobFairHandler) Detail(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	userID, _ := middleware.GetUserID(c)
	data, svcErr := h.svc.GetByID(c.Request.Context(), id, userID)
	if svcErr != nil {
		abortError(c, toAppError(svcErr))
		return
	}
	response.OK(c, data)
}

type RecommendationHandler struct {
	svc service.RecommendationService
}

func NewRecommendationHandler(svc service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{svc: svc}
}

func (h *RecommendationHandler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	result, err := h.svc.List(c.Request.Context(), userID, c)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, result)
}

func (h *RecommendationHandler) Dismiss(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	refID, err := parseUintParam(c, "refId")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	var req request.DismissRecommendationRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.svc.Dismiss(c.Request.Context(), userID, refID, req.EventType); err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, nil)
}

type CalendarHandler struct {
	svc service.CalendarService
}

func NewCalendarHandler(svc service.CalendarService) *CalendarHandler {
	return &CalendarHandler{svc: svc}
}

func (h *CalendarHandler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	list, err := h.svc.List(c.Request.Context(), userID, c)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, list)
}

func (h *CalendarHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	var req request.CreateCalendarEventRequest
	if !bindJSON(c, &req) {
		return
	}
	data, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *CalendarHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	var req request.UpdateCalendarEventRequest
	if !bindJSON(c, &req) {
		return
	}
	data, svcErr := h.svc.Update(c.Request.Context(), userID, id, &req)
	if svcErr != nil {
		abortError(c, toAppError(svcErr))
		return
	}
	response.OK(c, data)
}

func (h *CalendarHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), userID, id); err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, nil)
}

type ReminderHandler struct {
	svc service.ReminderService
}

func NewReminderHandler(svc service.ReminderService) *ReminderHandler {
	return &ReminderHandler{svc: svc}
}

func (h *ReminderHandler) ListLogs(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	list, total, page, pageSize, err := h.svc.ListLogs(c.Request.Context(), userID, c)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.Page(c, list, total, page, pageSize)
}

type AdminHandler struct {
	svc service.AdminService
}

func NewAdminHandler(svc service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) Sync(c *gin.Context) {
	var req request.SyncTriggerRequest
	_ = c.ShouldBindJSON(&req)
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	data, err := h.svc.TriggerSync(c.Request.Context(), operatorID, &req, c.ClientIP())
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *AdminHandler) ListSyncLogs(c *gin.Context) {
	var q request.SyncLogQuery
	if !bindQuery(c, &q) {
		return
	}
	list, total, page, pageSize, err := h.svc.ListSyncLogs(c.Request.Context(), &q)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.Page(c, list, total, page, pageSize)
}

func (h *AdminHandler) ListAuditLogs(c *gin.Context) {
	var q request.AuditLogQuery
	if !bindQuery(c, &q) {
		return
	}
	list, total, page, pageSize, err := h.svc.ListAuditLogs(c.Request.Context(), &q)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.Page(c, list, total, page, pageSize)
}

func (h *AdminHandler) ListCareerTalks(c *gin.Context) {
	var q request.AdminCareerTalkQuery
	if !bindQuery(c, &q) {
		return
	}
	list, total, page, pageSize, err := h.svc.ListCareerTalks(c.Request.Context(), &q)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.Page(c, list, total, page, pageSize)
}

func (h *AdminHandler) CreateCareerTalk(c *gin.Context) {
	var req request.AdminCareerTalkCreateRequest
	if !bindJSON(c, &req) {
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	data, err := h.svc.CreateCareerTalk(c.Request.Context(), operatorID, &req, c.ClientIP())
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *AdminHandler) UpdateCareerTalk(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	var req request.AdminCareerTalkUpdateRequest
	if !bindJSON(c, &req) {
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	data, svcErr := h.svc.UpdateCareerTalk(c.Request.Context(), operatorID, id, &req, c.ClientIP())
	if svcErr != nil {
		abortError(c, toAppError(svcErr))
		return
	}
	response.OK(c, data)
}

func (h *AdminHandler) DeleteCareerTalk(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	if err := h.svc.DeleteCareerTalk(c.Request.Context(), operatorID, id, c.ClientIP()); err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, nil)
}

func (h *AdminHandler) BatchCareerTalkStatus(c *gin.Context) {
	var req request.BatchPublishStatusRequest
	if !bindJSON(c, &req) {
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	if err := h.svc.BatchCareerTalkStatus(c.Request.Context(), operatorID, &req, c.ClientIP()); err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, nil)
}

func (h *AdminHandler) ListJobFairs(c *gin.Context) {
	var q request.AdminJobFairQuery
	if !bindQuery(c, &q) {
		return
	}
	list, total, page, pageSize, err := h.svc.ListJobFairs(c.Request.Context(), &q)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.Page(c, list, total, page, pageSize)
}

func (h *AdminHandler) CreateJobFair(c *gin.Context) {
	var req request.AdminJobFairCreateRequest
	if !bindJSON(c, &req) {
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	data, err := h.svc.CreateJobFair(c.Request.Context(), operatorID, &req, c.ClientIP())
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *AdminHandler) UpdateJobFair(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	var req request.AdminJobFairUpdateRequest
	if !bindJSON(c, &req) {
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	data, svcErr := h.svc.UpdateJobFair(c.Request.Context(), operatorID, id, &req, c.ClientIP())
	if svcErr != nil {
		abortError(c, toAppError(svcErr))
		return
	}
	response.OK(c, data)
}

func (h *AdminHandler) DeleteJobFair(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	if err := h.svc.DeleteJobFair(c.Request.Context(), operatorID, id, c.ClientIP()); err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, nil)
}

func (h *AdminHandler) BatchJobFairStatus(c *gin.Context) {
	var req request.BatchPublishStatusRequest
	if !bindJSON(c, &req) {
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	if err := h.svc.BatchJobFairStatus(c.Request.Context(), operatorID, &req, c.ClientIP()); err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, nil)
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	var q request.AdminUserQuery
	if !bindQuery(c, &q) {
		return
	}
	list, total, page, pageSize, err := h.svc.ListUsers(c.Request.Context(), &q)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.Page(c, list, total, page, pageSize)
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req request.AdminUserCreateRequest
	if !bindJSON(c, &req) {
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	data, err := h.svc.CreateUser(c.Request.Context(), operatorID, &req, c.ClientIP())
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	var req request.AdminUserUpdateRequest
	if !bindJSON(c, &req) {
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	data, svcErr := h.svc.UpdateUser(c.Request.Context(), operatorID, id, &req, c.ClientIP())
	if svcErr != nil {
		abortError(c, toAppError(svcErr))
		return
	}
	response.OK(c, data)
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	var req request.AdminUserStatusRequest
	if !bindJSON(c, &req) {
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	if err := h.svc.UpdateUserStatus(c.Request.Context(), operatorID, id, &req, c.ClientIP()); err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, nil)
}

func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		abortError(c, apperrors.ErrBadRequest)
		return
	}
	var req request.AdminResetPasswordRequest
	if !bindJSON(c, &req) {
		return
	}
	operatorID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	if err := h.svc.ResetUserPassword(c.Request.Context(), operatorID, id, &req, c.ClientIP()); err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, nil)
}

func toAppError(err error) *apperrors.AppError {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return apperrors.ErrInternalServer
}

func parseUintParam(c *gin.Context, name string) (uint64, error) {
	return parseUint(c.Param(name))
}

func parseUint(s string) (uint64, error) {
	var id uint64
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, apperrors.ErrBadRequest
		}
		id = id*10 + uint64(ch-'0')
	}
	if id == 0 && s != "0" {
		return 0, apperrors.ErrBadRequest
	}
	return id, nil
}
