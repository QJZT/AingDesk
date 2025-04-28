package main

import (
	"control-go/model"
	"control-go/router"
	"control-go/seed"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	p := flag.String("p", "p.db", "指定数据库环境 (dev/test/prod)")
	flag.Parse()
	// 确保启用 json1 扩展，并启用外键支持
	db, err := gorm.Open(sqlite.Open(*p+"?_pragma=json1&_pragma=foreign_keys(1)"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 启用调试日志
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 ModuleConfig、BaseModule 和 Product 表，确保顺序正确
	if err := db.AutoMigrate(&model.ModuleConfig{}); err != nil {
		panic("failed to migrate ModuleConfig: " + err.Error())
	}
	if err := db.AutoMigrate(&model.BaseModule{}); err != nil {
		panic("failed to migrate BaseModule: " + err.Error())
	}
	if err := db.AutoMigrate(&model.Product{}); err != nil {
		panic("failed to migrate Product: " + err.Error())
	}
	fmt.Println("Database migration completed")

	// 插入假数据
	if err := seed.SeedBaseModule(db); err != nil {
		panic("failed to seed BaseModule data: " + err.Error())
	}
	fmt.Println("Seeded BaseModule data successfully")
	if err := seed.SeedProduct(db); err != nil {
		panic("failed to seed Product data: " + err.Error())
	}
	fmt.Println("Seeded Product data successfully")

	// 设置 Gin 为 release 模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	// 设置信任的代理
	r.SetTrustedProxies([]string{"127.0.0.1"})

	router.SetupBaseModuleRoutes(r, db)
	router.SetupProductRoutes(r, db)
	router.SetupPingRoutes(r)

	fmt.Println("Starting server on :7072")
	r.Run(":7072")
}
