package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Name struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var names = []Name{}
var lastID = 0

func SetupNameRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/names", func(c *gin.Context) {
		var names []Name
		db.Find(&names)
		c.JSON(http.StatusOK, names)
	})
	r.POST("/set_name", func(c *gin.Context) {

		//读取请求的name string {"name": "123"}
		var data = struct {
			Name string `json:"name"`
		}{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//写入数据库
		db.Create(&Name{Name: data.Name})
		//返回成功
		c.JSON(http.StatusOK, gin.H{"message": "name set successfully"})
	})
	r.POST("/del_name", func(c *gin.Context) {
		var data = struct {
			Id int `json:"id"`
		}{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Where("id = ?", data.Id).Delete(&Name{}) //删除指定数据
		c.JSON(http.StatusOK, gin.H{"message": "name delete successfully"})
	})

}
