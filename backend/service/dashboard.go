package service

import (
	"fmt"
	"time"

	"tools-admin/backend/model"
	"tools-admin/backend/pkg/db"
	"tools-admin/backend/pkg/log"
)

type DashboardService struct{}

// GetOverview 获取仪表盘概览数据
func (s *DashboardService) GetOverview() (*model.DashboardOverview, int) {
	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	log.Info(fmt.Sprintf("正在获取仪表盘数据，今日：%s，昨日：%s", today, yesterday))

	// 获取今日数据
	var todayTask model.TaskStatistics
	var todaySms model.SmsStatistics
	if err := db.Db.Where("DATE(date) = ?", today).First(&todayTask).Error; err != nil {
		log.Error(fmt.Sprintf("获取今日任务统计失败: %v", err))
		todayTask = model.TaskStatistics{} // 使用空对象避免空指针
	}
	if err := db.Db.Where("DATE(date) = ?", today).First(&todaySms).Error; err != nil {
		log.Error(fmt.Sprintf("获取今日短信统计失败: %v", err))
		todaySms = model.SmsStatistics{} // 使用空对象避免空指针
	}

	// 获取昨日数据
	var yesterdayTask model.TaskStatistics
	var yesterdaySms model.SmsStatistics
	if err := db.Db.Where("DATE(date) = ?", yesterday).First(&yesterdayTask).Error; err != nil {
		log.Error(fmt.Sprintf("获取昨日任务统计失败: %v", err))
		yesterdayTask = model.TaskStatistics{} // 使用空对象避免空指针
	}
	if err := db.Db.Where("DATE(date) = ?", yesterday).First(&yesterdaySms).Error; err != nil {
		log.Error(fmt.Sprintf("获取昨日短信统计失败: %v", err))
		yesterdaySms = model.SmsStatistics{} // 使用空对象避免空指针
	}

	// 计算趋势
	var smsTrend float64
	if yesterdaySms.TotalCount > 0 {
		smsTrend = float64(todaySms.TotalCount-yesterdaySms.TotalCount) / float64(yesterdaySms.TotalCount) * 100
	}

	// 计算今日任务成功率
	var successRate float64
	if todayTask.TotalCount > 0 {
		successRate = float64(todayTask.SuccessCount) / float64(todayTask.TotalCount) * 100
	}

	// 计算昨日任务成功率和趋势
	var yesterdaySuccessRate float64
	if yesterdayTask.TotalCount > 0 {
		yesterdaySuccessRate = float64(yesterdayTask.SuccessCount) / float64(yesterdayTask.TotalCount) * 100
	}
	var successRateTrend float64
	if yesterdaySuccessRate > 0 {
		successRateTrend = successRate - yesterdaySuccessRate
	}

	// 计算任务趋势
	var taskTrend float64
	if yesterdayTask.TotalCount > 0 {
		taskTrend = float64(todayTask.TotalCount-yesterdayTask.TotalCount) / float64(yesterdayTask.TotalCount) * 100
	}

	overview := &model.DashboardOverview{
		SmsCount:         todaySms.TotalCount,
		SmsTrend:        smsTrend,
		TaskCount:        todayTask.TotalCount,
		TaskTrend:        taskTrend,
		SuccessRate:      successRate,
		SuccessRateTrend: successRateTrend,
	}

	log.Info(fmt.Sprintf("获取到的仪表盘数据: %+v", overview))
	return overview, 0
}

// GetTaskChart 获取任务图表数据
func (s *DashboardService) GetTaskChart(period string) ([]model.ChartData, int) {
	days := 7
	if period == "month" {
		days = 30
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days+1)

	log.Info(fmt.Sprintf("正在获取任务图表数据，周期：%s，开始日期：%s，结束日期：%s", 
		period, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))

	var statistics []model.TaskStatistics
	if err := db.Db.Where("DATE(date) >= ? AND DATE(date) <= ?", 
		startDate.Format("2006-01-02"), 
		endDate.Format("2006-01-02")).
		Order("date").
		Find(&statistics).Error; err != nil {
		log.Error(fmt.Sprintf("获取任务统计数据失败: %v", err))
		return nil, 1
	}

	log.Info(fmt.Sprintf("获取到 %d 条任务统计数据", len(statistics)))

	chartData := make([]model.ChartData, 0)
	for _, stat := range statistics {
		chartData = append(chartData, model.ChartData{
			Date:    stat.Date.Format("2006-01-02"),
			Success: stat.SuccessCount,
			Fail:    stat.FailCount,
		})
	}

	return chartData, 0
}

// GetSmsChart 获取短信图表数据
func (s *DashboardService) GetSmsChart(period string) ([]model.ChartData, int) {
	days := 7
	if period == "month" {
		days = 30
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days+1)

	log.Info(fmt.Sprintf("正在获取短信图表数据，周期：%s，开始日期：%s，结束日期：%s", 
		period, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))

	var statistics []model.SmsStatistics
	if err := db.Db.Where("DATE(date) >= ? AND DATE(date) <= ?", 
		startDate.Format("2006-01-02"), 
		endDate.Format("2006-01-02")).
		Order("date").
		Find(&statistics).Error; err != nil {
		log.Error(fmt.Sprintf("获取短信统计数据失败: %v", err))
		return nil, 1
	}

	log.Info(fmt.Sprintf("获取到 %d 条短信统计数据", len(statistics)))

	chartData := make([]model.ChartData, 0)
	for _, stat := range statistics {
		chartData = append(chartData, model.ChartData{
			Date:        stat.Date.Format("2006-01-02"),
			Success:     stat.SuccessCount,
			Fail:        stat.FailCount,
			Total:       stat.TotalCount,
			SuccessRate: stat.SuccessRate,
		})
	}

	return chartData, 0
}
