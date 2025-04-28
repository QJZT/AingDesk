package router

import (
	"control-go/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupProductRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/products", func(c *gin.Context) {
		var products []model.Product
		db.Find(&products)
		c.JSON(200, products)
	})

	r.POST("/products", func(c *gin.Context) {
		var product model.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return

		}
		db.Create(&product)
		c.JSON(200, product)
	})
}
