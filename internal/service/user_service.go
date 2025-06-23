package service

import (
	"errors"
	"strconv"
	"tier_up/app/internal/middleware/casbin"
	"tier_up/app/internal/middleware/jwt"
	"tier_up/app/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	DB         *gorm.DB
	JWTService *jwt.JWTService
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB, jwtService *jwt.JWTService) *UserService {
	return &UserService{
		DB:         db,
		JWTService: jwtService,
	}
}

// Register 注册用户
func (s *UserService) Register(req RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	var count int64
	if err := s.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := s.DB.Model(&model.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("邮箱已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1, // 默认激活状态
	}

	// 保存用户
	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(req LoginRequest) (string, *model.User, error) {
	var user model.User

	// 查找用户
	if err := s.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("用户不存在")
		}
		return "", nil, err
	}

	// 检查用户状态
	if user.Status != 1 {
		return "", nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", nil, errors.New("密码错误")
	}

	// 生成JWT令牌
	token, err := s.JWTService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

// GetUserByID 通过ID获取用户信息
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.DB.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *model.User) error {
	return s.DB.Save(user).Error
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user model.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	return s.DB.Save(&user).Error
}

// AssignRoleToUser 给用户分配角色
func (s *UserService) AssignRoleToUser(userID, roleID uint) error {
	var user model.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	var role model.Role
	if err := s.DB.First(&role, roleID).Error; err != nil {
		return err
	}

	// 添加角色关联
	if err := s.DB.Model(&user).Association("Roles").Append(&role); err != nil {
		return err
	}

	// 同时添加Casbin规则
	cs := casbin.GetInstance()
	_, err := cs.AddRoleForUser(strconv.Itoa(int(userID)), role.Name)
	return err
}

// RemoveRoleFromUser 从用户移除角色
func (s *UserService) RemoveRoleFromUser(userID, roleID uint) error {
	var user model.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	var role model.Role
	if err := s.DB.First(&role, roleID).Error; err != nil {
		return err
	}

	// 移除角色关联
	if err := s.DB.Model(&user).Association("Roles").Delete(&role); err != nil {
		return err
	}

	// 同时移除Casbin规则
	cs := casbin.GetInstance()
	_, err := cs.DeleteRoleForUser(strconv.Itoa(int(userID)), role.Name)
	return err
}
