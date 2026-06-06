package service

import (
	"context"
	"errors"

	"github.com/employment-center/campus-recruitment/internal/dto/request"
	dtoresp "github.com/employment-center/campus-recruitment/internal/dto/response"
	"github.com/employment-center/campus-recruitment/internal/model"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/apperrors"
	"gorm.io/gorm"
)

type userService struct {
	userRepo repository.UserRepository
	prefRepo repository.PreferenceRepository
}

func NewUserService(userRepo repository.UserRepository, prefRepo repository.PreferenceRepository) UserService {
	return &userService{userRepo: userRepo, prefRepo: prefRepo}
}

func (s *userService) GetProfile(ctx context.Context, userID uint64) (*dtoresp.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	return toProfileResponse(user), nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint64, req *request.UpdateProfileRequest) (*dtoresp.UserProfileResponse, error) {
	if len(req.TargetPositions) == 0 {
		return nil, apperrors.ErrProfileIncomplete
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.ErrInternalServer
	}

	taken, err := s.userRepo.ExistsByEmailExcludeID(ctx, req.Email, userID)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	if taken {
		return nil, apperrors.ErrEmailExists
	}

	user.Name = req.Name
	user.College = req.College
	user.Major = req.Major
	user.Grade = req.Grade
	user.TargetPositions = model.JSONStrings(req.TargetPositions)
	user.Phone = req.Phone
	user.Email = req.Email

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperrors.ErrInternalServer
	}

	_ = s.syncPreferenceTargetPositions(ctx, userID, req.TargetPositions)

	return toProfileResponse(user), nil
}

func (s *userService) GetPreferences(ctx context.Context, userID uint64) (interface{}, error) {
	pref, err := s.prefRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, userErr := s.userRepo.FindByID(ctx, userID)
			if userErr != nil {
				return defaultPreferenceResponse(), nil
			}
			return toPreferenceResponseFromUser(user), nil
		}
		return nil, apperrors.ErrInternalServer
	}
	if len(pref.TargetPositions) == 0 {
		user, userErr := s.userRepo.FindByID(ctx, userID)
		if userErr == nil && len(user.TargetPositions) > 0 {
			pref.TargetPositions = user.TargetPositions
		}
	}
	return toPreferenceResponse(pref), nil
}

func (s *userService) UpdatePreferences(ctx context.Context, userID uint64, req *request.UpdatePreferenceRequest) (interface{}, error) {
	pref := &model.UserPreference{
		UserID:             userID,
		TargetPositions:    model.JSONStrings(req.TargetPositions),
		PreferredCities:    model.JSONStrings(req.PreferredCities),
		PreferredCompanies: model.JSONStrings(req.PreferredCompanies),
		FocusCompanies:     model.JSONStrings(req.FocusCompanies),
		RemindBefore:       model.JSONStrings(req.RemindBefore),
	}
	if err := s.prefRepo.Upsert(ctx, pref); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	_ = s.syncUserTargetPositions(ctx, userID, req.TargetPositions)
	return toPreferenceResponse(pref), nil
}

func (s *userService) syncUserTargetPositions(ctx context.Context, userID uint64, positions []string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	user.TargetPositions = model.JSONStrings(positions)
	return s.userRepo.Update(ctx, user)
}

func (s *userService) syncPreferenceTargetPositions(ctx context.Context, userID uint64, positions []string) error {
	pref, err := s.prefRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.prefRepo.Upsert(ctx, &model.UserPreference{
				UserID:          userID,
				TargetPositions: model.JSONStrings(positions),
				RemindBefore:    model.JSONStrings{"1d"},
			})
		}
		return err
	}
	pref.TargetPositions = model.JSONStrings(positions)
	return s.prefRepo.Upsert(ctx, pref)
}

func toProfileResponse(user *model.User) *dtoresp.UserProfileResponse {
	positions := []string(user.TargetPositions)
	if positions == nil {
		positions = []string{}
	}
	return &dtoresp.UserProfileResponse{
		ID:               user.ID,
		Username:         user.Username,
		Role:             string(user.Role),
		Status:           string(user.Status),
		Name:             user.Name,
		Email:            user.Email,
		College:          user.College,
		Major:            user.Major,
		Grade:            user.Grade,
		TargetPositions:  positions,
		Phone:            user.Phone,
		Avatar:           user.Avatar,
		ProfileCompleted: user.Name != "" && user.College != "" && user.Major != "" && len(positions) > 0,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}
}

func toPreferenceResponse(pref *model.UserPreference) *dtoresp.UserPreferenceResponse {
	return &dtoresp.UserPreferenceResponse{
		TargetPositions:    stringsOrEmpty(pref.TargetPositions),
		PreferredCities:    stringsOrEmpty(pref.PreferredCities),
		PreferredCompanies: stringsOrEmpty(pref.PreferredCompanies),
		FocusCompanies:     stringsOrEmpty(pref.FocusCompanies),
		RemindBefore:       stringsOrEmpty(pref.RemindBefore),
		UpdatedAt:          &pref.UpdatedAt,
	}
}

func toPreferenceResponseFromUser(user *model.User) *dtoresp.UserPreferenceResponse {
	return &dtoresp.UserPreferenceResponse{
		TargetPositions:    stringsOrEmpty(user.TargetPositions),
		PreferredCities:    []string{},
		PreferredCompanies: []string{},
		FocusCompanies:     []string{},
		RemindBefore:       []string{"1d"},
	}
}

func defaultPreferenceResponse() *dtoresp.UserPreferenceResponse {
	return &dtoresp.UserPreferenceResponse{
		TargetPositions:    []string{},
		PreferredCities:    []string{},
		PreferredCompanies: []string{},
		FocusCompanies:     []string{},
		RemindBefore:       []string{"1d"},
	}
}

func stringsOrEmpty(v model.JSONStrings) []string {
	if v == nil {
		return []string{}
	}
	return []string(v)
}
