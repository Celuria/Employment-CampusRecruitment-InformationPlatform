package service

import (
	"context"

	"github.com/employment-center/campus-recruitment/internal/model"
	"github.com/employment-center/campus-recruitment/internal/repository"
)

func loadUserCalendarRefSet(ctx context.Context, calendarRepo repository.CalendarRepository, userID uint64) map[string]map[uint64]bool {
	if userID == 0 {
		return nil
	}
	set, err := calendarRepo.ListUserRefSet(ctx, userID)
	if err != nil {
		return nil
	}
	return set
}

func isInUserCalendar(set map[string]map[uint64]bool, eventType string, refID uint64) bool {
	if set == nil {
		return false
	}
	return set[eventType][refID]
}

func enrichCareerTalksInCalendar(list []model.CareerTalk, set map[string]map[uint64]bool) []model.CareerTalk {
	if set == nil {
		return list
	}
	for i := range list {
		list[i].InCalendar = isInUserCalendar(set, string(model.EventTypeCareerTalk), list[i].ID)
	}
	return list
}

func enrichCareerTalkInCalendar(talk *model.CareerTalk, set map[string]map[uint64]bool) {
	if talk == nil || set == nil {
		return
	}
	talk.InCalendar = isInUserCalendar(set, string(model.EventTypeCareerTalk), talk.ID)
}

func enrichJobFairsInCalendar(list []model.JobFair, set map[string]map[uint64]bool) []model.JobFair {
	if set == nil {
		return list
	}
	for i := range list {
		list[i].InCalendar = isInUserCalendar(set, string(model.EventTypeJobFair), list[i].ID)
	}
	return list
}

func enrichJobFairInCalendar(fair *model.JobFair, set map[string]map[uint64]bool) {
	if fair == nil || set == nil {
		return
	}
	fair.InCalendar = isInUserCalendar(set, string(model.EventTypeJobFair), fair.ID)
}
