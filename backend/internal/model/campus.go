package model

import "strings"

// 校区枚举
const (
	CampusNanhu        = "nanhu"
	CampusMafangshan   = "mafangshan"
	CampusYujiato      = "yujiato"
	CampusOnline       = "online"
)

var campusLabels = map[string]string{
	CampusNanhu:      "南湖校区",
	CampusMafangshan: "马房山校区",
	CampusYujiato:    "余家头校区",
	CampusOnline:     "线上",
}

// CampusLabel 返回校区中文名
func CampusLabel(campus string) string {
	if label, ok := campusLabels[campus]; ok {
		return label
	}
	return campus
}

// FormatEventLocation 组合校区与楼栋为展示用完整地点
func FormatEventLocation(campus, venue string) string {
	venue = strings.TrimSpace(venue)
	if campus == CampusOnline {
		if venue != "" {
			return "线上 · " + venue
		}
		return "线上"
	}
	label := CampusLabel(campus)
	if label == "" || label == campus {
		if venue != "" {
			return venue
		}
		return campus
	}
	if venue != "" {
		return label + " · " + venue
	}
	return label
}
