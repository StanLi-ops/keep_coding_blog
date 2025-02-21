package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

var (
	// strictPolicy 用于普通文本，完全转义所有HTML
	strictPolicy = bluemonday.StrictPolicy()

	// htmlPolicy 用于富文本，允许安全的HTML标签
	htmlPolicy = bluemonday.UGCPolicy()
)

func init() {
	// 配置富文本策略
	htmlPolicy.AllowStandardURLs()
	htmlPolicy.AllowStandardAttributes()
	// 允许常用的安全标签
	htmlPolicy.AllowElements("p", "br", "b", "i", "strong", "em", "ul", "ol", "li", "h1", "h2", "h3", "h4", "h5", "h6", "blockquote", "code", "pre", "hr")
	// 允许链接
	htmlPolicy.AllowAttrs("href").OnElements("a")
	// 允许图片
	htmlPolicy.AllowAttrs("src", "alt").OnElements("img")
	// 允许代码高亮相关属性
	htmlPolicy.AllowAttrs("class").OnElements("code", "pre")
}

// XSSProtection 中间件用于处理XSS防护
func XSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求前的操作
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			processRequest(c)
		}

		// 处理响应
		c.Writer = &xssResponseWriter{
			ResponseWriter: c.Writer,
			policy:         strictPolicy, // 默认使用严格策略
		}

		c.Next()
	}
}

// processRequest 处理请求数据
func processRequest(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")

	// 只处理 application/json 请求
	if contentType != "application/json" {
		return
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 重新设置请求体，因为读取后 body 会被清空
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
}

// xssResponseWriter 用于处理响应数据
type xssResponseWriter struct {
	gin.ResponseWriter
	policy *bluemonday.Policy
}

// Write 重写 Write 方法，对响应进行处理
func (w *xssResponseWriter) Write(data []byte) (int, error) {
	// 对响应数据进行处理
	sanitized := w.policy.SanitizeBytes(data)
	return w.ResponseWriter.Write(sanitized)
}

// SanitizeHTML 用于富文本内容的过滤，可在业务层调用
func SanitizeHTML(content string) string {
	return htmlPolicy.Sanitize(content)
}

// SanitizeText 用于普通文本的过滤，可在业务层调用
func SanitizeText(content string) string {
	return strictPolicy.Sanitize(content)
}
