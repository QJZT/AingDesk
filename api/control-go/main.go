package main

import (
	"control-go/global"
	"control-go/model"
	"control-go/router"
	"control-go/seed"
	"flag"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	p := flag.String("p", "test.db", "指定数据库环境")
	p1 := flag.String("p1", "control-go", "指定存储环境文件夹")
	p2 := flag.String("p2", "audio_modules", "指定音频模块文件夹") // 修改默认值为 "audio_modules"
	flag.Parse()

	// 打印参数值
	flag.Parse()
	// 确保启用 json1 扩展，并启用外键支持
	fmt.Println(*p)
	fmt.Println(*p1)
	db, err := gorm.Open(sqlite.Open(*p+"?_pragma=json1&_pragma=foreign_keys(1)"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 启用调试日志
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Name{}, &model.FileData{}, &model.BaseModule{}, &model.Product{}, &model.KvStr{})

	// 检查并创建文件夹
	global.ToneFilePath = *p1
	if _, err := os.Stat(*p1); os.IsNotExist(err) {
		if err := os.Mkdir(*p1, 0755); err != nil {
			panic("无法创建存储文件夹: " + err.Error())
		}
	}

	// 设置 AudioFilePath 并检查/创建目录
	global.AudioFilePath = *p2
	if _, err := os.Stat(*p2); os.IsNotExist(err) {
		if err := os.Mkdir(*p2, 0755); err != nil {
			panic("无法创建音频模块文件夹: " + err.Error())
		}
	}

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
	// 添加CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.SetupBaseModuleRoutes(r, db)
	router.SetupProductRoutes(r, db)
	router.SetupPingRoutes(r)
	router.SetupNameRoutes(r, db)
	router.SetupKvRoutes(r, db)
	router.ChatAiRoutes(r)
	err = r.Run(":7072")
	fmt.Println(err.Error())
}
