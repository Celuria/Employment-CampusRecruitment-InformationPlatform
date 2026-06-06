package handler

import (
	"github.com/employment-center/campus-recruitment/internal/dto/request"
	"github.com/employment-center/campus-recruitment/internal/middleware"
	"github.com/employment-center/campus-recruitment/internal/service"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/employment-center/campus-recruitment/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	data, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	var req request.UpdateProfileRequest
	if !bindJSON(c, &req) {
		return
	}
	data, err := h.svc.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *UserHandler) GetPreferences(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	data, err := h.svc.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *UserHandler) UpdatePreferences(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		abortError(c, apperrors.ErrUnauthorized)
		return
	}
	var req request.UpdatePreferenceRequest
	if !bindJSON(c, &req) {
		return
	}
	data, err := h.svc.UpdatePreferences(c.Request.Context(), userID, &req)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}
