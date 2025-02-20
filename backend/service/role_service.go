package service

import (
	"context"
	"errors"
	"keep_coding_blog/db"
	"keep_coding_blog/models"
	"log"
	"slices"
)

// RoleService 角色服务结构体
type RoleService struct{}

// CreateRole 创建角色 (insert)
func (s *RoleService) CreateRole(name, code, description string, permissionIDs []uint) (*models.Role, error) {
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

	// 检查角色名或编码是否已存在
	var existingRole models.Role
	if err := tx.Where("name = ? OR code = ?", name, code).First(&existingRole).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("role or code already exists")
	}

	// 检查权限是否存在
	var permissions []models.Permission
	if err := tx.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("some permissions do not exist")
	}

	// 创建角色
	role := &models.Role{
		Name:        name,
		Code:        code,
		Description: description,
	}

	if err := tx.Create(role).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 为角色分配权限
	if err := tx.Model(&role).Association("Permissions").Append(&permissions); err != nil {
		tx.Rollback()
		return nil, err
	}

	return role, tx.Commit().Error
}

// GetRole 获取单个角色 (select)
func (s *RoleService) GetRole(id uint) (*models.Role, error) {
	var role models.Role
	if err := db.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetAllRoles 获取所有角色 (select)
func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := db.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// UpdateRole 更新角色 (update)
func (s *RoleService) UpdateRole(id uint, name, code, description string) (*models.Role, error) {
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

	// 检查角色名或编码是否已存在 (排除当前id的记录)
	var existingRole models.Role
	if err := tx.Where("(name = ? OR code = ?) AND id != ?", name, code, id).First(&existingRole).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("role with same name or code already exists")
	}

	// 查找要更新的角色
	var role models.Role
	if err := tx.First(&role, id).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("role not found")
	}

	// 更新角色信息
	role.Name = name
	role.Code = code
	role.Description = description

	if err := tx.Save(&role).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &role, tx.Commit().Error
}

// DeleteRole 删除角色 (delete)
func (s *RoleService) DeleteRole(id uint) error {
	// 验证数据合法性
	if id == 0 {
		return errors.New("invalid role id")
	}

	// 使用事务处理
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除角色
	if err := tx.Delete(&models.Role{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdatePermissions 更新角色权限 (update)
func (s *RoleService) UpdatePermissions(roleID uint, permissionIDs []uint) (*models.Role, error) {
	// 验证数据合法性
	if roleID == 0 || len(permissionIDs) == 0 {
		return nil, errors.New("invalid input parameters")
	}

	// 对权限ID进行去重
	slices.Sort(permissionIDs)
	uniquePermissionIDs := slices.Compact(permissionIDs)

	// 使用事务处理
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 确保包含必要的默认权限
	var defaultPermissions []models.Permission
	if err := tx.Where("is_default = ?", true).Find(&defaultPermissions).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, dp := range defaultPermissions {
		if !slices.Contains(uniquePermissionIDs, dp.ID) {
			uniquePermissionIDs = append(uniquePermissionIDs, dp.ID)
		}
	}

	// 获取角色信息
	var role models.Role
	if err := tx.First(&role, roleID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("role not found")
	}

	// 预加载所有相关权限
	var permissions []models.Permission
	if err := tx.Where("id IN ?", uniquePermissionIDs).Find(&permissions).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 验证是否所有权限都存在
	if len(permissions) != len(uniquePermissionIDs) {
		tx.Rollback()
		return nil, errors.New("some permissions do not exist")
	}

	// 清除现有权限并分配新权限
	if err := tx.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 获取该角色关联的所有用户
	var users []models.User
	if err := tx.Model(&role).Association("Users").Find(&users); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 清除这些用户的Redis权限缓存
	for _, user := range users {
		if err := db.DeleteUserPermissions(context.Background(), user.ID); err != nil {
			// 仅记录日志,不中断事务
			log.Printf("Failed to delete permission cache for user %d: %v", user.ID, err)
		}
	}

	return &role, tx.Commit().Error
}
