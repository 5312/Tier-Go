package router

import (
	"tier_up/app/api/v1/controller"
	"tier_up/app/internal/middleware/auth"
	"tier_up/app/internal/middleware/jwt"
	"tier_up/app/internal/service"

	_ "tier_up/app/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine, db *gorm.DB) {
	// 创建JWT服务
	jwtService := jwt.NewJWTService()

	// 创建服务
	userService := service.NewUserService(db, jwtService)
	roleService := service.NewRoleService(db)

	// 创建控制器
	userController := controller.NewUserController(userService)
	roleController := controller.NewRole(roleService)

	// 设置API路由组
	api := r.Group("/api/v1")

	api.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	// 不需要认证的路由
	{
		// 用户认证
		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)
	}

	// 需要认证的路由
	authGroup := api.Group("")
	authGroup.Use(jwtService.JWTAuthMiddleware())
	{
		// 用户相关
		authGroup.GET("/user/info", userController.GetUserInfo)
		authGroup.PUT("/user/info", userController.UpdateUserInfo)
		authGroup.PUT("/user/password", userController.ChangePassword)

		// 需要权限验证的路由
		rbacGroup := authGroup.Group("")
		rbacGroup.Use(auth.AuthMiddleware())
		{
			// 用户角色管理
			rbacGroup.POST("/user/:id/role", userController.AssignRole)
			rbacGroup.DELETE("/user/:id/role", userController.RemoveRole)

			// 角色管理
			rbacGroup.POST("/role", roleController.CreateRole)
			rbacGroup.GET("/role/:id", roleController.GetRoleByID)
			rbacGroup.PUT("/role/:id", roleController.UpdateRole)
			rbacGroup.DELETE("/role/:id", roleController.DeleteRole)
			rbacGroup.GET("/roles", roleController.ListRoles)

			// 权限管理
			rbacGroup.POST("/permission", roleController.AddPermission)
			rbacGroup.DELETE("/permission", roleController.RemovePermission)
			rbacGroup.GET("/role-permissions/:name", roleController.GetPermissions)
		}
	}
}
