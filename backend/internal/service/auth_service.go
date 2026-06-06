package service

import (
	"context"
	"regexp"
	"time"

	"github.com/employment-center/campus-recruitment/config"
	"github.com/employment-center/campus-recruitment/internal/dto/request"
	dtoresp "github.com/employment-center/campus-recruitment/internal/dto/response"
	"github.com/employment-center/campus-recruitment/internal/model"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"github.com/employment-center/campus-recruitment/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo repository.UserRepository
	jwt      *jwt.Manager
	authCfg  config.AuthConfig
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.Manager, authCfg config.AuthConfig) AuthService {
	return &authService{userRepo: userRepo, jwt: jwtManager, authCfg: authCfg}
}

var (
	hasLetter       = regexp.MustCompile(`[A-Za-z]`)
	hasDigit        = regexp.MustCompile(`\d`)
	usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_]{4,32}$`)
)

func validatePassword(password string) bool {
	return len(password) >= 8 && hasLetter.MatchString(password) && hasDigit.MatchString(password)
}

func validateUsername(username string) bool {
	return usernamePattern.MatchString(username)
}

func (s *authService) Register(ctx context.Context, req *request.RegisterRequest) error {
	if !validateUsername(req.Username) {
		return apperrors.ErrBadRequest
	}
	if !validatePassword(req.Password) {
		return apperrors.ErrWeakPassword
	}
	if exists, _ := s.userRepo.ExistsByUsername(ctx, req.Username); exists {
		return apperrors.ErrUsernameExists
	}
	if exists, _ := s.userRepo.ExistsByEmail(ctx, req.Email); exists {
		return apperrors.ErrEmailExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	user := &model.User{
		Username:        req.Username,
		PasswordHash:    string(hash),
		Email:           req.Email,
		TargetPositions: model.JSONStrings{},
		Role:            model.RoleStudent,
		Status:          model.StatusActive,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return apperrors.ErrInternalServer
	}
	return nil
}

func (s *authService) Login(ctx context.Context, req *request.LoginRequest) (*dtoresp.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}
	if user.Status == model.StatusLocked {
		if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
			return nil, apperrors.ErrAccountLocked
		}
		_ = s.userRepo.ResetLoginAttempts(ctx, user.ID)
		user.Status = model.StatusActive
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		_ = s.userRepo.IncrementLoginAttempts(ctx, user.ID, s.authCfg.MaxLoginAttempts, s.authCfg.LockMinutes)
		return nil, apperrors.ErrInvalidCredentials
	}
	_ = s.userRepo.ResetLoginAttempts(ctx, user.ID)
	token, expiresIn, err := s.jwt.Generate(user.ID, user.Username, string(user.Role), req.Remember)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	return &dtoresp.LoginResponse{Token: token, ExpiresIn: expiresIn, TokenType: "Bearer"}, nil
}
