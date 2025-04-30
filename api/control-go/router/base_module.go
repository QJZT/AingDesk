package router

import (
	"control-go/dto"
	"control-go/model"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupBaseModuleRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/trigger-conditions", func(c *gin.Context) {
		conditions := model.ValidTriggerConditions()
		fmt.Printf("GET /trigger-conditions 返回数据: %v\n", conditions)
		c.JSON(200, conditions)
	})

	// 获取所有基础模块
	r.GET("/base-modules", func(c *gin.Context) {
		var baseModules []model.BaseModule
		if err := db.Order("order_num").Find(&baseModules).Error; err != nil {
			fmt.Printf("GET /base-modules 失败，错误: %v\n", err)
			c.JSON(500, gin.H{"error": "failed to query base modules"})
			return
		}
		response := dto.ToBaseModuleResponseSlice(baseModules)
		fmt.Printf("GET /base-modules 返回数据: %v\n", response)
		c.JSON(200, response)
	})

	// 创建基础模块
	r.POST("/base-modules", func(c *gin.Context) {
		var baseModule model.BaseModule
		if err := c.ShouldBindJSON(&baseModule); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 验证模块类型
		if baseModule.ModuleType != model.ModuleTypeBase && baseModule.ModuleType != model.ModuleTypeAudio {
			c.JSON(400, gin.H{"error": "invalid module type"})
			return
		}

		// 如果是音频模块，验证音频路径
		if baseModule.ModuleType == model.ModuleTypeAudio && baseModule.AudioPath == "" {
			c.JSON(400, gin.H{"error": "audio module requires audio path"})
			return
		}

		if baseModule.IntervalTimeStart >= baseModule.IntervalTimeEnd {
			fmt.Println("POST /base-modules 失败，IntervalTimeStart 必须小于 IntervalTimeEnd")
			c.JSON(400, gin.H{"error": "IntervalTimeStart must be less than IntervalTimeEnd"})
			return
		}

		validConditions := make(map[string]bool)
		for _, cond := range model.ValidTriggerConditions() {
			validConditions[cond] = true
		}
		for _, cond := range baseModule.TriggerConditions {
			if !validConditions[cond] {
				fmt.Printf("POST /base-modules 失败，无效的 TriggerCondition: %s\n", cond)
				c.JSON(400, gin.H{"error": "invalid TriggerCondition: " + cond})
				return
			}
		}

		// 创建基础模块部分的返回
		if err := db.Create(&baseModule).Error; err != nil {
			fmt.Printf("POST /base-modules 失败，数据库创建错误: %v\n", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		response := dto.ToBaseModuleResponse(baseModule)
		fmt.Printf("POST /base-modules 返回数据: %v\n", response)
		c.JSON(200, response)
	})

	// 更新基础模块
	r.PUT("/base-modules/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			fmt.Printf("PUT /base-modules/%d 失败，无效的 ID: %v\n", id, err)
			c.JSON(400, gin.H{"error": "id must be a positive integer"})
			return
		}

		var baseModule model.BaseModule
		if err := db.First(&baseModule, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				fmt.Printf("PUT /base-modules/%d 失败，模块未找到\n", id)
				c.JSON(404, gin.H{"error": "base module not found"})
			} else {
				fmt.Printf("PUT /base-modules/%d 失败，数据库查询错误: %v\n", id, err)
				c.JSON(500, gin.H{"error": "failed to query base module"})
			}
			return
		}

		var updateModule model.BaseModule
		if err := c.ShouldBindJSON(&updateModule); err != nil {
			fmt.Printf("PUT /base-modules/%d 失败，JSON 绑定错误: %v\n", id, err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 验证模块类型
		if updateModule.ModuleType != model.ModuleTypeBase && updateModule.ModuleType != model.ModuleTypeAudio {
			c.JSON(400, gin.H{"error": "invalid module type"})
			return
		}

		// 如果是音频模块，验证音频路径
		if updateModule.ModuleType == model.ModuleTypeAudio && updateModule.AudioPath == "" {
			c.JSON(400, gin.H{"error": "audio module requires audio path"})
			return
		}

		// 更新字段
		baseModule.ModuleType = updateModule.ModuleType // 新增：更新模块类型
		baseModule.OrderNum = updateModule.OrderNum
		baseModule.ModuleName = updateModule.ModuleName
		baseModule.IntervalTimeStart = updateModule.IntervalTimeStart
		baseModule.IntervalTimeEnd = updateModule.IntervalTimeEnd
		baseModule.TriggerConditions = updateModule.TriggerConditions
		baseModule.ReadStep = updateModule.ReadStep
		baseModule.ScriptContent = updateModule.ScriptContent
		baseModule.IsModelRewrite = updateModule.IsModelRewrite
		baseModule.RewriteFrequency = updateModule.RewriteFrequency
		baseModule.AudioPath = updateModule.AudioPath // 新增：更新音频路径

		// 更新基础模块部分的返回
		if err := db.Save(&baseModule).Error; err != nil {
			fmt.Printf("PUT /base-modules/%d 失败，数据库更新错误: %v\n", id, err)
			c.JSON(500, gin.H{"error": "failed to update base module"})
			return
		}

		response := dto.ToBaseModuleResponse(baseModule)
		fmt.Printf("PUT /base-modules/%d 返回数据: %v\n", id, response)
		c.JSON(200, response)
	})

	// 删除基础模块
	r.DELETE("/base-modules/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			fmt.Printf("DELETE /base-modules/%d 失败，无效的 ID: %v\n", id, err)
			c.JSON(400, gin.H{"error": "id must be a positive integer"})
			return
		}

		var baseModule model.BaseModule
		if err := db.First(&baseModule, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				fmt.Printf("DELETE /base-modules/%d 失败，模块未找到\n", id)
				c.JSON(404, gin.H{"error": "base module not found"})
			} else {
				fmt.Printf("DELETE /base-modules/%d 失败，数据库查询错误: %v\n", id, err)
				c.JSON(500, gin.H{"error": "failed to query base module"})
			}
			return
		}

		if err := db.Delete(&baseModule).Error; err != nil {
			fmt.Printf("DELETE /base-modules/%d 失败，数据库删除错误: %v\n", id, err)
			c.JSON(500, gin.H{"error": "failed to delete base module"})
			return
		}

		fmt.Printf("DELETE /base-modules/%d 返回数据: %v\n", id, gin.H{"message": "base module deleted successfully"})
		c.JSON(200, gin.H{"message": "base module deleted successfully"})
	})
}
