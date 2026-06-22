package main

import (
	"log"
	"online-ordering-system/config"
	"online-ordering-system/routes"
)

func main() {
	config.InitDB()
	config.SeedData()

	r := routes.SetupRouter()
	log.Println("网上订餐管理系统启动中...")
	log.Println("前台地址: http://localhost:8080")
	log.Println("后台地址: http://localhost:8080/admin")
	log.Println("默认管理员: admin / admin123")
	r.Run(config.AppConfig.Port)
}
