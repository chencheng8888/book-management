package service

import (
	"book-management/internal/pkg/tool"
	"time"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewBookSvc, NewActivitySvc, NewVolunteerSvc,
	NewUserSvc)

var (
	WeekPattern  = "week"
	MonthPattern = "month"
	YearPattern  = "year"
)

func getStartAndEndTime(pattern string) (startTime, endTime time.Time) {
	currentTime := tool.GetShanghaiTime()
	switch pattern {
	case WeekPattern:
		startTime, endTime = tool.GetWeekStartTime(currentTime), currentTime
	case MonthPattern:
		startTime, endTime = tool.GetMonthStartTime(currentTime), currentTime
	case YearPattern:
		startTime, endTime = tool.GetYearStartTime(currentTime), currentTime
	default:
		startTime, endTime = tool.GetWeekStartTime(currentTime), currentTime
	}
	return startTime, endTime
}
