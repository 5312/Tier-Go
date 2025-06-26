package crud

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCrudRoutes[T any](r *gin.RouterGroup, path string, db *gorm.DB) {
	newCrud := &Crud[T]{DB: db}
	group := r.Group(path)
	group.POST("/create", newCrud.Create)
	group.PUT("/update/:id", newCrud.Update)
	group.DELETE("/delete/:id", newCrud.Delete)
	group.GET("/page", newCrud.Page)
}
