package cronutil

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// GetNextRunTime 根据cron表达式计算下次执行时间
func GetNextRunTime(cronExpr string) (*time.Time, error) {
	return GetNextRunTimeFrom(cronExpr, time.Now())
}

// GetNextRunTimeFrom 从指定时间开始计算下次执行时间
func GetNextRunTimeFrom(cronExpr string, from time.Time) (*time.Time, error) {
	if cronExpr == "" {
		return nil, nil
	}

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression: %v", err)
	}

	next := schedule.Next(from)
	return &next, nil
}

// ValidateCronExpr 验证cron表达式是否有效
func ValidateCronExpr(cronExpr string) error {
	if cronExpr == "" {
		return nil
	}
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %v", err)
	}
	return nil
}

// CommonPatterns 常用的cron表达式模式（包含秒）
var CommonPatterns = map[string]string{
	"每秒":      "* * * * * *",
	"每分钟":     "0 * * * * *",
	"每小时":     "0 0 * * * *",
	"每天凌晨":    "0 0 0 * * *",
	"每周一凌晨":   "0 0 0 * * 1",
	"每月1号凌晨":  "0 0 0 1 * *",
	"工作日早上9点": "0 0 9 * * 1-5",
	"每天早上9点":  "0 0 9 * * *",
}
