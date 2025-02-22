package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"keep_coding_blog/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// encryptionConfig 加密配置
// 为方便测试暂时不添加至router.go
type encryptionConfig struct {
	key []byte
}

// NewEncryption 创建加密中间件
func NewEncryption(cfg *config.Config) gin.HandlerFunc {
	config := &encryptionConfig{
		key: []byte(cfg.Security.EncryptionKey),
	}
	return config.EncryptionMiddleware()
}

// EncryptionMiddleware 加密中间件
func (c *encryptionConfig) EncryptionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 只处理包含请求体的请求（POST, PUT, PATCH 等）
		if ctx.Request.Body != nil && ctx.Request.ContentLength > 0 {
			fmt.Println("解密请求体", ctx.Request.Body)

			encryptedData, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				return
			}

			// 解密请求数据
			decryptedData, err := c.decrypt(encryptedData)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to decrypt request"})
				return
			}

			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedData))
		}

		// 使用自定义ResponseWriter来加密响应
		writer := &encryptionResponseWriter{
			ResponseWriter: ctx.Writer,
			config:         c,
		}
		ctx.Writer = writer

		ctx.Next()
	}
}

// 加密响应数据的ResponseWriter
type encryptionResponseWriter struct {
	gin.ResponseWriter
	config *encryptionConfig
}

// 加密响应数据
func (w *encryptionResponseWriter) Write(data []byte) (int, error) {
	// 加密响应数据
	encryptedData, err := w.config.encrypt(data)
	if err != nil {
		return 0, err
	}
	return w.ResponseWriter.Write(encryptedData)
}

// AES加密方法
func (c *encryptionConfig) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// AES解密方法
func (c *encryptionConfig) decrypt(ciphertext []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, err
	}

	nonce, ciphertextData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
