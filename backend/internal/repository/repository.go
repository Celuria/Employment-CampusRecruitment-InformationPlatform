package repository

import (
	"context"
	"time"

	"github.com/employment-center/campus-recruitment/internal/dto/request"
	"github.com/employment-center/campus-recruitment/internal/model"
	"gorm.io/gorm"
)

// Repositories 聚合所有数据访问层
type Repositories struct {
	User           UserRepository
	Preference     PreferenceRepository
	CareerTalk     CareerTalkRepository
	JobFair        JobFairRepository
	Recommendation RecommendationRepository
	Calendar       CalendarRepository
	Reminder       ReminderRepository
	Admin          AdminRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:           NewUserRepository(db),
		Preference:     NewPreferenceRepository(db),
		CareerTalk:     NewCareerTalkRepository(db),
		JobFair:        NewJobFairRepository(db),
		Recommendation: NewRecommendationRepository(db),
		Calendar:       NewCalendarRepository(db),
		Reminder:       NewReminderRepository(db),
		Admin:          NewAdminRepository(db),
	}
}

// UserRepository 用户数据访问
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByEmailExcludeID(ctx context.Context, email string, excludeID uint64) (bool, error)
	Update(ctx context.Context, user *model.User) error
	IncrementLoginAttempts(ctx context.Context, id uint64, maxAttempts, lockMinutes int) error
	ResetLoginAttempts(ctx context.Context, id uint64) error
}

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepository{db: db} }

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) ExistsByEmailExcludeID(ctx context.Context, email string, excludeID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email = ? AND id <> ?", email, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) IncrementLoginAttempts(ctx context.Context, id uint64, maxAttempts, lockMinutes int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}
		user.LoginAttempts++
		if user.LoginAttempts >= maxAttempts {
			lockUntil := time.Now().Add(time.Duration(lockMinutes) * time.Minute)
			user.Status = model.StatusLocked
			user.LockedUntil = &lockUntil
		}
		return tx.Save(&user).Error
	})
}

func (r *userRepository) ResetLoginAttempts(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"login_attempts": 0,
		"status":         model.StatusActive,
		"locked_until":   nil,
		"last_login_at":  time.Now(),
	}).Error
}

// PreferenceRepository 偏好数据访问
type PreferenceRepository interface {
	FindByUserID(ctx context.Context, userID uint64) (*model.UserPreference, error)
	Upsert(ctx context.Context, pref *model.UserPreference) error
}

type preferenceRepository struct{ db *gorm.DB }

func NewPreferenceRepository(db *gorm.DB) PreferenceRepository { return &preferenceRepository{db: db} }

func (r *preferenceRepository) FindByUserID(ctx context.Context, userID uint64) (*model.UserPreference, error) {
	var pref model.UserPreference
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&pref).Error
	return &pref, err
}

func (r *preferenceRepository) Upsert(ctx context.Context, pref *model.UserPreference) error {
	var existing model.UserPreference
	err := r.db.WithContext(ctx).Where("user_id = ?", pref.UserID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(pref).Error
	}
	if err != nil {
		return err
	}
	// 仅更新业务字段，避免 Save 将 created_at 零值写入 MySQL
	return r.db.WithContext(ctx).Model(&existing).Updates(map[string]interface{}{
		"target_positions":    pref.TargetPositions,
		"preferred_cities":    pref.PreferredCities,
		"preferred_companies": pref.PreferredCompanies,
		"focus_companies":     pref.FocusCompanies,
		"remind_before":       pref.RemindBefore,
	}).Error
}

// CareerTalkRepository 宣讲会数据访问
type CareerTalkRepository interface {
	ListPublished(ctx context.Context, q *request.CareerTalkQuery, page, pageSize int) ([]model.CareerTalk, int64, error)
	FindPublishedByID(ctx context.Context, id uint64) (*model.CareerTalk, error)
}

type careerTalkRepository struct{ db *gorm.DB }

func NewCareerTalkRepository(db *gorm.DB) CareerTalkRepository { return &careerTalkRepository{db: db} }

func (r *careerTalkRepository) ListPublished(ctx context.Context, q *request.CareerTalkQuery, page, pageSize int) ([]model.CareerTalk, int64, error) {
	var list []model.CareerTalk
	var total int64
	query := r.db.WithContext(ctx).Model(&model.CareerTalk{}).Where("publish_status = ?", model.PublishPublished)
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		query = query.Where("title LIKE ? OR company LIKE ?", like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Order("start_time ASC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *careerTalkRepository) FindPublishedByID(ctx context.Context, id uint64) (*model.CareerTalk, error) {
	var talk model.CareerTalk
	err := r.db.WithContext(ctx).Where("id = ? AND publish_status = ?", id, model.PublishPublished).First(&talk).Error
	return &talk, err
}

// JobFairRepository 双选会数据访问
type JobFairRepository interface {
	ListPublished(ctx context.Context, q *request.JobFairQuery, page, pageSize int) ([]model.JobFair, int64, error)
	FindPublishedByID(ctx context.Context, id uint64) (*model.JobFair, error)
}

type jobFairRepository struct{ db *gorm.DB }

func NewJobFairRepository(db *gorm.DB) JobFairRepository { return &jobFairRepository{db: db} }

func (r *jobFairRepository) ListPublished(ctx context.Context, q *request.JobFairQuery, page, pageSize int) ([]model.JobFair, int64, error) {
	var list []model.JobFair
	var total int64
	query := r.db.WithContext(ctx).Model(&model.JobFair{}).Where("publish_status = ?", model.PublishPublished)
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		query = query.Where("title LIKE ?", like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Order("start_date ASC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *jobFairRepository) FindPublishedByID(ctx context.Context, id uint64) (*model.JobFair, error) {
	var fair model.JobFair
	err := r.db.WithContext(ctx).Where("id = ? AND publish_status = ?", id, model.PublishPublished).First(&fair).Error
	return &fair, err
}

// 以下 Repository 为占位，待业务迭代实现
type RecommendationRepository interface{}
type CalendarRepository interface{}
type ReminderRepository interface{}
type AdminRepository interface{}

func NewRecommendationRepository(_ *gorm.DB) RecommendationRepository { return struct{}{} }
func NewCalendarRepository(_ *gorm.DB) CalendarRepository           { return struct{}{} }
func NewReminderRepository(_ *gorm.DB) ReminderRepository           { return struct{}{} }
func NewAdminRepository(_ *gorm.DB) AdminRepository                 { return struct{}{} }
