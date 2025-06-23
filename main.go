package main

import (
	"fmt"
	"tier_up/app/internal/config"
	"tier_up/app/internal/db"
)

func main() {
	fmt.Println("|---------------------------|")
	fmt.Println("|----------admin------------|")
	fmt.Println("|---------------------------|")
	// 初始化配置
	config := config.Config{}
	config.InitConfig()
	// 初始化数据库
	db.InitDB(config)
	//
}
