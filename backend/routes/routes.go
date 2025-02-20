package routes

import (
	"keep_coding_blog/api"
	"keep_coding_blog/config"
	"keep_coding_blog/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, logger *logrus.Logger) {
	userController := api.NewUserController(logger)
	postController := api.NewPostController(logger)
	commentController := api.NewCommentController(logger)
	tagController := api.NewTagController(logger)
	roleController := api.NewRoleController(logger)
	permissionController := api.NewPermissionController(logger)

	// 添加静态文件服务
	r.Static("/uploads", "./uploads")

	// API 版本控制
	v1 := r.Group("/api")

	// 博客
	blog := v1.Group("")
	{
		// 公有路由
		public := blog.Group("")
		{
			// 用户相关
			public.POST("/register", userController.Register)
			public.POST("/login", userController.Login)
			public.POST("/refresh", userController.RefreshToken)

			// 文章相关
			public.GET("/posts", postController.GetAllPosts)
			public.GET("/posts/:id", postController.GetPost)
			public.GET("/posts/:id/comments", postController.GetPostComments)

			// 标签相关
			public.GET("/tags", tagController.GetAllTags)
			public.GET("/tag/:id", tagController.GetTag)

		}

		// 私有路由
		private := blog.Group("")
		private.Use(middleware.TokenAuth(&config.GetConfig().JWT))
		{
			// 用户相关
			private.POST("/logout", userController.Logout)

			private.Use(middleware.RBACAuth())
			{
				// 用户相关
				private.POST("/user", userController.CreateUser) // 创建用户

				private.GET("/users", userController.GetAllUsers)                 // 获取所有用户
				private.GET("/user/:id", userController.GetUser)                  // 获取指定用户
				private.GET("/user/:id/posts", userController.GetUserPosts)       // 获取指定用户所有文章
				private.GET("/user/:id/comments", userController.GetUserComments) // 获取指定用户所有评论

				private.PUT("/user/:id", userController.UpdateUser)           // 编辑指定用户
				private.PUT("/user/:id/role", userController.UpdateUserRoles) // 编辑指定用户角色

				private.DELETE("/user/:id", userController.DeleteUser) // 删除指定用户

				// 权限相关
				private.POST("/permission", permissionController.CreatePermission) // 创建权限

				private.GET("/permissions", permissionController.GetAllPermissions)
				private.GET("/permission/:id", permissionController.GetPermission)

				private.PUT("/permission/:id", permissionController.UpdatePermission)

				private.DELETE("/permission/:id", permissionController.DeletePermission)

				// 角色相关
				private.POST("/role", roleController.CreateRole)

				private.GET("/roles", roleController.GetAllRoles)
				private.GET("/role/:id", roleController.GetRole)

				private.PUT("/role/:id", roleController.UpdateRole)
				private.PUT("/role/:id/permission", roleController.UpdatePermissions)

				private.DELETE("/role/:id", roleController.DeleteRole)

				// 标签相关
				private.POST("/tag", tagController.CreateTag)

				private.PUT("/tag/:id", tagController.UpdateTag)

				private.DELETE("/tag/:id", tagController.DeleteTag)

				// 文章相关
				private.POST("/post", postController.CreatePost)

				private.PUT("/post/:id", postController.UpdatePost)

				private.DELETE("/post/:id", postController.DeletePost)

				// 评论相关
				private.POST("/comment", commentController.CreateComment)

				private.PUT("/comment/:id", commentController.UpdateComment)

				private.DELETE("/comment/:id", commentController.DeleteComment)
			}
		}

	}
}
