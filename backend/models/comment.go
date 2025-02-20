package models

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id"`
	Content   string    `gorm:"type:varchar(1000);not null" json:"content" binding:"required,max=1000"`
	PostID    uint      `gorm:"index" json:"post_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Post      Post      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"post"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
