package service

import (
	"errors"
	"keep_coding_blog/db"
	"keep_coding_blog/models"
)

// CommentService 评论服务结构体
type CommentService struct{}

// CreateComment 创建评论
func (s *CommentService) CreateComment(content string, postID, userID uint) (*models.Comment, error) {
	// 检查文章是否存在
	var post models.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		return nil, errors.New("post not found")
	}

	comment := &models.Comment{
		Content: content,
		PostID:  postID,
		UserID:  userID,
	}

	if err := db.DB.Create(comment).Error; err != nil {
		return nil, err
	}

	// 加载关联的用户信息
	if err := db.DB.Preload("User").First(comment, comment.ID).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

// GetCommentsByPost 获取文章的评论列表
func (s *CommentService) GetCommentsByPost(postID uint, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	// 获取评论总数
	if err := db.DB.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := db.DB.Where("post_id = ?", postID).
		Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(commentID, userID uint) error {
	var comment models.Comment
	if err := db.DB.First(&comment, commentID).Error; err != nil {
		return errors.New("comment not found")
	}

	// 检查是否是评论作者
	if comment.UserID != userID {
		return errors.New("unauthorized to delete this comment")
	}

	return db.DB.Delete(&comment).Error
}

// UpdateComment 更新评论
func (s *CommentService) UpdateComment(commentID, userID uint, content string) (*models.Comment, error) {
	var comment models.Comment
	if err := db.DB.First(&comment, commentID).Error; err != nil {
		return nil, errors.New("comment not found")
	}

	// 检查是否是评论作者
	if comment.UserID != userID {
		return nil, errors.New("unauthorized to update this comment")
	}

	comment.Content = content
	if err := db.DB.Save(&comment).Error; err != nil {
		return nil, err
	}

	// 重新加载评论信息，包括用户信息
	if err := db.DB.Preload("User").First(&comment, comment.ID).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}
