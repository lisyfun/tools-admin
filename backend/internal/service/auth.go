package service

import (
	"errors"
	"time"

	"tools-admin/backend/internal/model"
	"tools-admin/backend/pkg/auth"
	"tools-admin/backend/pkg/db"
	"tools-admin/backend/pkg/log"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrUserAlreadyExists  = errors.New("username already exists")
)

type AuthService struct {
	jwtConfig auth.JWTConfig
}

func NewAuthService(secret string, expire time.Duration) *AuthService {
	return &AuthService{
		jwtConfig: auth.JWTConfig{
			Secret: secret,
			Expire: expire,
		},
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	User      *model.User `json:"user"`
	ExpiresIn int64       `json:"expires_in"` // token过期时间（秒）
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	log.Info("用户尝试登录: %s", req.Username)

	// 查询用户
	var user model.User
	if err := db.Db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		log.Error("用户登录失败：用户不存在：%s, 错误：%v", req.Username, err)
		return nil, ErrInvalidCredentials
	}

	// 检查用户状态
	if user.Status != 1 {
		log.Warn("用户登录失败：账号已禁用：%s, 状态：%d", req.Username, user.Status)
		return nil, ErrUserDisabled
	}

	// 验证密码
	if !auth.CheckPassword(user.Password, req.Password) {
		log.Error("用户登录失败：密码错误：%s", req.Username)
		return nil, ErrInvalidCredentials
	}

	// 生成JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.RoleID, s.jwtConfig)
	if err != nil {
		log.Error("用户登录失败：生成token失败：%s, 错误：%v", req.Username, err)
		return nil, err
	}

	log.Info("用户登录成功：%s, ID：%d, 过期时间：%.0f秒",
		req.Username,
		user.ID,
		s.jwtConfig.Expire.Seconds(),
	)

	return &LoginResponse{
		Token:     token,
		User:      &user,
		ExpiresIn: int64(s.jwtConfig.Expire.Seconds()),
	}, nil
}

func (s *AuthService) Register(user *model.User) error {
	// 检查用户名是否已存在
	var existingUser model.User
	if err := db.Db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return ErrUserAlreadyExists
	}

	// 加密密码
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// 设置默认值
	user.Status = 1
	if user.Nickname == "" {
		user.Nickname = user.Username
	}

	// 创建用户
	if err := db.Db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

// ResetPassword 重置用户密码
func (s *AuthService) ResetPassword(username, newPassword string) error {
	log.Info("尝试重置密码: %s", username)

	// 查询用户
	var user model.User
	if err := db.Db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Error("重置密码失败：用户不存在：%s, 错误：%v", username, err)
		return ErrInvalidCredentials
	}

	// 检查用户状态
	if user.Status != 1 {
		log.Warn("重置密码失败：账号已禁用：%s, 状态：%d", username, user.Status)
		return ErrUserDisabled
	}

	// 生成新密码哈希
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		log.Error("重置密码失败：密码加密错误：%s, 错误：%v", username, err)
		return err
	}

	// 更新密码
	if err := db.Db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		log.Error("重置密码失败：数据库更新错误：%s, 错误：%v", username, err)
		return err
	}

	log.Info("密码重置成功：%s", username)
	return nil
}
