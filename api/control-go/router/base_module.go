package router

import (
	"control-go/dto"
	"control-go/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupBaseModuleRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/trigger-conditions", func(c *gin.Context) {
		c.JSON(200, model.ValidTriggerConditions())
	})

	// 获取指定用户的 ModuleConfig 及其关联的 BaseModules
	r.GET("/module-configs", func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(400, gin.H{"error": "user_id parameter is required"})
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID < 1 {
			c.JSON(400, gin.H{"error": "user_id must be a positive integer"})
			return
		}

		var moduleConfigs []model.ModuleConfig
		if err := db.Where("user_id = ?", userID).
			Preload("BaseModules", func(db *gorm.DB) *gorm.DB {
				return db.Order("order_num")
			}).
			Find(&moduleConfigs).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to query module configs"})
			return
		}
		c.JSON(200, moduleConfigs)
	})

	// 获取指定用户的 BaseModules
	r.GET("/base-modules", func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(400, gin.H{"error": "user_id parameter is required"})
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID < 1 {
			c.JSON(400, gin.H{"error": "user_id must be a positive integer"})
			return
		}

		var baseModules []model.BaseModule
		if err := db.Joins("JOIN module_configs ON module_configs.id = base_modules.module_config_id").
			Where("module_configs.user_id = ?", userID).
			Order("base_modules.order_num").
			Find(&baseModules).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to query base modules"})
			return
		}
		c.JSON(200, dto.ToBaseModuleResponseSlice(baseModules))
	})

	// 创建 BaseModule
	r.POST("/base-modules", func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(400, gin.H{"error": "user_id parameter is required"})
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID < 1 {
			c.JSON(400, gin.H{"error": "user_id must be a positive integer"})
			return
		}

		var baseModule model.BaseModule
		if err := c.ShouldBindJSON(&baseModule); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if baseModule.IntervalTimeStart >= baseModule.IntervalTimeEnd {
			c.JSON(400, gin.H{"error": "IntervalTimeStart must be less than IntervalTimeEnd"})
			return
		}

		validConditions := make(map[string]bool)
		for _, cond := range model.ValidTriggerConditions() {
			validConditions[cond] = true
		}
		for _, cond := range baseModule.TriggerConditions {
			if !validConditions[cond] {
				c.JSON(400, gin.H{"error": "invalid TriggerCondition: " + cond})
				return
			}
		}

		var moduleConfig model.ModuleConfig
		if err := db.Where("id = ? AND user_id = ?", baseModule.ModuleConfigID, userID).First(&moduleConfig).Error; err != nil {
			c.JSON(400, gin.H{"error": "invalid ModuleConfigID or not owned by this user"})
			return
		}

		if err := db.Create(&baseModule).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// 返回完整 BaseModule，包含 ID（前端需要 ID 用于编辑和删除）
		c.JSON(200, baseModule)
	})

	// 更新 BaseModule
	r.PUT("/base-modules/:id", func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(400, gin.H{"error": "user_id parameter is required"})
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID < 1 {
			c.JSON(400, gin.H{"error": "user_id must be a positive integer"})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			c.JSON(400, gin.H{"error": "id must be a positive integer"})
			return
		}

		var baseModule model.BaseModule
		if err := db.Joins("JOIN module_configs ON module_configs.id = base_modules.module_config_id").
			Where("module_configs.user_id = ? AND base_modules.id = ?", userID, id).
			First(&baseModule).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "base module not found or not owned by this user"})
			} else {
				c.JSON(500, gin.H{"error": "failed to query base module"})
			}
			return
		}

		var updateModule model.BaseModule
		if err := c.ShouldBindJSON(&updateModule); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 更新字段
		baseModule.ModuleConfigID = updateModule.ModuleConfigID
		baseModule.OrderNum = updateModule.OrderNum
		baseModule.ModuleName = updateModule.ModuleName
		baseModule.IntervalTimeStart = updateModule.IntervalTimeStart
		baseModule.IntervalTimeEnd = updateModule.IntervalTimeEnd
		baseModule.TriggerConditions = updateModule.TriggerConditions
		baseModule.ReadStep = updateModule.ReadStep
		baseModule.ScriptContent = updateModule.ScriptContent
		baseModule.IsModelRewrite = updateModule.IsModelRewrite
		baseModule.RewriteFrequency = updateModule.RewriteFrequency

		if err := db.Save(&baseModule).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to update base module"})
			return
		}

		// 返回完整 BaseModule，包含 ID
		c.JSON(200, baseModule)
	})

	// 删除 BaseModule
	r.DELETE("/base-modules/:id", func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(400, gin.H{"error": "user_id parameter is required"})
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID < 1 {
			c.JSON(400, gin.H{"error": "user_id must be a positive integer"})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			c.JSON(400, gin.H{"error": "id must be a positive integer"})
			return
		}

		var baseModule model.BaseModule
		if err := db.Joins("JOIN module_configs ON module_configs.id = base_modules.module_config_id").
			Where("module_configs.user_id = ? AND base_modules.id = ?", userID, id).
			First(&baseModule).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "base module not found or not owned by this user"})
			} else {
				c.JSON(500, gin.H{"error": "failed to query base module"})
			}
			return
		}

		if err := db.Delete(&baseModule).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to delete base module"})
			return
		}

		c.JSON(200, gin.H{"message": "base module deleted successfully"})
	})

	// 获取指定模板数据
	r.GET("/template-base-module", func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(400, gin.H{"error": "user_id parameter is required"})
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID < 1 {
			c.JSON(400, gin.H{"error": "user_id must be a positive integer"})
			return
		}

		idStr := c.Query("id")
		if idStr == "" {
			c.JSON(400, gin.H{"error": "id parameter is required"})
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			c.JSON(400, gin.H{"error": "id must be a positive integer"})
			return
		}

		var moduleConfigs []model.ModuleConfig
		if err := db.Where("user_id = ?", userID).Find(&moduleConfigs).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to query module configs"})
			return
		}

		if len(moduleConfigs) == 0 {
			c.JSON(404, gin.H{"error": "no module configs found for this user"})
			return
		}

		var baseModules []model.BaseModule
		moduleConfigIDs := make([]uint, len(moduleConfigs))
		for i, mc := range moduleConfigs {
			moduleConfigIDs[i] = mc.ID
		}
		if err := db.Where("module_config_id IN ?", moduleConfigIDs).
			Order("order_num").
			Find(&baseModules).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to query base modules"})
			return
		}

		if id > len(baseModules) {
			c.JSON(404, gin.H{"error": "template not found"})
			return
		}

		c.JSON(200, dto.ToBaseModuleResponse(baseModules[id-1]))
	})
}
