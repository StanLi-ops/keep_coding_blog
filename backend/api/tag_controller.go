package api

import (
	"keep_learning_blog/models"
	"keep_learning_blog/service"
	"net/http"
	"strconv"

	"keep_learning_blog/middleware"

	"keep_learning_blog/utils/logger"

	"github.com/gin-gonic/gin"
)

// TagController 标签控制器
type TagController struct {
	tagService service.TagService
}

// NewTagController 创建标签控制器
func NewTagController() *TagController {
	return &TagController{
		tagService: service.TagService{},
	}
}

// CreateTag 创建标签
func (c *TagController) CreateTag(ctx *gin.Context) {
	var req models.CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
		})).Error("Failed to bind tag request")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req.Name = middleware.SanitizeText(req.Name)

	tag, err := c.tagService.CreateTag(req.Name)
	if err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
			"name":  req.Name,
		})).Error("Failed to create tag")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"tag_id": tag.ID,
		"name":   tag.Name,
	})).Info("Tag created successfully")

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag created successfully",
		"tag":     tag,
	})
}

// GetTag 获取单个标签
func (c *TagController) GetTag(ctx *gin.Context) {
	// 解析标签ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	// 获取标签
	tag, err := c.tagService.GetTag(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag retrieved successfully",
		"tag":     tag,
	})
}

// GetAllTags 获取所有标签
func (c *TagController) GetAllTags(ctx *gin.Context) {
	// 获取所有标签
	tags, err := c.tagService.GetAllTags()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tags"})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tags retrieved successfully",
		"tags":    tags,
	})
}

// UpdateTag 更新标签
func (c *TagController) UpdateTag(ctx *gin.Context) {
	// 解析标签ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	// 解析请求体
	var req models.UpdateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 对标签名使用普通文本过滤
	req.Name = middleware.SanitizeText(req.Name)

	// 更新标签
	tag, err := c.tagService.UpdateTag(uint(id), req.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag updated successfully",
		"tag":     tag,
	})
}

// DeleteTag 删除标签
func (c *TagController) DeleteTag(ctx *gin.Context) {
	// 解析标签ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	// 删除标签
	if err := c.tagService.DeleteTag(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag deleted successfully",
	})
}
