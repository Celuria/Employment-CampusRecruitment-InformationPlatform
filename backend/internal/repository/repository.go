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

// HotCompany 热门公司（按规模聚合）
type HotCompany struct {
	Company     string `json:"company" gorm:"column:company"`
	CompanySize string `json:"companySize" gorm:"column:company_size"`
}

// CareerTalkRepository 宣讲会数据访问
type CareerTalkRepository interface {
	ListPublished(ctx context.Context, q *request.CareerTalkQuery, page, pageSize int) ([]model.CareerTalk, int64, error)
	FindPublishedByID(ctx context.Context, id uint64) (*model.CareerTalk, error)
	ListUpcomingWithin24h(ctx context.Context, limit int) ([]model.CareerTalk, error)
	ListHotCompanies(ctx context.Context, limit int) ([]HotCompany, error)
}

type careerTalkRepository struct{ db *gorm.DB }

func NewCareerTalkRepository(db *gorm.DB) CareerTalkRepository { return &careerTalkRepository{db: db} }

func (r *careerTalkRepository) ListPublished(ctx context.Context, q *request.CareerTalkQuery, page, pageSize int) ([]model.CareerTalk, int64, error) {
	var list []model.CareerTalk
	var total int64
	query := r.db.WithContext(ctx).Model(&model.CareerTalk{}).Where("publish_status = ?", model.PublishPublished)
	query = applyCareerTalkFilters(query, q)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Order(careerTalkOrderClause(q.SortBy)).Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func applyCareerTalkFilters(query *gorm.DB, q *request.CareerTalkQuery) *gorm.DB {
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		query = query.Where("title LIKE ? OR company LIKE ?", like, like)
	}
	if q.Company != "" {
		like := "%" + q.Company + "%"
		query = query.Where("company LIKE ?", like)
	}
	if q.Campus != "" && q.Campus != "all" {
		query = query.Where("campus = ?", q.Campus)
	}
	if q.Industry != "" && q.Industry != "all" {
		query = query.Where("industry_code = ?", q.Industry)
	}
	if q.DateRange != "" && q.DateRange != "all" {
		now := time.Now()
		loc := now.Location()
		switch q.DateRange {
		case "today":
			start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
			query = query.Where("start_time >= ? AND start_time < ?", start, start.Add(24*time.Hour))
		case "tomorrow":
			start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Add(24 * time.Hour)
			query = query.Where("start_time >= ? AND start_time < ?", start, start.Add(24*time.Hour))
		case "this_week":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			monday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -(weekday - 1))
			query = query.Where("start_time >= ? AND start_time < ?", monday, monday.AddDate(0, 0, 7))
		case "next_week":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			monday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -(weekday - 1))
			nextMonday := monday.AddDate(0, 0, 7)
			query = query.Where("start_time >= ? AND start_time < ?", nextMonday, nextMonday.AddDate(0, 0, 7))
		}
	}
	return query
}

func careerTalkOrderClause(sortBy string) string {
	switch sortBy {
	case "popularity":
		return "created_at DESC, start_time ASC"
	case "company_size":
		return companySizeRankSQL() + " DESC, start_time ASC"
	case "time_desc":
		return "start_time DESC"
	default:
		return "start_time ASC"
	}
}

func companySizeRankSQL() string {
	return "FIELD(company_size, '10000人以上', '1000-9999人', '500-999人', '150-500人', '50-150人', '50人以下')"
}

