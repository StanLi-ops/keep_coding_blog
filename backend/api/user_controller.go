package api

import (
	"keep_coding_blog/middleware"
	"keep_coding_blog/service"
	"net/http"

	"keep_coding_blog/config"
	"keep_coding_blog/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserController 用户控制器
type UserController struct {
	userService *service.UserService
	logger      *logrus.Logger
}

// NewUserController 创建用户控制器
func NewUserController(logger *logrus.Logger) *UserController {
	return &UserController{
		userService: &service.UserService{},
		logger:      logger,
	}
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 处理注册请求
func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.Register(req.Username, req.Password, req.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login 处理登录请求
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 生成令牌对
	cfg := config.GetConfig().JWT
	accessToken, refreshToken, err := middleware.GenerateTokenPair(user.ID, user.Username, &cfg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// RefreshToken 处理令牌刷新请求
func (c *UserController) RefreshToken(ctx *gin.Context) {
	cfg := config.GetConfig().JWT
	accessToken, refreshToken, err := middleware.RefreshToken(ctx, &cfg)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Logout 处理登出请求
func (c *UserController) Logout(ctx *gin.Context) {
	tokenID := ctx.GetString("token_id")
	if tokenID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No token found"})
		return
	}

	// 将 access token 加入黑名单
	cfg := config.GetConfig().JWT
	if err := db.AddToBlacklist(ctx, tokenID, cfg.AccessTokenTTL); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	// 将 refresh token 加入黑名单
	if err := db.AddToBlacklist(ctx, tokenID, cfg.RefreshTokenTTL); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
