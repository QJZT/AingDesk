package router

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"control-go/dto"
	"control-go/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Error string `json:"error"`
}

// respondWithError 统一错误响应函数
func respondWithError(c *gin.Context, status int, message string) {
	log.Error().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Int("status", status).
		Str("error", message).
		Msg("Request failed")
	c.JSON(status, ErrorResponse{Error: message})
}

// SetupBaseModuleRoutes 设置基础模块路由
func SetupBaseModuleRoutes(r *gin.Engine, db *gorm.DB) {
	// 公开路由
	r.GET("/trigger-conditions", func(c *gin.Context) {
		conditions := model.ValidTriggerConditions()
		log.Info().
			Str("method", "GET").
			Str("path", "/trigger-conditions").
			Interface("response", conditions).
			Msg("Returning trigger conditions")
		c.JSON(200, conditions)
	})

	// 基础模块路由
	baseModules := r.Group("/base-modules")
	{
		// // 处理 OPTIONS 请求以支持 CORS 预检
		// baseModules.OPTIONS("", func(c *gin.Context) {
		// 	log.Info().
		// 		Str("method", "OPTIONS").
		// 		Str("path", "/base-modules").
		// 		Msg("Received OPTIONS request for CORS preflight")
		// 	c.Header("Access-Control-Allow-Origin", "*")
		// 	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		// 	c.Header("Access-Control-Max-Age", "86400")
		// 	c.Status(204)
		// })

		// 获取所有基础模块
		baseModules.GET("", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
			defer cancel()

			var baseModules []model.BaseModule
			if err := db.WithContext(ctx).Order("order_num").Find(&baseModules).Error; err != nil {
				respondWithError(c, 500, "failed to query base modules")
				return
			}

			response := dto.ToBaseModuleResponseSlice(baseModules)
			log.Info().
				Str("method", "GET").
				Str("path", "/base-modules").
				Interface("response", response).
				Msg("Returning base modules")
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			c.JSON(200, response)
		})

		// 创建基础模块
		baseModules.POST("", func(c *gin.Context) {
			// 打印原始请求体
			// bodyBytes, _ := c.GetRawData()
			// log.Info().Str("raw_request_body", string(bodyBytes)).Msg("Raw request body received")

			var baseModule model.BaseModule
			if err := c.ShouldBindJSON(&baseModule); err != nil {
				fmt.Println("err:", err)
				// log.Error().Err(err).Interface("request_body", string(bodyBytes)).Msg("Failed to bind JSON")
				respondWithError(c, 400, err.Error())
				return
			}
			fmt.Println("1")
			log.Info().Interface("bound_module", baseModule).Msg("Successfully bound JSON to BaseModule")

			// 验证模块类型（必填）
			if baseModule.ModuleType == "" {
				respondWithError(c, 400, "ModuleType cannot be empty")
				return
			}
			fmt.Println("2")
			if baseModule.ModuleType != model.ModuleTypeBase && baseModule.ModuleType != model.ModuleTypeAudio {
				log.Error().Str("module_type", baseModule.ModuleType).Msg("Invalid module type detected")
				respondWithError(c, 400, "invalid module type")
				return
			}

			fmt.Println("3")

			// 如果是音频模块，验证音频路径
			if baseModule.ModuleType == model.ModuleTypeAudio && baseModule.AudioPath == "" {
				respondWithError(c, 400, "audio module requires audio path")
				return
			}

			// 验证 IntervalTimeStart 和 IntervalTimeEnd
			// if baseModule.IntervalTimeStart >= baseModule.IntervalTimeEnd {
			// 	respondWithError(c, 400, "时间设置有误")
			// 	return
			// }

			// 验证 TriggerConditions 不为空
			if len(baseModule.TriggerConditions) == 0 {
				respondWithError(c, 400, "TriggerConditions cannot be empty")
				return
			}

			// 验证 TriggerConditions 有效性
			validConditions := make(map[string]bool)
			for _, cond := range model.ValidTriggerConditions() {
				validConditions[cond] = true
			}
			for _, cond := range baseModule.TriggerConditions {
				if !validConditions[cond] {
					respondWithError(c, 400, "invalid TriggerCondition: "+cond)
					return
				}
			}

			// 验证 ModuleName 非空
			// if baseModule.ModuleName == "" {
			// 	respondWithError(c, 400, "ModuleName cannot be empty")
			// 	return
			// }

			// 验证 ReadStep
			validReadSteps := map[string]bool{"random": true, "sequential": true}
			if baseModule.ReadStep != "" && !validReadSteps[baseModule.ReadStep] {
				respondWithError(c, 400, "invalid ReadStep: "+baseModule.ReadStep)
				return
			}

			// 验证 RewriteFrequency（仅在 IsModelRewrite 为 true 时要求）
			if baseModule.IsModelRewrite && baseModule.RewriteFrequency <= 0 {
				respondWithError(c, 400, "RewriteFrequency must be greater than 0 when IsModelRewrite is true")
				return
			}

			// 使用事务创建模块
			tx := db.Begin()
			if err := tx.Create(&baseModule).Error; err != nil {
				log.Error().Err(err).Interface("module", baseModule).Msg("Failed to create base module")
				tx.Rollback()
				respondWithError(c, 500, err.Error())
				return
			}
			log.Info().Interface("created_module", baseModule).Msg("Base module created in transaction")

			if err := tx.Commit().Error; err != nil {
				log.Error().Err(err).Msg("Failed to commit transaction")
				tx.Rollback()
				respondWithError(c, 500, "failed to commit transaction")
				return
			}
			log.Info().Msg("Transaction committed successfully")

			response := dto.ToBaseModuleResponse(baseModule)
			log.Info().
				Str("method", "POST").
				Str("path", "/base-modules").
				Interface("response", response).
				Msg("Base module created")
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			c.JSON(200, response)
		})

		// 更新基础模块
		baseModules.PUT(":id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil || id < 1 {
				respondWithError(c, 400, "id must be a positive integer")
				return
			}

			var baseModule model.BaseModule
			if err := db.First(&baseModule, id).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					respondWithError(c, 404, "base module not found")
				} else {
					respondWithError(c, 500, "failed to query base module")
				}
				return
			}

			var updateModule model.BaseModule
			if err := c.ShouldBindJSON(&updateModule); err != nil {
				log.Error().Err(err).Msg("Failed to bind JSON")
				respondWithError(c, 400, err.Error())
				return
			}
			log.Info().Interface("update_module", updateModule).Msg("Successfully bound JSON to BaseModule")

			// 验证模块类型（必填）
			if updateModule.ModuleType == "" {
				respondWithError(c, 400, "ModuleType cannot be empty")
				return
			}
			if updateModule.ModuleType != model.ModuleTypeBase && updateModule.ModuleType != model.ModuleTypeAudio {
				log.Error().Str("module_type", updateModule.ModuleType).Msg("Invalid module type detected")
				respondWithError(c, 400, "invalid module type")
				return
			}

			// 如果是音频模块，验证音频路径
			if updateModule.ModuleType == model.ModuleTypeAudio && updateModule.AudioPath == "" {
				respondWithError(c, 400, "audio module requires audio path")
				return
			}

			// 验证 IntervalTimeStart 和 IntervalTimeEnd
			if updateModule.IntervalTimeStart >= updateModule.IntervalTimeEnd {
				respondWithError(c, 400, "IntervalTimeStart must be less than IntervalTimeEnd")
				return
			}

			// 验证 TriggerConditions 不为空
			if len(updateModule.TriggerConditions) == 0 {
				respondWithError(c, 400, "TriggerConditions cannot be empty")
				return
			}

			// 验证 TriggerConditions 有效性
			validConditions := make(map[string]bool)
			for _, cond := range model.ValidTriggerConditions() {
				validConditions[cond] = true
			}
			for _, cond := range updateModule.TriggerConditions {
				if !validConditions[cond] {
					respondWithError(c, 400, "invalid TriggerCondition: "+cond)
					return
				}
			}

			// 验证 ModuleName 非空
			if updateModule.ModuleName == "" {
				respondWithError(c, 400, "ModuleName cannot be empty")
				return
			}

			// 验证 ReadStep
			validReadSteps := map[string]bool{"random": true, "sequential": true}
			if updateModule.ReadStep != "" && !validReadSteps[updateModule.ReadStep] {
				respondWithError(c, 400, "invalid ReadStep: "+updateModule.ReadStep)
				return
			}

			// 验证 RewriteFrequency（仅在 IsModelRewrite 为 true 时要求）
			if updateModule.IsModelRewrite && updateModule.RewriteFrequency <= 0 {
				respondWithError(c, 400, "RewriteFrequency must be greater than 0 when IsModelRewrite is true")
				return
			}

			// 更新字段
			baseModule.ModuleType = updateModule.ModuleType
			baseModule.OrderNum = updateModule.OrderNum
			baseModule.ModuleName = updateModule.ModuleName
			baseModule.IntervalTimeStart = updateModule.IntervalTimeStart
			baseModule.IntervalTimeEnd = updateModule.IntervalTimeEnd
			baseModule.TriggerConditions = updateModule.TriggerConditions
			baseModule.ReadStep = updateModule.ReadStep
			baseModule.ScriptContent = updateModule.ScriptContent
			baseModule.IsModelRewrite = updateModule.IsModelRewrite
			baseModule.RewriteFrequency = updateModule.RewriteFrequency
			baseModule.AudioName = updateModule.AudioName
			baseModule.AudioPath = updateModule.AudioPath

			// 使用事务更新模块
			tx := db.Begin()
			if err := tx.Save(&baseModule).Error; err != nil {
				log.Error().Err(err).Interface("module", baseModule).Msg("Failed to update base module")
				tx.Rollback()
				respondWithError(c, 500, "failed to update base module")
				return
			}
			log.Info().Interface("updated_module", baseModule).Msg("Base module updated in transaction")

			if err := tx.Commit().Error; err != nil {
				log.Error().Err(err).Msg("Failed to commit transaction")
				tx.Rollback()
				respondWithError(c, 500, "failed to commit transaction")
				return
			}
			log.Info().Msg("Transaction committed successfully")

			response := dto.ToBaseModuleResponse(baseModule)
			log.Info().
				Str("method", "PUT").
				Str("path", "/base-modules/"+strconv.Itoa(id)).
				Interface("response", response).
				Msg("Base module updated")
			c.JSON(200, response)
		})

		// 删除基础模块
		baseModules.DELETE(":id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil || id < 1 {
				respondWithError(c, 400, "id must be a positive integer")
				return
			}

			var baseModule model.BaseModule
			if err := db.First(&baseModule, id).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					respondWithError(c, 404, "base module not found")
				} else {
					respondWithError(c, 500, "failed to query base module")
				}
				return
			}

			// 使用事务删除模块
			tx := db.Begin()
			if err := tx.Delete(&baseModule).Error; err != nil {
				log.Error().Err(err).Msg("Failed to delete base module")
				tx.Rollback()
				respondWithError(c, 500, "failed to delete base module")
				return
			}
			log.Info().Msg("Base module deleted in transaction")

			if err := tx.Commit().Error; err != nil {
				log.Error().Err(err).Msg("Failed to commit transaction")
				tx.Rollback()
				respondWithError(c, 500, "failed to commit transaction")
				return
			}
			log.Info().Msg("Transaction committed successfully")

			log.Info().
				Str("method", "DELETE").
				Str("path", "/base-modules/"+strconv.Itoa(id)).
				Msg("Base module deleted")
			c.JSON(200, gin.H{"message": "base module deleted successfully"})
		})
	}
}
