package service

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/employment-center/campus-recruitment/internal/dto/request"
	dtoresp "github.com/employment-center/campus-recruitment/internal/dto/response"
	"github.com/employment-center/campus-recruitment/internal/model"
	"github.com/employment-center/campus-recruitment/internal/repository"
	"github.com/employment-center/campus-recruitment/pkg/pagination"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type recommendationService struct {
	userRepo       repository.UserRepository
	prefRepo       repository.PreferenceRepository
	careerTalkRepo repository.CareerTalkRepository
	jobFairRepo    repository.JobFairRepository
	calendarRepo   repository.CalendarRepository
}

func NewRecommendationService(
	userRepo repository.UserRepository,
	prefRepo repository.PreferenceRepository,
	careerTalkRepo repository.CareerTalkRepository,
	jobFairRepo repository.JobFairRepository,
	calendarRepo repository.CalendarRepository,
) RecommendationService {
	return &recommendationService{
		userRepo:       userRepo,
		prefRepo:       prefRepo,
		careerTalkRepo: careerTalkRepo,
		jobFairRepo:    jobFairRepo,
		calendarRepo:   calendarRepo,
	}
}

type matchContext struct {
	targetPositions    []string
	preferredCities    []string
	preferredCompanies []string
	focusCompanies     []string
	major              string
	college            string
}

type rankedRecommendation struct {
	vo       dtoresp.RecommendationVO
	score    int
	sortTime time.Time
}

func (s *recommendationService) List(ctx context.Context, userID uint64, c *gin.Context) (*dtoresp.RecommendationListResult, error) {
	var pq pagination.Query
	_ = c.ShouldBindQuery(&pq)
	page, pageSize := pq.Normalize()
	if pageSize == pagination.DefaultPageSize {
		pageSize = 20
	}

	ctxMatch, err := s.buildMatchContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	talks, err := s.listAllCareerTalks(ctx)
	if err != nil {
		return nil, err
	}
	fairs, err := s.listAllJobFairs(ctx)
	if err != nil {
		return nil, err
	}

	ranked := s.rankItems(talks, fairs, ctxMatch)
	fallback := len(ranked) == 0 || !s.hasMatched(ranked)

	if fallback {
		ranked = s.buildFallbackList(talks, fairs)
	} else {
		sort.SliceStable(ranked, func(i, j int) bool {
			if ranked[i].score != ranked[j].score {
				return ranked[i].score > ranked[j].score
			}
			return ranked[i].sortTime.Before(ranked[j].sortTime)
		})
	}

	total := int64(len(ranked))
	start := (page - 1) * pageSize
	if start > len(ranked) {
		start = len(ranked)
	}
	end := start + pageSize
	if end > len(ranked) {
		end = len(ranked)
	}

	list := make([]dtoresp.RecommendationVO, 0, end-start)
	for _, item := range ranked[start:end] {
		list = append(list, item.vo)
	}
	refSet := loadUserCalendarRefSet(ctx, s.calendarRepo, userID)
	for i := range list {
		list[i].InCalendar = isInUserCalendar(refSet, list[i].EventType, list[i].RefID)
	}

	return &dtoresp.RecommendationListResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Fallback: fallback,
	}, nil
}

func (s *recommendationService) Dismiss(_ context.Context, _, _ uint64, _ string) error {
	return nil
}

func (s *recommendationService) buildMatchContext(ctx context.Context, userID uint64) (*matchContext, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &matchContext{}, nil
		}
		return nil, err
	}

	mc := &matchContext{
		targetPositions: uniqueNonEmpty(user.TargetPositions),
		major:           strings.TrimSpace(user.Major),
		college:         strings.TrimSpace(user.College),
	}

	pref, err := s.prefRepo.FindByUserID(ctx, userID)
	if err == nil {
		mc.targetPositions = uniqueNonEmpty(append(mc.targetPositions, pref.TargetPositions...))
		mc.preferredCities = uniqueNonEmpty(pref.PreferredCities)
		mc.preferredCompanies = uniqueNonEmpty(pref.PreferredCompanies)
		mc.focusCompanies = uniqueNonEmpty(pref.FocusCompanies)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return mc, nil
}

