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
	list, total, page, pageSize, err := h.svc.List(c.Request.Context(), userID, c)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.Page(c, list, total, page, pageSize)
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
	data, err := h.svc.TriggerSync(c.Request.Context(), &req)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *AdminHandler) Placeholder(c *gin.Context) {
	response.Fail(c, 501, 50100, "管理端接口框架已就绪，业务逻辑待实现")
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
