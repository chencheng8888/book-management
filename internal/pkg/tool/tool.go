package tool

import (
	"book-management/internal/pkg/common"
	"fmt"
	"math"
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

func CheckBorrowStatus(status string) bool {
	switch status {
	case common.BookStatusWaitingReturn, common.BookStatusReturned, common.BookStatusOverdue:
		return true
	}
	return false
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

// // Intersection 返回两个切片的交集（元素唯一，按第一个切片的出现顺序）
// func Intersection[T comparable](a, b []T) []T {
// 	// 创建b元素的快速查找集合
// 	setB := make(map[T]struct{})
// 	for _, v := range b {
// 		setB[v] = struct{}{}
// 	}

// 	var result []T
// 	seen := make(map[T]struct{}) // 用于记录已添加的元素

// 	// 遍历第一个切片
// 	for _, v := range a {
// 		// 检查元素是否同时存在于两个切片且未被记录
// 		if _, inB := setB[v]; inB {
// 			if _, exists := seen[v]; !exists {
// 				result = append(result, v)
// 				seen[v] = struct{}{}
// 			}
// 		}
// 	}
// 	return result
// }

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

func GetPage(num int, pageSize int) int {
	return int(math.Ceil(float64(num) / float64(pageSize)))
}

func Unique[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	var list []T
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func ConvertTimeFormat(t time.Time, format string) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return t.In(loc).Format(format)
}

// ParseToShanghaiTime 将时间字符串按指定格式解析为上海时区的时间
// timeStr: 时间字符串
// format: 时间格式，如"2006-01-02 15:04:05"
func ParseToShanghaiTime(timeStr, format string) (time.Time, error) {
	// 加载上海时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, fmt.Errorf("加载时区失败: %w", err)
	}

	// 解析时间
	t, err := time.ParseInLocation(format, timeStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("时间解析失败: %w", err)
	}

	return t, nil
}
