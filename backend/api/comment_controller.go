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

// CommentController 评论控制器
type CommentController struct {
	commentService service.CommentService
}

// NewCommentController 创建评论控制器
func NewCommentController() *CommentController {
	return &CommentController{
		commentService: service.CommentService{},
	}
}

// CreateComment 创建评论
func (c *CommentController) CreateComment(ctx *gin.Context) {
	var req models.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
		})).Error("Failed to bind comment request")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req.Content = middleware.SanitizeHTML(req.Content)
	userID := ctx.GetUint("user_id")

	comment, err := c.commentService.CreateComment(req.Content, req.PostID, userID)
	if err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"post_id": req.PostID,
		})).Error("Failed to create comment")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"comment_id": comment.ID,
		"user_id":    userID,
		"post_id":    req.PostID,
	})).Info("Comment created successfully")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "comment created successfully",
		"comment": comment,
	})
}

// UpdateComment 更新评论
func (c *CommentController) UpdateComment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
			"id":    ctx.Param("id"),
		})).Error("Invalid comment ID")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var req models.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
		})).Error("Failed to bind update comment request")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req.Content = middleware.SanitizeHTML(req.Content)
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
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
			"id":    ctx.Param("id"),
		})).Error("Invalid comment ID")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

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
