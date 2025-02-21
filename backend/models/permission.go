package models

import (
	"time"
)

// Permission 权限模型
type Permission struct {
	ID          uint      `gorm:"primarykey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null;unique" json:"name" binding:"required,max=50"`
	Code        string    `gorm:"type:varchar(50);not null;unique" json:"code" binding:"required,max=50"`
	Method      string    `gorm:"type:varchar(10);not null" json:"method" binding:"required,max=10"`
	Path        string    `gorm:"type:varchar(128);not null" json:"path" binding:"required,max=128"`
	Description string    `gorm:"type:text" json:"description"`
	IsDefault   bool      `gorm:"default:false" json:"is_default,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Roles       []Role    `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE" json:"roles,omitempty"`
}

// CreatePermissionRequest 创建权限请求结构体
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Method      string `json:"method" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsDefault   *bool  `json:"is_default"`
}

// UpdatePermissionRequest 更新权限请求结构体
type UpdatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsDefault   *bool  `json:"is_default"`
}
