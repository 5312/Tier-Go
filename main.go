package main

import (
	"fmt"

	"tier_up/app/api/v1/router"
	_ "tier_up/app/docs" // 导入swagger文档
	"tier_up/app/internal/config"
	"tier_up/app/internal/db"
	"tier_up/app/internal/middleware/casbin"

	"github.com/gin-gonic/gin"
)

// @title           Tier Up API
// @version         1.0
// @description     Tier Up项目的API服务
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	fmt.Println("|---------------------------|")
	fmt.Println("|----------admin------------|")
	fmt.Println("|---------------------------|")

	// 初始化配置
	config := config.Config{}
	config.InitConfig()
	// 初始化数据库
	sqlDB, gormDB := db.InitDB(config)
	defer sqlDB.Close()

	// 初始化Casbin
	casbin.InitCasbin(gormDB)

	// 初始化Gin
	r := gin.Default()
	// 设置路由
	router.SetupRouter(r, gormDB)
	// 启动服务器
	addr := fmt.Sprintf("%s:%s", config.WebConfig.Host, config.WebConfig.Port)
	fmt.Printf("服务器启动在 %s\n", addr)
	fmt.Println("Swagger文档地址: http://" + addr + "/api/v1/swagger/index.html")
	if err := r.Run(addr); err != nil {
		fmt.Printf("启动服务器失败: %v\n", err)
	}

	r.Run(":" + config.WebConfig.Port)

}
