package model

import (
	"time"
)

// TaskStatistics 任务统计
type TaskStatistics struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	Date         time.Time `json:"date" gorm:"type:date;uniqueIndex"`
	TotalCount   int       `json:"total_count"`
	SuccessCount int       `json:"success_count"`
	FailCount    int       `json:"fail_count"`
	SuccessRate  float64   `json:"success_rate" gorm:"type:decimal(5,2)"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SmsStatistics 短信统计
type SmsStatistics struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	Date         time.Time `json:"date" gorm:"type:date;uniqueIndex"`
	TotalCount   int       `json:"total_count"`
	SuccessCount int       `json:"success_count"`
	FailCount    int       `json:"fail_count"`
	SuccessRate  float64   `json:"success_rate" gorm:"type:decimal(5,2)"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// DashboardOverview 仪表盘概览数据
type DashboardOverview struct {
	SmsCount         int     `json:"sms_count"`          // 今日短信数
	SmsTrend         float64 `json:"sms_trend"`          // 短信数同比趋势
	TaskCount        int     `json:"task_count"`         // 今日任务数
	TaskTrend        float64 `json:"task_trend"`         // 任务数同比趋势
	SuccessRate      float64 `json:"success_rate"`       // 今日执行成功率
	SuccessRateTrend float64 `json:"success_rate_trend"` // 成功率同比趋势
}

// ChartData 图表数据
type ChartData struct {
	Date        string  `json:"date"`
	Count       int     `json:"count,omitempty"`
	Total       int     `json:"total,omitempty"`
	SuccessRate float64 `json:"success_rate,omitempty"`
	Success     int     `json:"success,omitempty"`
	Fail        int     `json:"fail,omitempty"`
}
