package api

import (
	"keep_coding_blog/models"
	"keep_coding_blog/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CommentController 评论控制器结构体
type CommentController struct {
	commentService service.CommentService
	logger         *logrus.Logger
}

// NewCommentController 创建评论控制器实例
func NewCommentController(logger *logrus.Logger) *CommentController {
	return &CommentController{
		commentService: service.CommentService{},
		logger:         logger,
	}
}

// CreateComment 创建评论
func (c *CommentController) CreateComment(ctx *gin.Context) {
	var req models.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 从上下文获取用户ID
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	comment, err := c.commentService.CreateComment(req.Content, req.PostID, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "comment created successfully",
		"comment": comment,
	})
}

// GetAllComments 获取所有评论
func (c *CommentController) GetAllComments(ctx *gin.Context) {
	comments, err := c.commentService.GetAllComments()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "comments retrieved successfully",
		"comments": comments,
	})
}

// UpdateComment 更新评论
func (c *CommentController) UpdateComment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var req models.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 从上下文获取用户ID
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	comment, err := c.commentService.UpdateComment(uint(id), userID, req.Content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "comment updated successfully",
		"comment": comment,
	})
}

// DeleteComment 删除评论
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// 从上下文获取用户ID
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.commentService.DeleteComment(uint(id), userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "comment deleted successfully",
	})
}
