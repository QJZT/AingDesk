package main

import (
  "github.com/gin-gonic/gin"
  "gorm.io/gorm"
  "flag"
  "github.com/glebarez/sqlite"
)

type Product struct {
  gorm.Model
  Code  string `json:"code"`
  Price uint   `json:"price"`
}

func main() {

  env := flag.String("p", "", "指定数据库环境 (dev/test/prod)")	
  flag.Parse()
//   sqlite-data.db
  db, err := gorm.Open(sqlite.Open(*env), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }
  db.AutoMigrate(&Product{})

  // 创建Gin路由
  r := gin.Default()

  // 定义API路由
  r.GET("/products", func(c *gin.Context) {
    var products []Product
    db.Find(&products)
    c.JSON(200, products)
  })

  r.POST("/products", func(c *gin.Context) {
    var product Product
    if err := c.ShouldBindJSON(&product); err != nil {
      c.JSON(400, gin.H{"error": err.Error()})
      return
    }
    db.Create(&product)
    c.JSON(200, product)
  })
  	// 添加ping接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

  // 启动服务
  r.Run(":7072")
}