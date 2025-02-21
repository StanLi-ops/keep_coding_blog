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
			public.POST("/register", userController.Register)    //注册
			public.POST("/login", userController.Login)          //登录
			public.POST("/refresh", userController.RefreshToken) //刷新token

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

				private.GET("/permissions", permissionController.GetAllPermissions) // 获取所有权限
				private.GET("/permission/:id", permissionController.GetPermission)  // 获取指定权限

				private.PUT("/permission/:id", permissionController.UpdatePermission) // 编辑指定权限

				private.DELETE("/permission/:id", permissionController.DeletePermission) // 删除指定权限

				// 角色相关
				private.POST("/role", roleController.CreateRole) // 创建角色

				private.GET("/roles", roleController.GetAllRoles) // 获取所有角色
				private.GET("/role/:id", roleController.GetRole)  // 获取指定角色

				private.PUT("/role/:id", roleController.UpdateRole)                    // 编辑指定角色
				private.PUT("/role/:id/permissions", roleController.UpdatePermissions) // 编辑指定角色权限

				private.DELETE("/role/:id", roleController.DeleteRole) // 删除指定角色

				// 标签相关
				private.POST("/tag", tagController.CreateTag) // 创建标签

				private.PUT("/tag/:id", tagController.UpdateTag) // 编辑指定标签

				private.DELETE("/tag/:id", tagController.DeleteTag) // 删除指定标签

				// 文章相关
				private.POST("/post", postController.CreatePost) //创建文章

				private.PUT("/post/:id", postController.UpdatePost) //编辑指定文章

				private.DELETE("/post/:id", postController.DeletePost) //删除指定文章

				// 评论相关
				private.POST("/comment", commentController.CreateComment) //创建评论

				private.PUT("/comment/:id", commentController.UpdateComment) //编辑指定评论

				private.DELETE("/comment/:id", commentController.DeleteComment) //删除指定评论
			}
		}

	}
}
