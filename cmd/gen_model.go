package main

import (
	"fmt"
	"tier_up/app/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// 初始化配置
	config := config.Config{}
	config.InitConfig()
	c := config
	// 初始化 generator
	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/dao", // 生成代码输出路径
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", c.DB.Host, c.DB.User, c.DB.Password, c.DB.DriverName, c.DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败")
	}
	g.UseDB(db)

	// 生成所有表对应的 model 和 CRUD
	g.GenerateAllTable()

	g.Execute()
}
