package middleware

import (
	"errors"
	"keep_learning_blog/config"
	"keep_learning_blog/db"
	"net/http"
	"strings"
	"time"

	"keep_learning_blog/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenAuther struct {
	config *config.JWTConfig
}

func NewTokenAuther(config *config.JWTConfig) *TokenAuther {
	return &TokenAuther{
		config: config,
	}
}

// TokenClaims 令牌声明
type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	TokenID   string `json:"token_id"`
	TokenType string `json:"token_type"`
	jwt.StandardClaims
}

// CreateTokenPair 生成访问令牌和刷新令牌
func CreateTokenPair(userID uint, username string, cfg *config.JWTConfig) (accessToken, refreshToken string, err error) {
	// 生成访问令牌
	accessTokenID := uuid.New().String()
	accessClaims := JWTClaims{
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
	refreshClaims := JWTClaims{
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
func RefreshJWTToken(c *gin.Context, cfg *config.JWTConfig) (string, string, error) {
	// 获取 refresh token
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		return "", "", errors.New("refresh token is required")
	}

	// 解析 refresh token
	claims := &JWTClaims{}
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
	return CreateTokenPair(claims.UserID, claims.Username, cfg)
}

// TokenAuth 令牌认证中间件
func (t *TokenAuther) TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Log.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Log.Warn("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// 获取 access token
		tokenString := parts[1]
		claims := &JWTClaims{}

		// 解析 access token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if claims.TokenType != "access" {
				return nil, errors.New("invalid token type")
			}
			return []byte(t.config.AccessTokenSecret), nil
		})

		// 检查 access token 是否有效
		if err != nil || !token.Valid {
			logger.Log.WithError(err).Warn("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 检查 access token 是否在黑名单中
		if db.IsBlacklisted(c, claims.TokenID) {
			logger.Log.WithFields(logger.Fields(map[string]interface{}{
				"token_id": claims.TokenID,
				"user_id":  claims.UserID,
			})).Warn("Token has been revoked")

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
