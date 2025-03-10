package tool

import (
	"book-management/internal/pkg/common"
	"time"
)

func CheckCategory(category string) bool {
	switch category {
	case common.ChildrenStory, common.ScienceKnowledge, common.ArtEnlightenment:
		return true
	default:
		return false
	}
}

const (
	Format1 = "2006-01-02 15:04:05"
	Format2 = "2006-01-02"
)

func IsTimeFormatValid(timeStr, format string) bool {
	t, err := time.Parse(format, timeStr)
	if err != nil {
		return false
	}
	// 将解析得到的时间重新格式化为目标格式，并与原始字符串比较
	return t.Format(format) == timeStr
}

// Intersection 返回两个切片的交集（元素唯一，按第一个切片的出现顺序）
func Intersection[T comparable](a, b []T) []T {
	// 创建b元素的快速查找集合
	setB := make(map[T]struct{})
	for _, v := range b {
		setB[v] = struct{}{}
	}

	var result []T
	seen := make(map[T]struct{}) // 用于记录已添加的元素

	// 遍历第一个切片
	for _, v := range a {
		// 检查元素是否同时存在于两个切片且未被记录
		if _, inB := setB[v]; inB {
			if _, exists := seen[v]; !exists {
				result = append(result, v)
				seen[v] = struct{}{}
			}
		}
	}
	return result
}

func GetShanghaiTime() time.Time {
	local, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(local)
}

func GetWeekStartTime(endTime time.Time) (startTime time.Time) {
	weekDay := endTime.Weekday()
	if weekDay == time.Sunday {
		return endTime.AddDate(0, 0, -6)
	} else {
		return endTime.AddDate(0, 0, -int(weekDay)+1)
	}
}
func GetMonthStartTime(endTime time.Time) (startTime time.Time) {
	return time.Date(endTime.Year(), endTime.Month(), 1, startTime.Hour(), startTime.Minute(), startTime.Second(), startTime.Nanosecond(), endTime.Location())
}

func GetYearStartTime(endTime time.Time) (startTime time.Time) {
	return time.Date(endTime.Year(), 1, 1, startTime.Hour(), startTime.Minute(), startTime.Second(), startTime.Nanosecond(), endTime.Location())
}
