package api

import (
	"keep_coding_blog/models"
	"keep_coding_blog/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PermissionController 权限控制器结构体
type PermissionController struct {
	permissionService service.PermissionService
	logger            *logrus.Logger
}

// NewPermissionController 创建权限控制器实例
func NewPermissionController(logger *logrus.Logger) *PermissionController {
	return &PermissionController{
		permissionService: service.PermissionService{},
		logger:            logger,
	}
}

// CreatePermission 创建权限
func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	var req models.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := c.permissionService.CreatePermission(req.Name, req.Code, req.Description)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Permission created successfully",
		"permission": permission,
	})
}

// GetPermission 获取单个权限
func (c *PermissionController) GetPermission(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	permission, err := c.permissionService.GetPermission(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Permission retrieved successfully",
		"permission": permission,
	})
}

// GetAllPermissions 获取所有权限
func (c *PermissionController) GetAllPermissions(ctx *gin.Context) {
	permissions, err := c.permissionService.GetAllPermissions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Permissions retrieved successfully",
		"permissions": permissions,
	})
}

// UpdatePermission 更新权限
func (c *PermissionController) UpdatePermission(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	var req models.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := c.permissionService.UpdatePermission(uint(id), req.Name, req.Code, req.Description)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Permission updated successfully",
		"permission": permission,
	})
}

// DeletePermission 删除权限
func (c *PermissionController) DeletePermission(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	if err := c.permissionService.DeletePermission(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permission deleted successfully",
	})
}
