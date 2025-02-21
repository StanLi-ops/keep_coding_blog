package api

import (
	"keep_coding_blog/models"
	"keep_coding_blog/service"
	"net/http"
	"strconv"

	"keep_coding_blog/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TagController 标签控制器结构体
type TagController struct {
	tagService service.TagService
	logger     *logrus.Logger
}

// NewTagController 创建标签控制器实例
func NewTagController(logger *logrus.Logger) *TagController {
	return &TagController{
		tagService: service.TagService{},
		logger:     logger,
	}
}

// CreateTag 创建标签
func (c *TagController) CreateTag(ctx *gin.Context) {
	var req models.CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 对标签名使用普通文本过滤
	req.Name = middleware.SanitizeText(req.Name)

	tag, err := c.tagService.CreateTag(req.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag created successfully",
		"tag":     tag,
	})
}

// GetTag 获取单个标签
func (c *TagController) GetTag(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	tag, err := c.tagService.GetTag(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag retrieved successfully",
		"tag":     tag,
	})
}

// GetAllTags 获取所有标签
func (c *TagController) GetAllTags(ctx *gin.Context) {
	tags, err := c.tagService.GetAllTags()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tags"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tags retrieved successfully",
		"tags":    tags,
	})
}

// UpdateTag 更新标签
func (c *TagController) UpdateTag(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	var req models.UpdateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 对标签名使用普通文本过滤
	req.Name = middleware.SanitizeText(req.Name)

	tag, err := c.tagService.UpdateTag(uint(id), req.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag updated successfully",
		"tag":     tag,
	})
}

// DeleteTag 删除标签
func (c *TagController) DeleteTag(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := c.tagService.DeleteTag(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag deleted successfully",
	})
}
