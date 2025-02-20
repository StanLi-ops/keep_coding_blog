package models

// Tag 标签模型
type Tag struct {
	ID    uint   `gorm:"primarykey;autoIncrement" json:"id"`
	Name  string `gorm:"type:varchar(50);unique;not null;index" json:"name" binding:"required,max=50"`
	Posts []Post `gorm:"many2many:post_tags;constraint:OnDelete:CASCADE" json:"-"`
}

// CreateTagRequest 创建标签请求结构体
type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateTagRequest 更新标签请求结构体
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required"`
}
