package api

import (
	"keep_learning_blog/middleware"
	"keep_learning_blog/models"
	"keep_learning_blog/service"
	"net/http"
	"strconv"

	"keep_learning_blog/utils/logger"

	"github.com/gin-gonic/gin"
)

// PostController 文章控制器
type PostController struct {
	postService service.PostService
}

// NewPostController 创建文章控制器
func NewPostController() *PostController {
	return &PostController{
		postService: service.PostService{},
	}
}

// CreatePost 创建文章
func (c *PostController) CreatePost(ctx *gin.Context) {
	var req models.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
		})).Error("Failed to bind post request")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Title = middleware.SanitizeText(req.Title)
	req.Content = middleware.SanitizeHTML(req.Content)

	userID, exists := ctx.Get("user_id")
	if !exists {
		logger.Log.Error("User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	post, err := c.postService.CreatePost(req.Title, req.Content, userID.(uint), req.TagNames)
	if err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"title":   req.Title,
		})).Error("Failed to create post")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"post_id": post.ID,
		"user_id": userID,
		"title":   post.Title,
	})).Info("Post created successfully")

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "post created successfully",
		"post":    post,
	})
}

// GetPost 获取单个文章
func (c *PostController) GetPost(ctx *gin.Context) {
	// 解析文章ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	// 获取文章
	post, err := c.postService.GetPost(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "post retrieved successfully",
		"post":    post,
	})
}

// GetAllPosts 获取文章列表
func (c *PostController) GetAllPosts(ctx *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	// 获取文章列表
	posts, total, err := c.postService.GetAllPosts(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "posts retrieved successfully",
		"posts":   posts,
		"total":   total,
	})
}

// UpdatePost 更新文章
func (c *PostController) UpdatePost(ctx *gin.Context) {
	// 解析文章ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	// 解析请求体
	var req models.UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 对标题使用普通文本过滤，对内容使用 HTML 过滤
	req.Title = middleware.SanitizeText(req.Title)
	req.Content = middleware.SanitizeHTML(req.Content)

	// 从上下文获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 更新文章
	post, err := c.postService.UpdatePost(uint(id), userID.(uint), req.Title, req.Content, req.TagNames)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "post updated successfully",
		"post":    post,
	})
}

// DeletePost 删除文章
func (c *PostController) DeletePost(ctx *gin.Context) {
	// 解析文章ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 删除文章
	if err := c.postService.DeletePost(uint(id), userID.(uint)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "post deleted successfully",
	})
}

// GetPostComments 获取文章的所有评论
func (c *PostController) GetPostComments(ctx *gin.Context) {
	postID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	// 获取文章的所有评论
	comments, total, err := c.postService.GetPostComments(uint(postID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "comments retrieved successfully",
		"comments": comments,
		"total":    total,
	})
}

/*
// GetPostTags 获取文章标签
func (c *PostController) GetPostTags(ctx *gin.Context) {
	postID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	tags, err := c.postService.GetPostTags(uint(postID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

*/
