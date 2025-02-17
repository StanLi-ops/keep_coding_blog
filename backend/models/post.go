package models

import (
	"time"
)

// Post 文章模型
type Post struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Tags      []Tag     `gorm:"many2many:post_tags;" json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Tag 标签模型
type Tag struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Name  string `gorm:"unique;not null" json:"name"`
	Posts []Post `gorm:"many2many:post_tags;" json:"-"`
}
