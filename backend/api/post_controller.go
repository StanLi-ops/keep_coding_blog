package api

import (
	"keep_coding_blog/models"
	"keep_coding_blog/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PostController 文章控制器结构体
type PostController struct {
	postService service.PostService
	logger      *logrus.Logger
}

// NewPostController 创建文章控制器实例
func NewPostController(logger *logrus.Logger) *PostController {
	return &PostController{
		postService: service.PostService{},
		logger:      logger,
	}
}

// CreatePost 创建文章
func (c *PostController) CreatePost(ctx *gin.Context) {
	var req models.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	post, err := c.postService.CreatePost(req.Title, req.Content, userID.(uint), req.TagNames)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "post created successfully",
		"post":    post,
	})
}

// GetPost 获取单个文章
func (c *PostController) GetPost(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, err := c.postService.GetPost(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "post retrieved successfully",
		"post":    post,
	})
}

// GetAllPosts 获取文章列表
func (c *PostController) GetAllPosts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	posts, total, err := c.postService.GetAllPosts(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "posts retrieved successfully",
		"posts":   posts,
		"total":   total,
	})
}

// UpdatePost 更新文章
func (c *PostController) UpdatePost(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	var req models.UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	post, err := c.postService.UpdatePost(uint(id), userID.(uint), req.Title, req.Content, req.TagNames)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "post updated successfully",
		"post":    post,
	})
}

// DeletePost 删除文章
func (c *PostController) DeletePost(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

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

	comments, total, err := c.postService.GetPostComments(uint(postID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
