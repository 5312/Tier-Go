package controller

import (
	"net/http"
	"strconv"
	"tier_up/app/internal/service"

	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type Role struct {
	RoleService *service.RoleService
}

// NewRoleController 创建角色控制器
func NewRole(roleService *service.RoleService) *Role {
	return &Role{
		RoleService: roleService,
	}
}

// CreateRole 创建角色
// @Summary 创建新角色
// @Description 创建一个新的角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body service.RoleRequest true "角色信息"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "创建角色失败"
// @Router /role [post]
func (c *Role) CreateRole(ctx *gin.Context) {
	var req service.RoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	role, err := c.RoleService.CreateRole(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建角色失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建角色成功", "data": role})
}

// GetRoleByID 获取角色
// @Summary 获取角色详情
// @Description 根据ID获取角色详情
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]interface{} "角色详情"
// @Failure 400 {object} map[string]interface{} "无效的角色ID"
// @Failure 500 {object} map[string]interface{} "获取角色失败"
// @Router /role/{id} [get]
func (c *Role) GetRoleByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的角色ID"})
		return
	}

	role, err := c.RoleService.GetRoleByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取角色失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取角色成功", "data": role})
}

// UpdateRole 更新角色
// @Summary 更新角色信息
// @Description 根据ID更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param data body service.RoleRequest true "角色信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "更新角色失败"
// @Router /role/{id} [put]
func (c *Role) UpdateRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的角色ID"})
		return
	}

	var req service.RoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	role, err := c.RoleService.UpdateRole(uint(id), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新角色失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新角色成功", "data": role})
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 根据ID删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "无效的角色ID"
// @Failure 500 {object} map[string]interface{} "删除角色失败"
// @Router /role/{id} [delete]
func (c *Role) DeleteRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的角色ID"})
		return
	}

	if err := c.RoleService.DeleteRole(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除角色失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除角色成功"})
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 分页获取角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页数量，默认10"
// @Success 200 {object} map[string]interface{} "角色列表"
// @Failure 500 {object} map[string]interface{} "获取角色列表失败"
// @Router /roles [get]
func (c *Role) ListRoles(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	roles, total, err := c.RoleService.ListRoles(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取角色列表失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取角色列表成功",
		"data": gin.H{
			"list":      roles,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AddPermission 添加权限
// @Summary 添加权限
// @Description 为角色添加访问路径的权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body service.PermissionRequest true "权限信息"
// @Success 200 {object} map[string]interface{} "添加成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "添加权限失败"
// @Router /permission [post]
func (c *Role) AddPermission(ctx *gin.Context) {
	var req service.PermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := c.RoleService.AddPermission(req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "添加权限失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "添加权限成功"})
}

// RemovePermission 移除权限
// @Summary 移除权限
// @Description 移除角色的访问路径权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body service.PermissionRequest true "权限信息"
// @Success 200 {object} map[string]interface{} "移除成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "移除权限失败"
// @Router /permission [delete]
func (c *Role) RemovePermission(ctx *gin.Context) {
	var req service.PermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := c.RoleService.RemovePermission(req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "移除权限失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "移除权限成功"})
}

// GetPermissions 获取角色权限
// @Summary 获取角色权限
// @Description 获取指定角色的所有权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name path string true "角色名称"
// @Success 200 {object} map[string]interface{} "权限列表"
// @Failure 400 {object} map[string]interface{} "角色名称不能为空"
// @Failure 500 {object} map[string]interface{} "获取角色权限失败"
// @Router /role-permissions/{name} [get]
func (c *Role) GetPermissions(ctx *gin.Context) {
	roleName := ctx.Param("name")
	if roleName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "角色名称不能为空"})
		return
	}

	permissions, err := c.RoleService.GetPermissions(roleName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取角色权限失败: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取角色权限成功",
		"data":    permissions,
	})
}
