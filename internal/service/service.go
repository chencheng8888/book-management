package service

import (
	"book-management/internal/pkg/tool"
	"github.com/google/wire"
	"time"
)

var ProviderSet = wire.NewSet()

var (
	WeekPattern  = "week_pattern"
	MonthPattern = "month_pattern"
	YearPattern  = "year_pattern"
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
