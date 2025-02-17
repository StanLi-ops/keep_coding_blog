package service

import (
	"errors"
	"keep_coding_blog/db"
	"keep_coding_blog/models"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务结构体
type UserService struct{}

// Register 注册用户
func (s *UserService) Register(username, password, email string) error {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if err := db.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建新用户
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Role:     "user", // 默认角色
	}

	return db.DB.Create(&user).Error
}

// Login 登录用户
func (s *UserService) Login(username, password string) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
