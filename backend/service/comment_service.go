package service

import (
	"errors"
	"keep_coding_blog/db"
	"keep_coding_blog/models"
)

// CommentService 评论服务结构体
type CommentService struct{}

// CreateComment 创建评论 (insert)
func (s *CommentService) CreateComment(content string, postID, userID uint) (*models.Comment, error) {
	// 验证数据合法性
	if content == "" || postID == 0 || userID == 0 {
		return nil, errors.New("content, postID, and userID cannot be empty")
	}
	if len(content) > 1000 {
		return nil, errors.New("content cannot be longer than 1000 characters")
	}

	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查文章是否存在
	var post models.Post
	if err := tx.First(&post, postID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("post not found")
	}

	// 创建评论
	comment := &models.Comment{
		Content: content,
		PostID:  postID,
		UserID:  userID,
	}

	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 加载关联的用户信息
	if err := tx.Preload("User").First(comment, comment.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return comment, tx.Commit().Error
}

// GetAllComments 获取所有评论 (select)
func (s *CommentService) GetAllComments() ([]models.Comment, error) {
	var comments []models.Comment
	if err := db.DB.Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// UpdateComment 更新评论 (update)
func (s *CommentService) UpdateComment(commentID, userID uint, content string) (*models.Comment, error) {
	// 验证数据合法性
	if content == "" || commentID == 0 || userID == 0 {
		return nil, errors.New("content, commentID, and userID cannot be empty")
	}
	if len(content) > 1000 {
		return nil, errors.New("content cannot be longer than 1000 characters")
	}

	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查评论是否存在
	var comment models.Comment
	if err := tx.First(&comment, commentID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("comment not found")
	}

	// 检查是否是评论作者
	if comment.UserID != userID {
		tx.Rollback()
		return nil, errors.New("unauthorized to update this comment")
	}

	// 更新评论内容
	comment.Content = content
	if err := tx.Save(&comment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 重新加载评论信息，包括用户信息
	if err := tx.Preload("User").First(&comment, comment.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &comment, tx.Commit().Error
}

// DeleteComment 删除评论 (delete)
func (s *CommentService) DeleteComment(commentID, userID uint) error {
	// 验证数据合法性
	if commentID == 0 || userID == 0 {
		return errors.New("commentID and userID cannot be empty")
	}

	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查评论是否存在
	var comment models.Comment
	if err := tx.First(&comment, commentID).Error; err != nil {
		tx.Rollback()
		return errors.New("comment not found")
	}

	// 检查是否是评论作者
	if comment.UserID != userID {
		tx.Rollback()
		return errors.New("unauthorized to delete this comment")
	}

	// 删除评论
	if err := tx.Delete(&comment).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
