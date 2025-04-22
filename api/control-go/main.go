package main

import (
	"control-go/router"
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	p := flag.String("p", "", "指定数据库环境 (dev/test/prod)")
	flag.Parse()
	db, err := gorm.Open(sqlite.Open(*p), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&router.Product{})

	r := gin.Default()

	router.SetupProductRoutes(r, db)
	router.SetupPingRoutes(r)

	r.Run(":7072")
}