func (s *recommendationService) listAllCareerTalks(ctx context.Context) ([]model.CareerTalk, error) {
	list, _, err := s.careerTalkRepo.ListPublished(ctx, &request.CareerTalkQuery{}, 1, 200)
	return list, err
}

func (s *recommendationService) listAllJobFairs(ctx context.Context) ([]model.JobFair, error) {
	list, _, err := s.jobFairRepo.ListPublished(ctx, &request.JobFairQuery{}, 1, 200)
	return list, err
}

func (s *recommendationService) rankItems(
	talks []model.CareerTalk,
	fairs []model.JobFair,
	mc *matchContext,
) []rankedRecommendation {
	var ranked []rankedRecommendation
	now := time.Now()

	for i := range talks {
		talk := talks[i]
		score, reasons := scoreCareerTalk(&talk, mc)
		if score <= 0 {
			continue
		}
		vo := toCareerTalkVO(&talk, now, score, reasons)
		ranked = append(ranked, rankedRecommendation{vo: vo, score: score, sortTime: talk.StartTime})
	}

	for i := range fairs {
		fair := fairs[i]
		score, reasons := scoreJobFair(&fair, mc)
		if score <= 0 {
			continue
		}
		vo := toJobFairVO(&fair, score, reasons)
		ranked = append(ranked, rankedRecommendation{vo: vo, score: score, sortTime: fair.StartDate})
	}

	return ranked
}

func (s *recommendationService) hasMatched(ranked []rankedRecommendation) bool {
	return len(ranked) > 0
}

func (s *recommendationService) buildFallbackList(
	talks []model.CareerTalk,
	fairs []model.JobFair,
) []rankedRecommendation {
	now := time.Now()
	var ranked []rankedRecommendation

	for i := range talks {
		talk := talks[i]
		vo := toCareerTalkVO(&talk, now, 0, nil)
		ranked = append(ranked, rankedRecommendation{vo: vo, score: 0, sortTime: talk.StartTime})
	}
	for i := range fairs {
		fair := fairs[i]
		vo := toJobFairVO(&fair, 0, nil)
		ranked = append(ranked, rankedRecommendation{vo: vo, score: 0, sortTime: fair.StartDate})
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		return ranked[i].sortTime.Before(ranked[j].sortTime)
	})
	return ranked
}

func scoreCareerTalk(talk *model.CareerTalk, mc *matchContext) (int, []string) {
	var score int
	var reasons []string

	for _, pos := range mc.targetPositions {
		for _, talkPos := range talk.Positions {
			if textMatch(talkPos, pos) {
				score += 20
				reasons = append(reasons, "匹配您的意向岗位："+pos)
				break
			}
		}
	}

	for _, company := range mc.focusCompanies {
		if textMatch(talk.Company, company) {
			score += 30
			reasons = append(reasons, "您特别关注的公司："+company)
			break
		}
	}

	for _, company := range mc.preferredCompanies {
		if textMatch(talk.Company, company) {
			score += 15
			reasons = append(reasons, "匹配偏好公司："+company)
			break
		}
	}

	for _, city := range mc.preferredCities {
		if textMatch(talk.Location, city) {
			score += 10
			reasons = append(reasons, "匹配意向城市："+city)
			break
		}
	}

	if mc.major != "" {
		for _, major := range talk.TargetMajors {
			if textMatch(major, mc.major) {
				score += 25
				reasons = append(reasons, "匹配您的专业："+mc.major)
				break
			}
		}
	}

	if mc.college != "" && (textMatch(talk.Title, mc.college) || textMatch(talk.Location, mc.college)) {
		score += 5
		reasons = append(reasons, "与您的学院相关："+mc.college)
	}

	return score, uniqueNonEmpty(reasons)
}