func (r *careerTalkRepository) ListUpcomingWithin24h(ctx context.Context, limit int) ([]model.CareerTalk, error) {
	var list []model.CareerTalk
	now := time.Now()
	deadline := now.Add(24 * time.Hour)
	err := r.db.WithContext(ctx).
		Where("publish_status = ?", model.PublishPublished).
		Where("start_time >= ? AND start_time < ?", now, deadline).
		Order("start_time ASC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

func (r *careerTalkRepository) ListHotCompanies(ctx context.Context, limit int) ([]HotCompany, error) {
	var rows []HotCompany
	rankSQL := companySizeRankSQL()
	err := r.db.WithContext(ctx).Raw(`
		SELECT company,
		       SUBSTRING_INDEX(GROUP_CONCAT(company_size ORDER BY `+rankSQL+` DESC SEPARATOR ','), ',', 1) AS company_size
		FROM career_talks
		WHERE publish_status = ?
		GROUP BY company
		ORDER BY MAX(`+rankSQL+`) DESC
		LIMIT ?
	`, model.PublishPublished, limit).Scan(&rows).Error
	return rows, err
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
	query = applyJobFairFilters(query, q)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Order(jobFairOrderClause(q.SortBy)).Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func applyJobFairFilters(query *gorm.DB, q *request.JobFairQuery) *gorm.DB {
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		query = query.Where("title LIKE ? OR location LIKE ?", like, like)
	}
	if q.Campus != "" && q.Campus != "all" {
		query = query.Where("campus = ?", q.Campus)
	}
	if q.StartDate != "" {
		query = query.Where("COALESCE(end_date, start_date) >= ?", q.StartDate)
	}
	if q.EndDate != "" {
		query = query.Where("start_date <= ?", q.EndDate)
	}
	return query
}

func jobFairOrderClause(sortBy string) string {
	switch sortBy {
	case "start_date_desc":
		return "start_date DESC"
	case "deadline_asc":
		return "deadline ASC"
	case "company_count_desc":
		return "company_count DESC, start_date ASC"
	default:
		return "start_date ASC"
	}
}

func (r *jobFairRepository) FindPublishedByID(ctx context.Context, id uint64) (*model.JobFair, error) {
	var fair model.JobFair
	err := r.db.WithContext(ctx).Where("id = ? AND publish_status = ?", id, model.PublishPublished).First(&fair).Error
	return &fair, err
}

// CalendarRepository 日历数据访问
type CalendarRepository interface {
	ListByUser(ctx context.Context, userID uint64, startDate, endDate, eventType string) ([]model.CalendarEvent, error)
	FindByID(ctx context.Context, userID, id uint64) (*model.CalendarEvent, error)
	ExistsByUserEvent(ctx context.Context, userID, refID uint64, eventType string) (bool, error)
	Create(ctx context.Context, event *model.CalendarEvent) error
	Update(ctx context.Context, userID, id uint64, note string, remindBefore model.JSONStrings) error
	Delete(ctx context.Context, userID, id uint64) error
}

type calendarRepository struct{ db *gorm.DB }

func NewCalendarRepository(db *gorm.DB) CalendarRepository { return &calendarRepository{db: db} }

func (r *calendarRepository) ListByUser(ctx context.Context, userID uint64, startDate, endDate, eventType string) ([]model.CalendarEvent, error) {
	var list []model.CalendarEvent
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if startDate != "" {
		query = query.Where("COALESCE(end_time, start_time) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("start_time < ?", endDate+" 23:59:59")
	}
	err := query.Order("start_time ASC").Find(&list).Error
	return list, err
}

func (r *calendarRepository) FindByID(ctx context.Context, userID, id uint64) (*model.CalendarEvent, error) {
	var event model.CalendarEvent
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&event).Error
	return &event, err
}

func (r *calendarRepository) ExistsByUserEvent(ctx context.Context, userID, refID uint64, eventType string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CalendarEvent{}).
		Where("user_id = ? AND event_type = ? AND ref_id = ?", userID, eventType, refID).
		Count(&count).Error
	return count > 0, err
}

func (r *calendarRepository) Create(ctx context.Context, event *model.CalendarEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *calendarRepository) Update(ctx context.Context, userID, id uint64, note string, remindBefore model.JSONStrings) error {
	result := r.db.WithContext(ctx).Model(&model.CalendarEvent{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"custom_note":   note,
			"remind_before": remindBefore,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *calendarRepository) Delete(ctx context.Context, userID, id uint64) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.CalendarEvent{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 以下 Repository 为占位，待业务迭代实现
type RecommendationRepository interface{}
type ReminderRepository interface{}
type AdminRepository interface{}

func NewRecommendationRepository(_ *gorm.DB) RecommendationRepository { return struct{}{} }
func NewReminderRepository(_ *gorm.DB) ReminderRepository           { return struct{}{} }
func NewAdminRepository(_ *gorm.DB) AdminRepository                 { return struct{}{} }
