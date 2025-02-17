package api

import (
	"keep_coding_blog/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PostController 文章控制器
type PostController struct {
	postService *service.PostService
	logger      *logrus.Logger
}

// NewPostController 创建文章控制器
func NewPostController(logger *logrus.Logger) *PostController {
	return &PostController{
		postService: &service.PostService{},
		logger:      logger,
	}
}

// CreatePostRequest 创建文章请求结构体
type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	TagNames []string `json:"tags"`
}

// UpdatePostRequest 更新文章请求结构体
type UpdatePostRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	TagNames []string `json:"tags"`
}

// CreatePost 创建文章
func (c *PostController) CreatePost(ctx *gin.Context) {
	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("user_id")
	post, err := c.postService.CreatePost(req.Title, req.Content, userID, req.TagNames)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"post": post})
}

// GetPost 获取单个文章
func (c *PostController) GetPost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	post, err := c.postService.GetPost(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"post": post})
}

// GetPosts 获取文章列表
func (c *PostController) GetPosts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	posts, total, err := c.postService.GetPosts(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// UpdatePost 更新文章
func (c *PostController) UpdatePost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("user_id")
	post, err := c.postService.UpdatePost(uint(id), userID, req.Title, req.Content, req.TagNames)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"post": post})
}

// DeletePost 删除文章
func (c *PostController) DeletePost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	userID := ctx.GetUint("user_id")
	if err := c.postService.DeletePost(uint(id), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}

// SearchPosts 搜索文章
func (c *PostController) SearchPosts(ctx *gin.Context) {
	// 获取搜索参数
	query := ctx.Query("q")
	tags := ctx.QueryArray("tags")

	// 解析时间范围
	var startTime, endTime *time.Time
	if startStr := ctx.Query("start_time"); startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err == nil {
			startTime = &t
		}
	}
	if endStr := ctx.Query("end_time"); endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err == nil {
			endTime = &t
		}
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	posts, total, err := c.postService.SearchPosts(query, tags, startTime, endTime, page, pageSize)
	if err != nil {
		c.logger.WithError(err).Error("Failed to search posts")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    500001,
			Message: "Failed to search posts",
			Detail:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// GetTags 获取所有标签
func (c *PostController) GetTags(ctx *gin.Context) {
	tags, err := c.postService.GetAllTags()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tags": tags})
}

// 添加统一的错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}
