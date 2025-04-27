package main

import (
	"control-go/global"
	"control-go/model"
	"control-go/router"
	"flag"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	p := flag.String("p", "test.db", "指定数据库环境")
	p1 := flag.String("p1", "control-go", "指定存储环境文件夹")
	flag.Parse()
	db, err := gorm.Open(sqlite.Open(*p), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Name{}, &model.FileData{})

	// 检查并创建文件夹
	global.ToneFilePath = *p1
	if _, err := os.Stat(*p1); os.IsNotExist(err) {
		if err := os.Mkdir(*p1, 0755); err != nil {
			panic("无法创建存储文件夹: " + err.Error())
		}
	}

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
