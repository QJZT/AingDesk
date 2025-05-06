package router

import (
	"net/http"

	"control-go/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupKvRoutes(r *gin.Engine, db *gorm.DB) {
	//获取数据map值
	r.POST("/get_kv", func(c *gin.Context) {
		var data = struct {
			Key string `json:"key"`
		}{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		KvStr := model.KvStr{}
		db.Where("key =?", data.Key).Find(&KvStr) //删除指定数据

		//如果没数据就创建个空的
		if KvStr.Key == "" {
			db.Create(&model.KvStr{Key: data.Key, V: make(map[string]interface{})})
			c.JSON(http.StatusOK, gin.H{"kv": make(map[string]interface{})})
			return
		}
		c.JSON(http.StatusOK, gin.H{"kv": KvStr.V})
	})
	//获取数据map值
	r.POST("/set_kv", func(c *gin.Context) {
		var data = struct {
			Key string                 `json:"key"`
			Kv  map[string]interface{} `json:"kv"`
		}{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		KvStr := model.KvStr{}
		db.Where("key =?", data.Key).Find(&KvStr) //删除指定数据

		//如果没数据就创建个空的
		if KvStr.Key == "" {
			db.Create(&model.KvStr{Key: data.Key, V: data.Kv})
			c.JSON(http.StatusOK, gin.H{"kv": KvStr.V})
			return
		}
		for k, v := range data.Kv {
			KvStr.V[k] = v
		}
		db.Save(&KvStr)
		c.JSON(http.StatusOK, gin.H{"kv": KvStr.V})
	})
}
