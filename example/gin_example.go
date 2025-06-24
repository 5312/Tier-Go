package main

import (
	"net/http"
	"tier-up/internal/config"
	"tier-up/internal/crud"
	"tier-up/internal/db"
	"tier-up/internal/model"
	"tier-up/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func main() {

	// 初始化配置
	config := config.Config{}
	config.InitConfig()
	// 初始化数据库
	sqlDB, gormDB := db.InitDB(config)
	defer sqlDB.Close()

	router := gin.Default()

	router.POST("/role/create", func(ctx *gin.Context) {
		var req service.RoleRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
			return
		}
		crud := crud.Crud[model.Role]{
			DB: gormDB,
		}
		var role model.Role

		copier.Copy(&role, &req)

		err := crud.Create(&role)
		if err != nil {
			ctx.JSON(200, gin.H{
				"message": err,
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(":88")
}
