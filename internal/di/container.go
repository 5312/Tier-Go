package di

import (
	"tier-up/api/v1/controller"
	"tier-up/internal/middleware/jwt"
	"tier-up/internal/service"

	"go.uber.org/dig"
	"gorm.io/gorm"
)

func BuildContainer(db *gorm.DB) *dig.Container {
	container := dig.New()

	// 基础服务
	container.Provide(func() *gorm.DB { return db })
	container.Provide(jwt.NewJWTService)

	// 业务服务
	container.Provide(service.NewUserService)
	container.Provide(service.NewRoleService)

	// 控制器
	container.Provide(controller.NewUserController)
	container.Provide(controller.NewRoleController)

	return container
}
