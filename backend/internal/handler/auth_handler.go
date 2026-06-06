package handler

import (
	"github.com/employment-center/campus-recruitment/internal/dto/request"
	"github.com/employment-center/campus-recruitment/internal/service"
	"github.com/employment-center/campus-recruitment/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.svc.Register(c.Request.Context(), &req); err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if !bindJSON(c, &req) {
		return
	}
	data, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		abortError(c, toAppError(err))
		return
	}
	response.OK(c, data)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.OK(c, nil)
}

func (h *AuthHandler) Captcha(c *gin.Context) {
	response.OK(c, gin.H{"captchaKey": "", "captchaImage": ""})
}
