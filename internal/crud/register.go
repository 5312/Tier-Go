package crud

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCrudRoutes[T any](
	r *gin.RouterGroup,
	config RouteConfig,
	db *gorm.DB,
) {
	crud := &Crud[T]{DB: db}
	group := r.Group(config.Prefix)

	// 按需注册路由
	if config.Create {
		group.POST("/create", crud.Create)
	}
	if config.Update {
		group.PUT("/update/:id", crud.Update)
	}
	if config.Delete {
		group.DELETE("/delete/:id", crud.Delete)
	}
	if config.Page {
		group.GET("/page", crud.Page)
	}
}
