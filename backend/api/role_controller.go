package api

import (
	"keep_learning_blog/config"
	"keep_learning_blog/models"
	"keep_learning_blog/service"
	"keep_learning_blog/utils/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct {
	config      *config.Config
	roleService service.RoleService
}

// NewRoleController 创建角色控制器
func NewRoleController(config *config.Config) *RoleController {
	return &RoleController{
		config:      config,
		roleService: service.RoleService{},
	}
}

// CreateRole 创建角色
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req models.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
		})).Error("Failed to bind role request")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := c.roleService.CreateRole(req.Name, req.Code, req.Description, req.PermissionIDs, req.IsDefault)
	if err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
			"name":  req.Name,
			"code":  req.Code,
		})).Error("Failed to create role")

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"role_id": role.ID,
		"name":    role.Name,
		"code":    role.Code,
	})).Info("Role created successfully")

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role created successfully",
		"role":    role,
	})
}

// GetRole 获取单个角色
func (c *RoleController) GetRole(ctx *gin.Context) {
	// 解析角色ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	// 获取角色
	role, err := c.roleService.GetRole(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role retrieved successfully",
		"role":    role,
	})
}

// GetAllRoles 获取所有角色
func (c *RoleController) GetAllRoles(ctx *gin.Context) {
	// 获取所有角色
	roles, err := c.roleService.GetAllRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Roles retrieved successfully",
		"roles":   roles,
	})
}

// UpdateRole 更新角色
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	// 解析角色ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	// 解析请求体
	var req models.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新角色
	role, err := c.roleService.UpdateRole(uint(id), req.Name, req.Code, req.Description, req.IsDefault)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role updated successfully",
		"role":    role,
	})
}

// DeleteRole 删除角色
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	// 解析角色ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	// 删除角色
	if err := c.roleService.DeleteRole(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role deleted successfully",
	})
}

// UpdatePermissions 更新角色权限
func (c *RoleController) UpdatePermissions(ctx *gin.Context) {
	// 解析角色ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	// 解析请求体
	var req models.UpdatePermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新角色权限
	role, err := c.roleService.UpdatePermissions(uint(id), req.PermissionIDs, c.config)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role permissions updated successfully",
		"role":    role,
	})
}