func scoreJobFair(fair *model.JobFair, mc *matchContext) (int, []string) {
	var score int
	var reasons []string

	if mc.major != "" {
		for _, major := range fair.TargetMajors {
			if textMatch(major, mc.major) {
				score += 25
				reasons = append(reasons, "匹配您的专业："+mc.major)
				break
			}
		}
	}

	for _, city := range mc.preferredCities {
		if textMatch(fair.Location, city) {
			score += 10
			reasons = append(reasons, "匹配意向城市："+city)
			break
		}
	}

	if mc.college != "" && textMatch(fair.TargetAudience, mc.college) {
		score += 15
		reasons = append(reasons, "面向您的学院："+mc.college)
	}

	for _, pos := range mc.targetPositions {
		if textMatch(fair.TargetAudience, pos) || textMatch(fair.Title, pos) {
			score += 15
			reasons = append(reasons, "匹配您的意向岗位："+pos)
			break
		}
	}

	for _, company := range append(mc.focusCompanies, mc.preferredCompanies...) {
		if textMatch(fair.Title, company) || textMatch(fair.TargetAudience, company) {
			score += 10
			reasons = append(reasons, "匹配偏好公司："+company)
			break
		}
	}

	return score, uniqueNonEmpty(reasons)
}

func toCareerTalkVO(talk *model.CareerTalk, now time.Time, score int, reasons []string) dtoresp.RecommendationVO {
	status := string(model.EventUpcoming)
	if talk.StartTime.Before(now) {
		status = string(model.EventEnded)
	}
	if reasons == nil {
		reasons = []string{}
	}
	return dtoresp.RecommendationVO{
		ID:           talk.ID,
		EventType:    "career_talk",
		RefID:        talk.ID,
		Title:        talk.Title,
		Company:      talk.Company,
		StartTime:    talk.StartTime.Format(time.RFC3339),
		Location:     talk.Location,
		IndustryCode: talk.IndustryCode,
		Industry:     industryLabel(talk.IndustryCode),
		CompanySize:  talk.CompanySize,
		Format:       string(talk.Format),
		Positions:    []string(talk.Positions),
		Status:       status,
		MatchScore:   score,
		MatchReasons: reasons,
	}
}

func toJobFairVO(fair *model.JobFair, score int, reasons []string) dtoresp.RecommendationVO {
	if reasons == nil {
		reasons = []string{}
	}
	vo := dtoresp.RecommendationVO{
		ID:             fair.ID,
		EventType:      "job_fair",
		RefID:          fair.ID,
		Title:          fair.Title,
		StartDate:      fair.StartDate.Format("2006-01-02"),
		StartTime:      fair.StartDate.Format(time.RFC3339),
		Location:       fair.Location,
		CompanyCount:   fair.CompanyCount,
		TargetAudience: fair.TargetAudience,
		DetailURL:      fair.DetailURL,
		MatchScore:     score,
		MatchReasons:   reasons,
	}
	if fair.EndDate != nil {
		vo.EndDate = fair.EndDate.Format("2006-01-02")
	}
	if fair.Deadline != nil {
		vo.Deadline = fair.Deadline.Format(time.RFC3339)
	}
	return vo
}

var industryLabels = map[string]string{
	"internet":       "互联网",
	"finance":        "金融",
	"manufacturing":  "制造",
	"consulting":     "咨询",
}

func industryLabel(code string) string {
	if label, ok := industryLabels[code]; ok {
		return label
	}
	return code
}

func textMatch(haystack, needle string) bool {
	h := strings.ToLower(strings.TrimSpace(haystack))
	n := strings.ToLower(strings.TrimSpace(needle))
	if h == "" || n == "" {
		return false
	}
	return strings.Contains(h, n) || strings.Contains(n, h)
}

func uniqueNonEmpty(items []string) []string {
	seen := make(map[string]bool, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}
