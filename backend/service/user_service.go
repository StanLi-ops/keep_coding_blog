package service

import (
	"context"
	"errors"
	"keep_learning_blog/config"
	"keep_learning_blog/db"
	"keep_learning_blog/models"
	"log"

	"keep_learning_blog/utils/logger"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务结构体
type UserService struct{}

// Register 注册用户 (insert)
func (s *UserService) Register(username, password, email string) (*models.User, error) {
	log := logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"username": username,
		"email":    email,
	}))

	// 验证数据合法性
	if username == "" || password == "" || email == "" {
		return nil, errors.New("username, password and email cannot be empty")
	}
	if len(username) > 64 || len(password) > 64 || len(email) > 128 {
		return nil, errors.New("username, password and email cannot be longer than 64 and 128 characters")
	}

	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查用户名或邮箱是否已存在
	var existingUser models.User
	if err := tx.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("username or email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建新用户
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	if err := tx.Create(&user).Error; err != nil {
		log.WithError(err).Error("Failed to create user")
		tx.Rollback()
		return nil, err
	}

	// 获取默认角色
	var defaultRole models.Role
	if err := tx.Where("is_default = ?", true).First(&defaultRole).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 为新用户分配默认角色
	if err := tx.Model(&user).Association("Roles").Append(&defaultRole); err != nil {
		tx.Rollback()
		return nil, err
	}

	log.Info("User registered successfully")
	return &user, tx.Commit().Error
}

// CreateUser 创建用户 (insert)
func (s *UserService) CreateUser(username, password, email string, rolesID uint) (*models.User, error) {
	// 验证数据合法性
	if username == "" || password == "" || email == "" || rolesID == 0 {
		return nil, errors.New("username, password, email and rolesID cannot be empty")
	}
	if len(username) > 64 || len(password) > 64 || len(email) > 128 {
		return nil, errors.New("username, password and email cannot be longer than 64 and 128 characters")
	}

	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查用户名或邮箱是否已存在
	var existingUser models.User
	if err := tx.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("username or email already exists")
	}

	// 检查角色是否存在
	var existingRole models.Role
	if err := tx.Where("id = ?", rolesID).First(&existingRole).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("role not found")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建新用户
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 用户分配角色
	if err := tx.Model(&user).Association("Roles").Append(&existingRole); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &user, tx.Commit().Error
}

// Login 登录用户 (select)
func (s *UserService) Login(username, password string) (*models.User, error) {
	log := logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"username": username,
	}))

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.WithError(err).Warn("Login failed: user not found")
		return nil, errors.New("user not found")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.WithError(err).Warn("Login failed: invalid password")
		return nil, errors.New("invalid password")
	}

	log.Info("User logged in successfully")
	return &user, nil
}

// GetUser 获取用户及其角色和权限信息 (select)
func (s *UserService) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := db.DB.Preload("Roles.Permissions").First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// GetAllUsers 获取所有用户及其角色和权限信息 (select)
func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := db.DB.Preload("Roles.Permissions").Find(&users).Error; err != nil {
		return nil, errors.New("failed to get all users")
	}

	return users, nil
}

// UpdateUser 更新用户 (update)
func (s *UserService) UpdateUser(id uint, username, password, email string) (*models.User, error) {
	// 验证输入不为空
	if id == 0 || username == "" || password == "" || email == "" {
		return nil, errors.New("id, username, password and email cannot be empty")
	}
	if len(username) > 64 || len(password) > 64 || len(email) > 128 {
		return nil, errors.New("username, password and email cannot be longer than 64 and 128 characters")
	}

	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查用户名或邮箱是否已存在 (排除当前id的记录)
	var existingUser models.User
	if err := tx.Where("(username = ? OR email = ?) AND id != ?", username, email, id).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("username or email already exists")
	}

	// 查找要更新的角色
	var user models.User
	if err := tx.First(&user, id).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("user not found")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	user.Username = username
	user.Password = string(hashedPassword)
	user.Email = email

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &user, tx.Commit().Error
}

// DeleteUser 删除用户 (delete)
func (s *UserService) DeleteUser(id uint) error {
	// 验证输入不为空
	if id == 0 {
		return errors.New("invalid user id")
	}

	// 开始事务
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除用户
	if err := tx.Delete(&models.User{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UpdateUserRoles 修改用户角色 (update)
func (s *UserService) UpdateUserRoles(userID uint, roleID uint, cfg *config.Config) (*models.User, error) {
	// 验证输入不为空
	if userID == 0 || roleID == 0 {
		return nil, errors.New("userID and roleID cannot be empty")
	}

	// 使用事务处理
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取用户信息
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("user not found")
	}

	// 获取角色信息
	var role models.Role
	if err := tx.First(&role, roleID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("role not found")
	}

	// 清除现有角色并分配新角色
	if err := tx.Model(&user).Association("Roles").Replace(&role); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 清除用户Redis权限缓存
	if err := db.DeleteUserPermissions(context.Background(), userID, cfg); err != nil {
		// 仅记录日志,不中断事务
		log.Printf("Failed to delete permission cache: %v", err)
	}

	return &user, tx.Commit().Error
}

// GetUserPosts 获取用户发表的文章 (select)
func (s *UserService) GetUserPosts(userID uint) ([]models.Post, error) {
	var posts []models.Post
	if err := db.DB.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, errors.New("failed to get user posts")
	}

	return posts, nil
}

// GetUserComments 获取用户发表的评论 (select)
func (s *UserService) GetUserComments(userID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := db.DB.Where("user_id = ?", userID).Find(&comments).Error; err != nil {
		return nil, errors.New("failed to get user comments")
	}

	return comments, nil
}

// GetUserPermissions 获取用户权限 (select)
func (s *UserService) GetUserPermissions(userID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := db.DB.Distinct().
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN roles ON roles.id = role_permissions.role_id").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&permissions).Error; err != nil {
		return nil, errors.New("failed to get user permissions")
	}

	return permissions, nil
}
