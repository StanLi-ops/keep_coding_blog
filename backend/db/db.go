package db

import (
	"fmt"
	"keep_coding_blog/config"
	"keep_coding_blog/models"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(cfg *config.Config, logger *logrus.Logger) error {
	// 构建数据库连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.WithError(err).Error("Failed to connect to database")
		return err
	}

	// 自动迁移数据库结构
	err = db.AutoMigrate(
		&models.Permission{}, &models.Role{}, &models.User{},
		&models.Post{}, &models.Tag{}, &models.Comment{},
	)
	if err != nil {
		logger.WithError(err).Error("Failed to migrate database")
		return err
	}

	// 初始化基础数据
	if err := InitBaseData(db, logger); err != nil {
		logger.WithError(err).Error("Failed to initialize base data")
		return err
	}

	DB = db
	logger.Info("Database connected and migrated successfully")
	return nil
}

// InitBaseData 初始化基础数据
func InitBaseData(db *gorm.DB, logger *logrus.Logger) error {
	// 初始化权限
	permissions := []models.Permission{
		// 用户管理权限
		{Name: "查看用户", Code: "user:select", Method: "GET", Path: "/user", Description: "查看用户的信息"},
		{Name: "创建用户", Code: "user:create", Method: "POST", Path: "/user", Description: "创建新用户"},
		{Name: "编辑用户", Code: "user:edit", Method: "PUT", Path: "/user", Description: "编辑用户信息"},
		{Name: "删除用户", Code: "user:delete", Method: "DELETE", Path: "/user", Description: "删除用户"},

		// 权限管理权限
		{Name: "查看权限", Code: "permission:select", Method: "GET", Path: "/permission", Description: "查看权限信息"},
		{Name: "创建权限", Code: "permission:create", Method: "POST", Path: "/permission", Description: "创建新权限"},
		{Name: "编辑权限", Code: "permission:edit", Method: "PUT", Path: "/permission", Description: "编辑权限信息"},
		{Name: "删除权限", Code: "permission:delete", Method: "DELETE", Path: "/permission", Description: "删除权限"},

		// 角色管理权限
		{Name: "查看角色", Code: "role:select", Method: "POST", Path: "/role", Description: "查看角色信息"},
		{Name: "创建角色", Code: "role:create", Method: "POST", Path: "/role", Description: "创建新角色"},
		{Name: "编辑角色", Code: "role:edit", Method: "PUT", Path: "/role", Description: "编辑角色信息"},
		{Name: "删除角色", Code: "role:delete", Method: "DELETE", Path: "/role", Description: "删除角色"},

		// 标签管理权限
		{Name: "查看标签", Code: "tag:select", Method: "GET", Path: "/tag", Description: "查看标签", IsDefault: true},
		{Name: "创建标签", Code: "tag:create", Method: "POST", Path: "/tag", Description: "创建新标签"},
		{Name: "编辑标签", Code: "tag:edit", Method: "PUT", Path: "/tag", Description: "编辑标签"},
		{Name: "删除标签", Code: "tag:delete", Method: "DELETE", Path: "/tag", Description: "删除标签"},

		// 文章管理权限
		{Name: "查看文章", Code: "post:select", Method: "GET", Path: "/post", Description: "查看文章", IsDefault: true},
		{Name: "创建文章", Code: "post:create", Method: "POST", Path: "/post", Description: "创建新文章", IsDefault: true},
		{Name: "编辑文章", Code: "post:edit", Method: "PUT", Path: "/post", Description: "编辑文章", IsDefault: true},
		{Name: "删除文章", Code: "post:delete", Method: "DELETE", Path: "/post", Description: "删除文章", IsDefault: true},

		// 评论管理权限
		{Name: "查看评论", Code: "comment:select", Method: "GET", Path: "/comment", Description: "查看评论", IsDefault: true},
		{Name: "创建评论", Code: "comment:create", Method: "POST", Path: "/comment", Description: "创建新评论", IsDefault: true},
		{Name: "编辑评论", Code: "comment:edit", Method: "PUT", Path: "/comment", Description: "编辑评论", IsDefault: true},
		{Name: "删除评论", Code: "comment:delete", Method: "DELETE", Path: "/comment", Description: "删除评论", IsDefault: true},
	}

	// 使用FirstOrCreate避免重复创建
	for _, perm := range permissions {
		if err := db.Where("code = ?", perm.Code).FirstOrCreate(&perm).Error; err != nil {
			fmt.Println(err)
			return err
		}
	}

	// 初始化角色
	roles := []models.Role{
		{
			Name:        "超级管理员",
			Code:        "SUPER_ADMIN",
			Description: "系统超级管理员，拥有所有权限",
		},
		{
			Name:        "内容管理员",
			Code:        "CONTENT_ADMIN",
			Description: "内容管理员，负责管理文章和评论",
		},
		{
			Name:        "普通用户",
			Code:        "USER",
			Description: "普通用户，可以发布文章和评论",
			IsDefault:   true,
		},
	}

	// 使用FirstOrCreate避免重复创建
	for _, role := range roles {
		if err := db.Where("code = ?", role.Code).FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}

	// 为超级管理员角色分配所有权限
	var superAdmin models.Role
	if err := db.Where("code = ?", "SUPER_ADMIN").First(&superAdmin).Error; err != nil {
		return err
	}

	var allPermissions []models.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	if err := db.Model(&superAdmin).Association("Permissions").Replace(allPermissions); err != nil {
		return err
	}

	// 为内容管理员分配权限
	var contentAdmin models.Role
	if err := db.Where("code = ?", "CONTENT_ADMIN").First(&contentAdmin).Error; err != nil {
		return err
	}

	var contentPermissions []models.Permission
	if err := db.Where("code LIKE ? OR code LIKE ? OR code LIKE ?",
		"post:%", "comment:%", "tag:%").
		Find(&contentPermissions).Error; err != nil {
		return err
	}

	if err := db.Model(&contentAdmin).Association("Permissions").Replace(contentPermissions); err != nil {
		return err
	}

	// 为普通用户分配基本权限
	var user models.Role
	if err := db.Where("code = ?", "USER").First(&user).Error; err != nil {
		return err
	}

	var defaultPermissions []models.Permission
	if err := db.Where("is_default = ?", true).Find(&defaultPermissions).Error; err != nil {
		return err
	}

	if err := db.Model(&user).Association("Permissions").Replace(defaultPermissions); err != nil {
		return err
	}

	// 创建默认管理员用户
	password := "123456"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	superAdminUser := models.User{
		Username: "SuperAdmin",
		Password: string(passwordHash),
		Email:    "SuperAdmin@example.com",
	}

	if err := db.Where("username = ?", superAdminUser.Username).FirstOrCreate(&superAdminUser).Error; err != nil {
		return err
	}

	// 为管理员用户分配管理员角色
	if err := db.Model(&superAdminUser).Association("Roles").Replace(&superAdmin); err != nil {
		return err
	}

	logger.Info("Base data initialized successfully")
	return nil
}

// GetDB 返回数据库连接实例
func GetDB() *gorm.DB {
	return DB
}
