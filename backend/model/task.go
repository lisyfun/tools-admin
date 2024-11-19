package model

import (
	"time"

	"gorm.io/gorm"
)

// TaskType 任务类型
type TaskType int8

// TaskStatus 任务状态
type TaskStatus int8

// TaskExecStatus 任务执行状态
type TaskExecStatus int8

const (
	TaskTypeShell     TaskType = iota + 1 // Shell脚本任务
	TaskTypeHttp                          // HTTP请求任务
	TaskTypeDatax                         // Datax数据同步任务
	TaskTypeRegular                       // 常规任务
	TaskTypeUrgent                        // 紧急任务
	TaskTypeLongTerm                      // 长期任务
	TaskTypeRecurring                     // 循环任务
)

const (
	TaskStatusStarted TaskStatus = 1 // 启动
	TaskStatusStopped TaskStatus = 2 // 停止
)

const (
	TaskExecStatusPending  TaskExecStatus = 1 // 待执行
	TaskExecStatusRunning  TaskExecStatus = 2 // 执行中
	TaskExecStatusSuccess  TaskExecStatus = 3 // 执行成功
	TaskExecStatusFailed   TaskExecStatus = 4 // 执行失败
)

// 任务类型映射
var TaskTypeMap = map[TaskType]string{
	TaskTypeShell:     "shell",
	TaskTypeHttp:      "http",
	TaskTypeDatax:     "datax",
	TaskTypeRegular:   "regular",
	TaskTypeUrgent:    "urgent",
	TaskTypeLongTerm:  "longterm",
	TaskTypeRecurring: "recurring",
}

// 任务状态映射
var TaskStatusMap = map[TaskStatus]string{
	TaskStatusStarted: "started",
	TaskStatusStopped: "stopped",
}

// 任务执行状态映射
var TaskExecStatusMap = map[TaskExecStatus]string{
	TaskExecStatusPending: "pending",
	TaskExecStatusRunning: "running",
	TaskExecStatusSuccess: "success",
	TaskExecStatusFailed:  "failed",
}

// Task 任务模型
type Task struct {
	gorm.Model
	Name        string        `json:"name" gorm:"type:varchar(100);not null"`
	Type        TaskType      `json:"type" gorm:"type:tinyint;default:1"`
	Description string        `json:"description" gorm:"type:text"`
	Status      TaskStatus    `json:"status" gorm:"type:tinyint;default:2"` // 默认为停止状态
	ExecStatus  TaskExecStatus `json:"execStatus" gorm:"type:tinyint;default:1"` // 默认为待执行状态
	Priority    string        `json:"priority" gorm:"type:varchar(20);default:'medium'"`
	CronExpr    string        `json:"cronExpr" gorm:"type:varchar(100)"`
	NextRunTime *time.Time    `json:"nextRunTime"`
	LastRunTime *time.Time    `json:"lastRunTime"`
	TaskContent string        `json:"taskContent" gorm:"type:text"`
	TaskParams  string        `json:"taskParams" gorm:"type:text"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Type        int8          `json:"type"`
	Description string        `json:"description"`
	Status      int8          `json:"status"`
	ExecStatus  int8          `json:"execStatus"`
	Priority    string        `json:"priority"`
	CreateTime  time.Time     `json:"createTime"`
	CronExpr    string        `json:"cronExpr"`
	NextRunTime *time.Time    `json:"nextRunTime"`
	LastRunTime *time.Time    `json:"lastRunTime"`
	TaskContent string        `json:"taskContent"`
	TaskParams  string        `json:"taskParams"`
}

// ToResponse 转换为响应对象
func (t *Task) ToResponse() *TaskResponse {
	return &TaskResponse{
		ID:          t.ID,
		Name:        t.Name,
		Type:        int8(t.Type),
		Description: t.Description,
		Status:      int8(t.Status),
		ExecStatus:  int8(t.ExecStatus),
		Priority:    t.Priority,
		CreateTime:  t.CreatedAt,
		CronExpr:    t.CronExpr,
		NextRunTime: t.NextRunTime,
		LastRunTime: t.LastRunTime,
		TaskContent: t.TaskContent,
		TaskParams:  t.TaskParams,
	}
}

// String 实现 TaskType 的字符串方法
func (t TaskType) String() string {
	return TaskTypeMap[t]
}

// String 实现 TaskStatus 的字符串方法
func (s TaskStatus) String() string {
	return TaskStatusMap[s]
}

// String 实现 TaskExecStatus 的字符串方法
func (s TaskExecStatus) String() string {
	return TaskExecStatusMap[s]
}

// TaskLog 任务日志模型
type TaskLog struct {
	gorm.Model
	TaskID     uint          `json:"taskId" gorm:"not null"`                // 任务ID
	Status     TaskExecStatus `json:"status" gorm:"type:tinyint;not null"` // 执行状态
	Output     string        `json:"output" gorm:"type:text"`              // 执行输出
	Error      string        `json:"error" gorm:"type:text"`               // 错误信息
	StartTime  time.Time     `json:"startTime" gorm:"not null"`            // 开始时间
	EndTime    time.Time     `json:"endTime" gorm:"not null"`             // 结束时间
	Duration   int64         `json:"duration" gorm:"not null"`             // 执行时长（秒）
}

// TaskLogResponse 任务日志响应
type TaskLogResponse struct {
	ID        uint          `json:"id"`
	TaskID    uint          `json:"taskId"`
	Status    int8          `json:"status"`
	Output    string        `json:"output"`
	Error     string        `json:"error"`
	StartTime time.Time     `json:"startTime"`
	EndTime   time.Time     `json:"endTime"`
	Duration  int64         `json:"duration"`
}

// ToResponse 转换为响应对象
func (l *TaskLog) ToResponse() *TaskLogResponse {
	return &TaskLogResponse{
		ID:        l.ID,
		TaskID:    l.TaskID,
		Status:    int8(l.Status),
		Output:    l.Output,
		Error:     l.Error,
		StartTime: l.StartTime,
		EndTime:   l.EndTime,
		Duration:  l.Duration,
	}
}
