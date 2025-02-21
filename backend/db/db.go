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
		{Name: "创建用户", Code: "user:create", Method: "POST", Path: "/user", Description: "创建新用户"},

		{Name: "查看所有用户", Code: "users:select", Method: "GET", Path: "/users", Description: "查看所有用户信息"},
		{Name: "查看指定用户", Code: "user:select", Method: "GET", Path: "/user/:id", Description: "查看指定用户信息"},
		{Name: "查看指定用户所有文章", Code: "user:select:posts", Method: "GET", Path: "/user/:id/posts", Description: "查看指定用户所有文章信息"},
		{Name: "查看指定用户所有评论", Code: "user:select:comments", Method: "GET", Path: "/user/:id/comments", Description: "查看指定用户所有评论信息"},

		{Name: "编辑指定用户", Code: "user:edit", Method: "PUT", Path: "/user/:id", Description: "编辑指定用户信息"},
		{Name: "编辑指定用户角色", Code: "user:edit:roles", Method: "PUT", Path: "/user/:id/role", Description: "编辑指定用户角色"},

		{Name: "删除指定用户", Code: "user:delete", Method: "DELETE", Path: "/user/:id", Description: "删除指定用户"},

		// 权限管理权限
		{Name: "创建权限", Code: "permission:create", Method: "POST", Path: "/permission", Description: "创建新权限"},

		{Name: "查看所有权限", Code: "permissions:select", Method: "GET", Path: "/permissions", Description: "查看所有权限信息"},
		{Name: "查看指定权限", Code: "permission:select", Method: "GET", Path: "/permission/:id", Description: "查看指定权限信息"},

		{Name: "编辑指定权限", Code: "permission:edit", Method: "PUT", Path: "/permission/:id", Description: "编辑指定权限信息"},

		{Name: "删除指定权限", Code: "permission:delete", Method: "DELETE", Path: "/permission/:id", Description: "删除指定权限"},

		// 角色管理权限
		{Name: "创建角色", Code: "role:create", Method: "POST", Path: "/role", Description: "创建新角色"},

		{Name: "查看所有角色", Code: "roles:select", Method: "GET", Path: "/roles", Description: "查看所有角色信息"},
		{Name: "查看指定角色", Code: "role:select", Method: "GET", Path: "/role/:id", Description: "查看指定角色信息"},

		{Name: "编辑指定角色", Code: "role:edit", Method: "PUT", Path: "/role/:id", Description: "编辑指定角色信息"},
		{Name: "编辑指定角色权限", Code: "role:edit:permissions", Method: "PUT", Path: "/role/:id/permissions", Description: "编辑指定角色权限"},

		{Name: "删除指定角色", Code: "role:delete", Method: "DELETE", Path: "/role/:id", Description: "删除指定角色"},

		// 标签管理权限
		{Name: "创建标签", Code: "tag:create", Method: "POST", Path: "/tag", Description: "创建新标签"},

		{Name: "编辑指定标签", Code: "tag:edit", Method: "PUT", Path: "/tag/:id", Description: "编辑指定标签信息"},

		{Name: "删除指定标签", Code: "tag:delete", Method: "DELETE", Path: "/tag/:id", Description: "删除指定标签"},

		// 文章管理权限
		{Name: "创建文章", Code: "post:create", Method: "POST", Path: "/post", Description: "创建新文章", IsDefault: true},

		{Name: "编辑指定文章", Code: "post:edit", Method: "PUT", Path: "/post/:id", Description: "编辑指定文章信息", IsDefault: true},

		{Name: "删除指定文章", Code: "post:delete", Method: "DELETE", Path: "/post/:id", Description: "删除指定文章", IsDefault: true},

		// 评论管理权限
		{Name: "创建评论", Code: "comment:create", Method: "POST", Path: "/comment", Description: "创建新评论", IsDefault: true},

		{Name: "编辑指定评论", Code: "comment:edit", Method: "PUT", Path: "/comment/:id", Description: "编辑指定评论信息", IsDefault: true},

		{Name: "删除指定评论", Code: "comment:delete", Method: "DELETE", Path: "/comment/:id", Description: "删除指定评论", IsDefault: true},
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
