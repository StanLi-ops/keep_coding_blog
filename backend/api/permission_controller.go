package api

import (
	"keep_learning_blog/models"
	"keep_learning_blog/service"
	"net/http"
	"strconv"

	"keep_learning_blog/utils/logger"

	"github.com/gin-gonic/gin"
)

// PermissionController 权限控制器
type PermissionController struct {
	permissionService service.PermissionService
}

// NewPermissionController 创建权限控制器
func NewPermissionController() *PermissionController {
	return &PermissionController{
		permissionService: service.PermissionService{},
	}
}

// CreatePermission 创建权限
func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	// 解析请求体
	var req models.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
		})).Error("Failed to bind permission request")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建权限
	permission, err := c.permissionService.CreatePermission(req.Name, req.Code, req.Method, req.Path, req.Description, req.IsDefault)
	if err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
			"name":  req.Name,
			"code":  req.Code,
		})).Error("Failed to create permission")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"permission_id": permission.ID,
		"name":          permission.Name,
		"code":          permission.Code,
	})).Info("Permission created successfully")

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Permission created successfully",
		"permission": permission,
	})
}

// GetPermission 获取单个权限
func (c *PermissionController) GetPermission(ctx *gin.Context) {
	// 解析权限ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	// 获取权限
	permission, err := c.permissionService.GetPermission(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Permission retrieved successfully",
		"permission": permission,
	})
}

// GetAllPermissions 获取所有权限
func (c *PermissionController) GetAllPermissions(ctx *gin.Context) {
	// 获取所有权限
	permissions, err := c.permissionService.GetAllPermissions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Permissions retrieved successfully",
		"permissions": permissions,
	})
}

// UpdatePermission 更新权限
func (c *PermissionController) UpdatePermission(ctx *gin.Context) {
	// 解析权限ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	// 解析请求体
	var req models.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新权限
	permission, err := c.permissionService.UpdatePermission(uint(id), req.Name, req.Code, req.Description, req.IsDefault)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Permission updated successfully",
		"permission": permission,
	})
}

// DeletePermission 删除权限
func (c *PermissionController) DeletePermission(ctx *gin.Context) {
	// 解析权限ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	// 删除权限
	if err := c.permissionService.DeletePermission(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Permission deleted successfully",
	})
}
