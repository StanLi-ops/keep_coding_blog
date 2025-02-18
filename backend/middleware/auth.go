package middleware

import (
	"errors"
	"keep_coding_blog/config"
	"keep_coding_blog/db"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// TokenClaims 令牌声明
type TokenClaims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	TokenID   string `json:"token_id"`
	TokenType string `json:"token_type"`
	jwt.StandardClaims
}

// GenerateTokenPair 生成访问令牌和刷新令牌
func GenerateTokenPair(userID uint, username string, cfg *config.JWTConfig) (accessToken, refreshToken string, err error) {
	// 生成访问令牌
	accessTokenID := uuid.New().String()
	accessClaims := TokenClaims{
		UserID:    userID,
		Username:  username,
		TokenID:   accessTokenID,
		TokenType: "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cfg.AccessTokenTTL).Unix(),
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(cfg.AccessTokenSecret))
	if err != nil {
		return "", "", err
	}

	// 生成 refresh token
	refreshTokenID := uuid.New().String()
	refreshClaims := TokenClaims{
		UserID:    userID,
		Username:  username,
		TokenID:   refreshTokenID,
		TokenType: "refresh",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cfg.RefreshTokenTTL).Unix(),
		},
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(cfg.RefreshTokenSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken 刷新访问令牌
func RefreshToken(c *gin.Context, cfg *config.JWTConfig) (string, string, error) {
	// 获取 refresh token
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		return "", "", errors.New("refresh token is required")
	}

	// 解析 refresh token
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if claims.TokenType != "refresh" {
			return nil, errors.New("invalid token type")
		}
		return []byte(cfg.RefreshTokenSecret), nil
	})

	// 检查 refresh token 是否有效
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	// 检查 refresh token 是否在黑名单中
	if db.IsBlacklisted(c, claims.TokenID) {
		return "", "", errors.New("refresh token has been revoked")
	}

	// 将 refresh token 加入黑名单
	if err := db.AddToBlacklist(c, claims.TokenID, cfg.RefreshTokenTTL); err != nil {
		return "", "", errors.New("failed to invalidate refresh token")
	}

	// 生成新的令牌对
	return GenerateTokenPair(claims.UserID, claims.Username, cfg)
}

// AuthMiddleware 认证中间件
func AuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查 Authorization header 格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// 获取 access token
		tokenString := parts[1]
		claims := &TokenClaims{}

		// 解析 access token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if claims.TokenType != "access" {
				return nil, errors.New("invalid token type")
			}
			return []byte(cfg.AccessTokenSecret), nil
		})

		// 检查 access token 是否有效
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 检查 access token 是否在黑名单中
		if db.IsBlacklisted(c, claims.TokenID) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("token_id", claims.TokenID)

		c.Next()
	}
}
