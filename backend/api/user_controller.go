package api

import (
	"keep_learning_blog/middleware"
	"keep_learning_blog/models"
	"keep_learning_blog/service"
	"net/http"
	"strconv"

	"keep_learning_blog/config"
	"keep_learning_blog/db"
	"keep_learning_blog/utils/logger"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	config      *config.Config
	userService *service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(config *config.Config) *UserController {
	return &UserController{
		config:      config,
		userService: &service.UserService{},
	}
}

// Register 注册用户
func (c *UserController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error": err.Error(),
		})).Error("Failed to bind register request")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Username = middleware.SanitizeText(req.Username)
	req.Email = middleware.SanitizeText(req.Email)

	user, err := c.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error":    err.Error(),
			"username": req.Username,
			"email":    req.Email,
		})).Error("Failed to register user")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	})).Info("User registered successfully")

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user registered successfully",
		"user":    user,
	})
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	// 解析请求体
	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建用户
	user, err := c.userService.CreateUser(req.Username, req.Password, req.Email, req.RoleID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
		"user":    user,
	})
}

// Login 登录用户
func (c *UserController) Login(ctx *gin.Context) {
	identifier := ctx.MustGet("login_identifier").(string)
	loginLimiter := ctx.MustGet("login_limiter").(*middleware.LoginLimiter)
	login_request := ctx.MustGet("login_request").(models.LoginRequest)

	user, err := c.userService.Login(login_request.Username, login_request.Password)
	if err != nil {
		logger.Log.WithFields(logger.Fields(map[string]interface{}{
			"error":      err.Error(),
			"username":   login_request.Username,
			"identifier": identifier,
		})).Warn("Failed login attempt")

		loginLimiter.RecordLoginAttempt(ctx, false, identifier)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	logger.Log.WithFields(logger.Fields(map[string]interface{}{
		"user_id":    user.ID,
		"username":   user.Username,
		"identifier": identifier,
	})).Info("User logged in successfully")

	// 登录成功，记录成功并清除失败计数
	loginLimiter.RecordLoginAttempt(ctx, true, identifier)

	// 生成令牌对
	accessToken, refreshToken, err := middleware.CreateTokenPair(user.ID, user.Username, &c.config.JWT)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message":       "login successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// RefreshToken 处理令牌刷新请求
func (c *UserController) RefreshToken(ctx *gin.Context) {
	// 刷新令牌
	accessToken, refreshToken, err := middleware.RefreshJWTToken(ctx, &c.config.JWT)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message":       "token refreshed successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Logout 处理登出请求
func (c *UserController) Logout(ctx *gin.Context) {
	// 获取 tokenID
	tokenID := ctx.GetString("token_id")
	if tokenID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No token found"})
		return
	}

	// 将 access token 加入黑名单
	if err := db.AddToBlacklist(ctx, tokenID, c.config.JWT.AccessTokenTTL); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// GetUser 获取用户信息
func (c *UserController) GetUser(ctx *gin.Context) {
	// 解析用户ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 获取用户
	user, err := c.userService.GetUser(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user retrieved successfully",
		"user":    user,
	})
}

// GetAllUsers 获取所有用户及其角色和权限信息
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	// 获取所有用户
	users, err := c.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "users info retrieved successfully",
		"users":   users,
	})
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// 解析用户ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 解析请求体
	var req models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 对用户名和邮箱使用普通文本过滤
	req.Username = middleware.SanitizeText(req.Username)
	req.Email = middleware.SanitizeText(req.Email)

	// 更新用户
	user, err := c.userService.UpdateUser(uint(id), req.Username, req.Password, req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"user":    user,
	})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// 解析用户ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 删除用户
	if err := c.userService.DeleteUser(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UpdateUserRoles 更新用户角色
func (c *UserController) UpdateUserRoles(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 解析请求体
	var req models.UpdateUserRolesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新用户角色
	user, err := c.userService.UpdateUserRoles(uint(id), req.RoleID, c.config)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user roles updated successfully",
		"user":    user,
	})
}

// GetUserPosts 获取用户发表的文章
func (c *UserController) GetUserPosts(ctx *gin.Context) {
	// 解析用户ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 获取用户发表的文章
	posts, err := c.userService.GetUserPosts(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"posts": posts})
}

// GetUserComments 获取用户发表的评论
func (c *UserController) GetUserComments(ctx *gin.Context) {
	// 解析用户ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 获取用户发表的评论
	comments, err := c.userService.GetUserComments(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"comments": comments})
}
