package controller

import (
	"book-management/internal/pkg/resp"
	"book-management/internal/repository/do"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type HomeController struct {
	db *gorm.DB
}

func NewHomeController(db *gorm.DB) *HomeController {
	return &HomeController{
		db: db,
	}
}

func (h *HomeController) RegisterRoute(r *gin.Engine) {
	g := r.Group("/api/v1/home")
	{
		g.GET("/get_statics", h.GetHomePage)
	}
}

// GetHomePage 获取首页统计数据
// @Summary 获取首页统计数据
// @Description 获取库存总数、今日借阅数量、本月借阅数量、活跃用户数量、热门书籍、平均借阅时长、新增用户数、库存不足数量以及逾期未归还书本数量
// @Tags 首页
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Success 200 {object} HomePageResp "统计数据"
// @Router /api/v1/home/get_statics [get]
func (h *HomeController) GetHomePage(c *gin.Context) {
	var (
		totalStock        int64
		todayBorrowed     int64
		monthBorrowed     int64
		activeUsers       int64
		hotBook           string
		avgBorrowDuration float64
		newUsers          int64
		lowStockCount     int64
		overdueBooks      int64
		monthlyBorrowed   []MonthlyBorrowed
	)

	// 获取库存总数
	h.db.Model(&do.BookStock{}).Select("SUM(stock)").Scan(&totalStock)

	// 今日借阅的书籍数量
	h.db.Model(&do.BookBorrow{}).
		Where("DATE(created_time) = ?", time.Now().Format("2006-01-02")).
		Count(&todayBorrowed)

	// 本月借阅的书籍数量
	h.db.Model(&do.BookBorrow{}).
		Where("MONTH(created_time) = ? AND YEAR(created_time) = ?", time.Now().Month(), time.Now().Year()).
		Count(&monthBorrowed)

	// 活跃用户数量
	h.db.Raw(`
		SELECT COUNT(DISTINCT borrower_id) 
		FROM (
			SELECT borrower_id, COUNT(*) AS borrow_count 
			FROM book_borrow 
			WHERE created_time >= ? AND return_time IS NOT NULL 
			GROUP BY borrower_id 
			HAVING borrow_count > 3
		) AS active_users
	`, time.Now().AddDate(0, 0, -7)).Scan(&activeUsers)

	var hotBookID uint64

	// 热门书籍（借阅最多的书籍）
	h.db.Raw(`
		SELECT book_id 
		FROM book_borrow 
		GROUP BY book_id 
		ORDER BY COUNT(*) DESC 
		LIMIT 1
	`).Scan(&hotBookID)

	if hotBookID != 0 {
		h.db.Model(&do.BookInfo{}).
			Select("name").
			Where("id = ?", hotBookID).
			Scan(&hotBook)
	}

	// 平均借阅时长（单位天）
	h.db.Raw(`
		SELECT AVG(DATEDIFF(return_time, created_time)) 
		FROM book_borrow 
		WHERE return_time IS NOT NULL
	`).Scan(&avgBorrowDuration)

	// 新增用户数（最近一周内创建的用户数目）
	h.db.Model(&User{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Count(&newUsers)

	// 库存不足的数量（库存少于20本的）
	h.db.Model(&do.BookStock{}).
		Where("stock < ?", 20).
		Count(&lowStockCount)

	// 逾期未归还的书本数量
	h.db.Model(&do.BookBorrow{}).
		Where("expected_return_time < ? AND return_time IS NULL", time.Now()).
		Count(&overdueBooks)

	// 近六个月每月借阅数目
	h.db.Raw(`
		SELECT DATE_FORMAT(created_time, '%Y-%m') AS month, COUNT(*) AS borrow_count
		FROM book_borrow
		WHERE created_time >= ?
		GROUP BY month
		ORDER BY month ASC
	`, time.Now().AddDate(0, -6, 0)).Scan(&monthlyBorrowed)

	// 返回结果
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"total_stock":         totalStock,
		"today_borrowed":      todayBorrowed,
		"month_borrowed":      monthBorrowed,
		"active_users":        activeUsers,
		"hot_book":            hotBook,
		"avg_borrow_duration": avgBorrowDuration,
		"new_users":           newUsers,
		"low_stock_count":     lowStockCount,
		"overdue_books":       overdueBooks,
		"monthly_borrowed":    monthlyBorrowed,
	}))
}

// HomePageResp 首页统计数据响应体
type HomePageResp struct {
	TotalStock        int64             `json:"total_stock"`         // 库存总数
	TodayBorrowed     int64             `json:"today_borrowed"`      // 今日借阅数量
	MonthBorrowed     int64             `json:"month_borrowed"`      // 本月借阅数量
	ActiveUsers       int64             `json:"active_users"`        // 活跃用户数量
	HotBook           string            `json:"hot_book"`            // 热门书籍
	AvgBorrowDuration float64           `json:"avg_borrow_duration"` // 平均借阅时长（单位天）
	NewUsers          int64             `json:"new_users"`           // 新增用户数
	LowStockCount     int64             `json:"low_stock_count"`     // 库存不足数量
	OverdueBooks      int64             `json:"overdue_books"`       // 逾期未归还书本数量
	MonthlyBorrowed   []MonthlyBorrowed `json:"monthly_borrowed"`    // 近六个月每月借阅数目
}

// MonthlyBorrowed 近六个月每月借阅数目
type MonthlyBorrowed struct {
	Month       string `json:"month"`        // 月份
	BorrowCount int64  `json:"borrow_count"` // 借阅数目
}
