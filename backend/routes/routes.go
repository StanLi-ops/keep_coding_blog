package routes

import (
	"keep_coding_blog/api"
	"keep_coding_blog/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, logger *logrus.Logger) {
	userController := api.NewUserController(logger)
	postController := api.NewPostController(logger)
	commentController := api.NewCommentController(logger)

	// 添加静态文件服务
	r.Static("/uploads", "./uploads")

	// 公开路由
	public := r.Group("/api")
	{
		// 用户相关
		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)

		// 文章相关（公开访问）
		public.GET("/posts", postController.GetPosts)
		public.GET("/posts/:id", postController.GetPost)

		// 评论相关（公开访问）
		public.GET("/comments", commentController.GetComments)

		// 搜索相关
		public.GET("/posts/search", postController.SearchPosts)
		public.GET("/tags", postController.GetTags)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// 文章相关（需要认证）
		protected.POST("/posts", postController.CreatePost)
		protected.PUT("/posts/:id", postController.UpdatePost)
		protected.DELETE("/posts/:id", postController.DeletePost)

		// 评论相关（需要认证）
		protected.POST("/comments", commentController.CreateComment)
		protected.PUT("/comments/:id", commentController.UpdateComment)
		protected.DELETE("/comments/:id", commentController.DeleteComment)
	}

}
