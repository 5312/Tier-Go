package service

import (
	"tier_up/app/internal/middleware/casbin"
	"tier_up/app/internal/model"

	"gorm.io/gorm"
)

// RoleService 角色服务
type RoleService struct {
	DB *gorm.DB
}

// RoleRequest 角色请求
type RoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
}

// PermissionRequest 权限请求
type PermissionRequest struct {
	Role   string `json:"role" binding:"required"`
	Path   string `json:"path" binding:"required"`
	Method string `json:"method" binding:"required"`
}

// NewRoleService 创建角色服务
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		DB: db,
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(req RoleRequest) (*model.Role, error) {
	role := &model.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
	}

	if err := s.DB.Create(role).Error; err != nil {
		return nil, err
	}

	return role, nil
}

// GetRoleByID 通过ID获取角色
func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := s.DB.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(id uint, req RoleRequest) (*model.Role, error) {
	var role model.Role
	if err := s.DB.First(&role, id).Error; err != nil {
		return nil, err
	}

	role.Name = req.Name
	role.DisplayName = req.DisplayName
	role.Description = req.Description

	if err := s.DB.Save(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	return s.DB.Delete(&model.Role{}, id).Error
}

// ListRoles 获取角色列表
func (s *RoleService) ListRoles(page, pageSize int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	if err := s.DB.Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := s.DB.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// AddPermission 添加权限
func (s *RoleService) AddPermission(req PermissionRequest) error {
	cs := casbin.GetInstance()
	_, err := cs.AddPolicy(req.Role, req.Path, req.Method)
	return err
}

// RemovePermission 移除权限
func (s *RoleService) RemovePermission(req PermissionRequest) error {
	cs := casbin.GetInstance()
	_, err := cs.RemovePolicy(req.Role, req.Path, req.Method)
	return err
}

// GetPermissions 获取角色的所有权限
func (s *RoleService) GetPermissions(roleName string) ([][]string, error) {
	cs := casbin.GetInstance()
	return cs.GetEnforcer().GetFilteredPolicy(0, roleName)
}
