package models

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
