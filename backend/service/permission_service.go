package service

import (
	"errors"
	"fmt"
	"keep_coding_blog/db"
	"keep_coding_blog/models"
)

// PermissionService 权限服务结构体
type PermissionService struct{}

// CreatePermission 创建权限 (insert)
func (s *PermissionService) CreatePermission(name, code, description string) (*models.Permission, error) {
	// 验证数据合法性
	if name == "" || code == "" {
		return nil, errors.New("name and code cannot be empty")
	}
	if len(name) > 50 || len(code) > 50 {
		return nil, errors.New("name and code must be less than 50 characters")
	}

	// 使用事务处理
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查权限名或编码是否已存在
	var existingPermission models.Permission
	if err := tx.Where("name = ? OR code = ?", name, code).First(&existingPermission).Error; err == nil {
		tx.Rollback()
		return nil, fmt.Errorf("permission with name '%s' or code '%s' already exists", name, code)
	}

	// 创建权限
	permission := &models.Permission{
		Name:        name,
		Code:        code,
		Description: description,
	}

	if err := tx.Create(permission).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, tx.Commit().Error
}

// GetPermission 获取单个权限 (select)
func (s *PermissionService) GetPermission(id uint) (*models.Permission, error) {
	var permission models.Permission
	if err := db.DB.First(&permission, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get permission with id %d: %w", id, err)
	}
	return &permission, nil
}

// GetAllPermissions 获取所有权限 (select)
func (s *PermissionService) GetAllPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := db.DB.Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get all permissions: %w", err)
	}
	return permissions, nil
}

// UpdatePermission 更新权限 (update)
func (s *PermissionService) UpdatePermission(id uint, name, code, description string) (*models.Permission, error) {
	// 验证数据合法性
	if id == 0 || name == "" || code == "" {
		return nil, errors.New("invalid input parameters")
	}
	if len(name) > 50 || len(code) > 50 {
		return nil, errors.New("name and code must be less than 50 characters")
	}

	// 使用事务处理
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查权限名或编码是否已存在 (排除当前id的记录)
	var existingPermission models.Permission
	if err := tx.Where("(name = ? OR code = ?) AND id != ?", name, code, id).First(&existingPermission).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("permission with same name or code already exists")
	}

	// 查找要更新的权限
	var permission models.Permission
	if err := tx.First(&permission, id).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("permission not found")
	}

	// 更新权限信息
	permission.Name = name
	permission.Code = code
	permission.Description = description

	if err := tx.Save(&permission).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &permission, tx.Commit().Error
}

// DeletePermission 删除权限 (delete)
func (s *PermissionService) DeletePermission(id uint) error {
	// 验证数据合法性
	if id == 0 {
		return errors.New("invalid permission id")
	}

	// 使用事务处理
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除权限
	if err := tx.Delete(&models.Permission{}, id).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return tx.Commit().Error
}
