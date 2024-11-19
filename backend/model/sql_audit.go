package model

import "time"

// SQLAudit SQL审计日志
type SQLAudit struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	DatabaseID  uint      `json:"database_id" gorm:"not null;comment:数据库ID"`
	UserID      uint      `json:"user_id" gorm:"not null;comment:用户ID"`
	Username    string    `json:"username" gorm:"size:50;not null;comment:用户名"`
	SQL         string    `json:"sql" gorm:"type:text;not null;comment:SQL语句"`
	Duration    int64     `json:"duration" gorm:"not null;comment:执行时长(毫秒)"`
	Status      string    `json:"status" gorm:"size:20;not null;comment:执行状态(success/failed)"`
	Error       string    `json:"error" gorm:"type:text;comment:错误信息"`
	AffectedRows int64    `json:"affected_rows" gorm:"comment:影响行数"`
	ClientIP    string    `json:"client_ip" gorm:"size:50;not null;comment:客户端IP"`
	CreatedAt   time.Time `json:"created_at"`
}

// SQLRisk SQL风险级别
type SQLRisk string

const (
	RiskLow    SQLRisk = "low"    // 低风险
	RiskMedium SQLRisk = "medium" // 中风险
	RiskHigh   SQLRisk = "high"   // 高风险
)

// SQLAnalysisResult SQL分析结果
type SQLAnalysisResult struct {
	Risk        SQLRisk `json:"risk"`         // 风险级别
	Description string  `json:"description"`   // 风险描述
	Suggestion  string  `json:"suggestion"`    // 优化建议
}
