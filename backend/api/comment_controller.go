package api

import (
	"keep_coding_blog/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CommentController 评论控制器
type CommentController struct {
	commentService *service.CommentService
	logger         *logrus.Logger
}

// NewCommentController 创建评论控制器
func NewCommentController(logger *logrus.Logger) *CommentController {
	return &CommentController{
		commentService: &service.CommentService{},
		logger:         logger,
	}
}

// CreateCommentRequest 创建评论请求结构体
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	PostID  uint   `json:"post_id" binding:"required"`
}

// UpdateCommentRequest 更新评论请求结构体
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateComment 创建评论
func (c *CommentController) CreateComment(ctx *gin.Context) {
	var req CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("user_id")
	comment, err := c.commentService.CreateComment(req.Content, req.PostID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"comment": comment})
}

// GetComments 获取文章的评论列表
func (c *CommentController) GetComments(ctx *gin.Context) {
	// 从查询参数获取 post_id，而不是路径参数
	postIDStr := ctx.Query("post_id")
	if postIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}

	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	comments, total, err := c.commentService.GetCommentsByPost(uint(postID), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    total,
		"page":     page,
		"size":     pageSize,
	})
}

// UpdateComment 更新评论
func (c *CommentController) UpdateComment(ctx *gin.Context) {
	commentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	var req UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("user_id")
	comment, err := c.commentService.UpdateComment(uint(commentID), userID, req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"comment": comment})
}

// DeleteComment 删除评论
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	commentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.logger.WithError(err).Error("Invalid comment ID format")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	userID := ctx.GetUint("user_id")
	if err := c.commentService.DeleteComment(uint(commentID), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}
