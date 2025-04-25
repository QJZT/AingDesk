package main

import (
	"control-go/model"
	"control-go/router"
	"flag"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	p := flag.String("p", "test.db", "指定数据库环境 (dev/test/prod)")
	flag.Parse()
	db, err := gorm.Open(sqlite.Open(*p), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Name{})

	r := gin.Default()
	// 添加CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.SetupProductRoutes(r, db)
	router.SetupPingRoutes(r)
	router.SetupNameRoutes(r, db)
	r.Run(":7072")
}
