package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"username" binding:"required,max=64"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-" binding:"required,max=255"`
	Email     string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"email" binding:"required,email,max=128"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Roles     []Role    `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE" json:"roles,omitempty"`
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// CreateUserRequest 创建用户请求结构体
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户请求结构体
type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// UpdateUserRolesRequest 更新用户角色请求结构体
type UpdateUserRolesRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}
