package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(32);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(128);not null" json:"-"` // 密码不返回给前端
	Nickname string `gorm:"type:varchar(32)" json:"nickname"`
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`
	Email    string `gorm:"type:varchar(128)" json:"email"`
	Mobile   string `gorm:"type:varchar(16)" json:"mobile"`
	RoleID   uint   `gorm:"not null" json:"role_id"`
	Status   int    `gorm:"type:tinyint;default:1" json:"status"` // 1: 正常, 0: 禁用
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeSave 保存前的钩子函数
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 如果密码被修改，则加密
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// UserCreateReq 创建用户请求
type UserCreateReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname" binding:"max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Mobile   string `json:"mobile" binding:"omitempty,len=11"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

// UserUpdateReq 更新用户请求
type UserUpdateReq struct {
	Password string `json:"password" binding:"omitempty,min=6,max=32"`
	Nickname string `json:"nickname" binding:"max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Mobile   string `json:"mobile" binding:"omitempty,len=11"`
	RoleID   uint   `json:"role_id"`
	Status   int    `json:"status" binding:"oneof=0 1"`
}

// UserListReq 用户列表请求
type UserListReq struct {
	Username string `form:"username"`
	Nickname string `form:"nickname"`
	Mobile   string `form:"mobile"`
	RoleID   uint   `form:"role_id"`
	Status   int    `form:"status"`
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
}

// UserListResp 用户列表响应
type UserListResp struct {
	Total int64   `json:"total"`
	List  []*User `json:"list"`
}

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}
