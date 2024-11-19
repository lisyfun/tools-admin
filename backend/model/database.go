package model

import "time"

// Database 数据库连接配置
type Database struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;not null;comment:数据库连接名称"`
	Type      string    `json:"type" gorm:"size:20;not null;comment:数据库类型(mysql/postgresql)"`
	Host      string    `json:"host" gorm:"size:255;not null;comment:主机地址"`
	Port      int       `json:"port" gorm:"not null;comment:端口"`
	Username  string    `json:"username" gorm:"size:50;not null;comment:用户名"`
	Password  string    `json:"password" gorm:"size:255;not null;comment:密码"`
	Database  string    `json:"database" gorm:"size:50;not null;comment:数据库名"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DatabaseListReq 数据库列表请求
type DatabaseListReq struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"pageSize" binding:"required,min=1,max=100"`
	Name     string `form:"name"`
	Type     string `form:"type"`
}

// DatabaseListResp 数据库列表响应
type DatabaseListResp struct {
	Total int64      `json:"total"`
	List  []Database `json:"list"`
}

// DatabaseCreateReq 创建数据库连接请求
type DatabaseCreateReq struct {
	Name     string `json:"name" binding:"required,max=50"`
	Type     string `json:"type" binding:"required,oneof=mysql postgresql"`
	Host     string `json:"host" binding:"required,max=255"`
	Port     int    `json:"port" binding:"required,min=1,max=65535"`
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,max=255"`
	Database string `json:"database" binding:"required,max=50"`
}

// DatabaseUpdateReq 更新数据库连接请求
type DatabaseUpdateReq struct {
	ID       string `json:"-"`
	Name     string `json:"name" binding:"required,max=50"`
	Type     string `json:"type" binding:"required,oneof=mysql postgresql"`
	Host     string `json:"host" binding:"required,max=255"`
	Port     int    `json:"port" binding:"required,min=1,max=65535"`
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,max=255"`
	Database string `json:"database" binding:"required,max=50"`
}

// DatabaseTestReq 测试数据库连接请求
type DatabaseTestReq struct {
	Type     string `json:"type" binding:"required,oneof=mysql postgresql"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Database string `json:"database" binding:"required"`
}

// QueryExecuteReq SQL查询请求
type QueryExecuteReq struct {
	DatabaseID uint   `json:"database_id" binding:"required"`
	SQL        string `json:"sql" binding:"required"`
}

// QueryExecuteResp SQL查询响应
type QueryExecuteResp struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

// TableListReq 获取表列表请求
type TableListReq struct {
	DatabaseID uint `form:"database_id" binding:"required"`
}

// TableInfo 表信息
type TableInfo struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// TableSchemaReq 获取表结构请求
type TableSchemaReq struct {
	DatabaseID uint   `form:"database_id" binding:"required"`
	TableName  string `form:"table_name" binding:"required"`
}

// ColumnInfo 列信息
type ColumnInfo struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Length        int64  `json:"length"`
	Nullable      bool   `json:"nullable"`
	IsPrimaryKey  bool   `json:"is_primary_key"`
	IsAutoIncrement bool `json:"is_auto_increment"`
	DefaultValue  string `json:"default_value"`
	Comment       string `json:"comment"`
}
