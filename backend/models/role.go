package models

import (
	"time"
)

// Role 角色模型
type Role struct {
	ID          uint         `gorm:"primarykey;autoIncrement" json:"id"`
	Name        string       `gorm:"type:varchar(50);not null;unique" json:"name" binding:"required,max=50"`
	Code        string       `gorm:"type:varchar(50);not null;unique" json:"code" binding:"required,max=50"`
	Description string       `gorm:"type:text" json:"description"`
	IsDefault   bool         `gorm:"default:false" json:"is_default,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE" json:"users,omitempty"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Code          string `json:"code" binding:"required"`
	Description   string `json:"description" binding:"required"`
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
	IsDefault     *bool  `json:"is_default"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsDefault   *bool  `json:"is_default"`
}

// UpdatePermissionsRequest 更新角色权限请求
type UpdatePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}
