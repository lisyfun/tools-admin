package service

import (
	"errors"
	"tools-admin/backend/model"
	"tools-admin/backend/pkg/log"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrUsernameExists    = errors.New("用户名已存在")
	ErrInvalidPassword   = errors.New("密码错误")
	ErrUserStatusInvalid = errors.New("用户状态无效")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Create 创建用户
func (s *UserService) Create(req *model.UserCreateReq) error {
	// 检查用户名是否存在
	var count int64
	if err := s.db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		log.Error("检查用户名是否存在失败: %v", err)
		return err
	}
	if count > 0 {
		return ErrUsernameExists
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Mobile:   req.Mobile,
		RoleID:   req.RoleID,
		Status:   1, // 默认启用
	}

	if err := s.db.Create(user).Error; err != nil {
		log.Error("创建用户失败: %v", err)
		return err
	}

	return nil
}

// Update 更新用户
func (s *UserService) Update(id uint, req *model.UserUpdateReq) error {
	user := &model.User{}
	if err := s.db.First(user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		log.Error("查询用户失败: %v", err)
		return err
	}

	// 更新用户信息
	updates := map[string]interface{}{}
	if req.Password != "" {
		updates["password"] = req.Password
	}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Mobile != "" {
		updates["mobile"] = req.Mobile
	}
	if req.RoleID > 0 {
		updates["role_id"] = req.RoleID
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}

	if err := s.db.Model(user).Updates(updates).Error; err != nil {
		log.Error("更新用户失败: %v", err)
		return err
	}

	return nil
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	result := s.db.Delete(&model.User{}, id)
	if result.Error != nil {
		log.Error("删除用户失败: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id uint) (*model.User, error) {
	user := &model.User{}
	if err := s.db.First(user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		log.Error("查询用户失败: %v", err)
		return nil, err
	}
	return user, nil
}

// GetByUsername 根据用户名获取用户
func (s *UserService) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}
	if err := s.db.Where("username = ?", username).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		log.Error("查询用户失败: %v", err)
		return nil, err
	}
	return user, nil
}

// List 获取用户列表
func (s *UserService) List(req *model.UserListReq) (*model.UserListResp, error) {
	resp := &model.UserListResp{}

	// 构建查询条件
	query := s.db.Model(&model.User{})
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Nickname != "" {
		query = query.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.Mobile != "" {
		query = query.Where("mobile = ?", req.Mobile)
	}
	if req.RoleID > 0 {
		query = query.Where("role_id = ?", req.RoleID)
	}
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 获取总数
	if err := query.Count(&resp.Total).Error; err != nil {
		log.Error("获取用户总数失败: %v", err)
		return nil, err
	}

	// 获取列表
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&resp.List).Error; err != nil {
		log.Error("获取用户列表失败: %v", err)
		return nil, err
	}

	return resp, nil
}

// VerifyPassword 验证用户密码
func (s *UserService) VerifyPassword(username, password string) (*model.User, error) {
	user, err := s.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if !user.CheckPassword(password) {
		return nil, ErrInvalidPassword
	}

	if user.Status != 1 {
		return nil, ErrUserStatusInvalid
	}

	return user, nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(req *model.ResetPasswordReq) error {
	// 检查用户是否存在
	user, err := s.GetByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		log.Error("查询用户失败: %v", err)
		return err
	}

	// 更新密码
	user.Password = req.Password
	if err := s.db.Save(user).Error; err != nil {
		log.Error("更新密码失败: %v", err)
		return err
	}

	return nil
}
