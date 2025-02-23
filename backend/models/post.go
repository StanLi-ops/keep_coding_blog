package models

import (
	"time"
)

// Post 文章模型
type Post struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id"`
	Title     string    `gorm:"type:varchar(200);not null;index" json:"title" binding:"required,max=200"`
	Content   string    `gorm:"type:text" json:"content"`
	UserID    uint      `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Tags      []Tag     `gorm:"many2many:post_tags;constraint:OnDelete:CASCADE" json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	TagNames []string `json:"tagNames"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	TagNames []string `json:"tagNames"`
}
