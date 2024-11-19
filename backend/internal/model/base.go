package model

import (
	"time"
)

// Base 基础模型
type Base struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User 用户模型
type User struct {
	Base
	Username string `gorm:"size:32;not null;unique" json:"username"`
	Password string `gorm:"size:128;not null" json:"-"`
	Nickname string `gorm:"size:32" json:"nickname"`
	Email    string `gorm:"size:128" json:"email"`
	Phone    string `gorm:"size:16" json:"phone"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	RoleID   uint   `json:"role_id"`
	Role     *Role  `json:"role,omitempty"`
}

// Role 角色模型
type Role struct {
	Base
	Name   string   `gorm:"size:32;not null;unique" json:"name"`
	Desc   string   `gorm:"size:128" json:"desc"`
	Status int      `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	Menus  []*Menu  `gorm:"many2many:role_menus;" json:"menus,omitempty"`
	Users  []*User  `json:"users,omitempty"`
}

// Menu 菜单模型
type Menu struct {
	Base
	Name       string  `gorm:"size:32;not null" json:"name"`
	Path       string  `gorm:"size:128" json:"path"`
	Component  string  `gorm:"size:128" json:"component"`
	Sort       int     `gorm:"default:0" json:"sort"`
	ParentID   uint    `json:"parent_id"`
	Parent     *Menu   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children   []*Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Icon       string  `gorm:"size:32" json:"icon"`
	Status     int     `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	Hidden     bool    `gorm:"default:false" json:"hidden"`
	Roles      []*Role `gorm:"many2many:role_menus;" json:"roles,omitempty"`
}

// Task 任务模型
type Task struct {
	Base
	Name        string    `gorm:"size:64;not null" json:"name"`
	Type        string    `gorm:"size:16;not null" json:"type"` // datax, shell, http
	Content     string    `gorm:"type:text" json:"content"`
	Cron        string    `gorm:"size:32" json:"cron"`
	Status      int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	LastRunTime time.Time `json:"last_run_time"`
	CreatorID   uint      `json:"creator_id"`
	Creator     *User     `json:"creator,omitempty"`
}

// TaskLog 任务日志模型
type TaskLog struct {
	Base
	TaskID    uint      `json:"task_id"`
	Task      *Task     `json:"task,omitempty"`
	Status    int       `json:"status"` // 0: 失败, 1: 成功
	Output    string    `gorm:"type:text" json:"output"`
	Error     string    `gorm:"type:text" json:"error"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// SMSTemplate 短信模板模型
type SMSTemplate struct {
	Base
	Name      string `gorm:"size:64;not null" json:"name"`
	Content   string `gorm:"type:text;not null" json:"content"`
	Status    int    `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatorID uint   `json:"creator_id"`
	Creator   *User  `json:"creator,omitempty"`
}

// SMSRecipient 短信接收人模型
type SMSRecipient struct {
	Base
	Name    string `gorm:"size:32;not null" json:"name"`
	Phone   string `gorm:"size:16;not null" json:"phone"`
	Status  int    `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	GroupID uint   `json:"group_id"`
}

// Server 服务器模型
type Server struct {
	Base
	Name     string `gorm:"size:64;not null" json:"name"`
	Host     string `gorm:"size:128;not null" json:"host"`
	Port     int    `gorm:"default:22" json:"port"`
	Username string `gorm:"size:32" json:"username"`
	Password string `gorm:"size:128" json:"-"`
	Status   int    `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
}

// Database 数据库连接模型
type Database struct {
	Base
	Name     string `gorm:"size:64;not null" json:"name"`
	Type     string `gorm:"size:16;not null" json:"type"` // mysql, postgresql
	Host     string `gorm:"size:128;not null" json:"host"`
	Port     int    `json:"port"`
	Username string `gorm:"size:32" json:"username"`
	Password string `gorm:"size:128" json:"-"`
	DBName   string `gorm:"size:64" json:"db_name"`
	Status   int    `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
}

// OperationLog 操作日志模型
type OperationLog struct {
	Base
	UserID    uint   `json:"user_id"`
	User      *User  `json:"user,omitempty"`
	IP        string `gorm:"size:32" json:"ip"`
	Method    string `gorm:"size:16" json:"method"`
	Path      string `gorm:"size:128" json:"path"`
	Status    int    `json:"status"`
	Latency   int64  `json:"latency"` // 请求耗时(ms)
	Agent     string `gorm:"size:256" json:"agent"`
	Request   string `gorm:"type:text" json:"request"`  // 请求参数
	Response  string `gorm:"type:text" json:"response"` // 响应内容
}
