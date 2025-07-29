package router

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"control-go/dto"
	"control-go/global"
	"control-go/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse 统一成功响应格式
type SuccessResponse struct {
	Message string `json:"message"`
	Path    string `json:"path"`
}

// respondWithError 统一错误响应函数
func respondWithError(c *gin.Context, status int, message string) {
	log.Error().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Int("status", status).
		Str("error", message).
		Msg("请求失败")
	c.JSON(status, ErrorResponse{Error: message})
}

// respondWithSuccess 统一成功响应函数
func respondWithSuccess(c *gin.Context, message, path string) {
	log.Info().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("message", message).
		Str("path", path).
		Msg("请求成功")
	c.JSON(200, SuccessResponse{Message: message, Path: path})
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
			Msg("返回触发条件列表")
		c.JSON(200, conditions)
	})

	// 上传音频文件路由
	r.POST("/upload-audio", func(c *gin.Context) {
		// 获取关联 ID（可选，假设前端提供）
		id := c.PostForm("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			idInt = 0 // 默认值，如果解析失败
		}

		// 获取上传文件
		file, err := c.FormFile("file")
		if err != nil {
			respondWithError(c, 400, "获取文件失败："+err.Error())
			return
		}

		// 验证文件类型
		allowedTypes := []string{
			"audio/mpeg", "audio/wav", "audio/ogg", "audio/aac", "audio/mp4",
			"audio/webm", "audio/flac", "audio/x-m4a", "audio/x-ms-wma",
			"audio/x-aiff", "audio/x-wav", "audio/x-mpegurl", "audio/x-scpls",
		}
		valid := false
		for _, t := range allowedTypes {
			if file.Header.Get("Content-Type") == t {
				valid = true
				break
			}
		}
		if !valid {
			respondWithError(c, 400, "仅支持音频文件（MP3、WAV、OGG 等）")
			return
		}

		// 创建保存路径
		saveDir := filepath.Join(global.AudioFilePath, fmt.Sprintf("%d", idInt))
		if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
			respondWithError(c, 500, "创建上传目录失败："+err.Error())
			return
		}

		// 构造文件路径
		filePath := filepath.Join(saveDir, file.Filename)
		if _, err := os.Stat(filePath); err == nil {
			respondWithError(c, 400, "文件已存在："+file.Filename)
			return
		}

		// 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			respondWithError(c, 500, "保存文件失败："+err.Error())
			return
		}

		// 打印保存地址到日志
		log.Info().Str("saved_file_path", filePath).Msg("音频文件保存成功")

		// 返回成功响应
		respondWithSuccess(c, "音频文件上传成功", filePath)
	})

	// 基础模块路由（保持原有的逻辑不变）
	baseModules := r.Group("/base-modules")
	{
		// 获取所有基础模块
		baseModules.GET("", func(c *gin.Context) {
			var baseModules []model.BaseModule
			if err := db.Order("order_num").Find(&baseModules).Error; err != nil {
				respondWithError(c, 500, "查询基础模块失败")
				return
			}

			response := dto.ToBaseModuleResponseSlice(baseModules)
			// log.Info().
			// 	Str("method", "GET").
			// 	Str("path", "/base-modules").
			// 	Interface("response", response).
			// 	Msg("返回基础模块列表")
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			c.JSON(200, response)
		})

		// 创建基础模块（保持原有的逻辑不变）
		baseModules.POST("", func(c *gin.Context) {
			var baseModule model.BaseModule
			if err := c.ShouldBindJSON(&baseModule); err != nil {
				fmt.Println("错误:", err)
				respondWithError(c, 400, "绑定 JSON 失败："+err.Error())
				return
			}
			log.Info().Interface("绑定模块", baseModule).Msg("成功将 JSON 绑定到 BaseModule")

			// 验证模块类型（必填）
			if baseModule.ModuleType == "" {
				respondWithError(c, 400, "ModuleType 不能为空")
				return
			}
			if baseModule.ModuleType != model.ModuleTypeBase && baseModule.ModuleType != model.ModuleTypeAudio {
				log.Error().Str("module_type", baseModule.ModuleType).Msg("检测到无效的模块类型")
				respondWithError(c, 400, "无效的模块类型")
				return
			}

			// 如果是音频模块，验证音频路径
			if baseModule.ModuleType == model.ModuleTypeAudio && baseModule.AudioPath == "" {
				respondWithError(c, 400, "音频模块需要音频路径")
				return
			}

			// 验证 TriggerConditions 不为空
			if len(baseModule.TriggerConditions) == 0 {
				respondWithError(c, 400, "TriggerConditions 不能为空")
				return
			}

			// 验证 TriggerConditions 有效性
			validConditions := make(map[string]bool)
			for _, cond := range model.ValidTriggerConditions() {
				validConditions[cond] = true
			}
			for _, cond := range baseModule.TriggerConditions {
				if !validConditions[cond] {
					respondWithError(c, 400, "无效的 TriggerCondition: "+cond)
					return
				}
			}

			// 验证 ReadStep
			validReadSteps := map[string]bool{"random": true, "sequential": true}
			if baseModule.ReadStep != "" && !validReadSteps[baseModule.ReadStep] {
				respondWithError(c, 400, "无效的 ReadStep: "+baseModule.ReadStep)
				return
			}

			// 验证 RewriteFrequency（仅在 IsModelRewrite 为 true 时要求）
			if baseModule.IsModelRewrite && baseModule.RewriteFrequency <= 0 {
				respondWithError(c, 400, "当 IsModelRewrite 为 true 时，RewriteFrequency 必须大于 0")
				return
			}

			// 使用事务创建模块 改
			if err := db.Create(&baseModule).Error; err != nil {
				log.Error().Err(err).Interface("module", baseModule).Msg("创建基础模块失败")
				respondWithError(c, 500, "创建基础模块失败："+err.Error())
				return
			}
			log.Info().Interface("创建的模块", baseModule).Msg("基础模块创建成功")

			response := dto.ToBaseModuleResponse(baseModule)
			log.Info().
				Str("method", "POST").
				Str("path", "/base-modules").
				Interface("response", response).
				Msg("基础模块创建成功")
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			c.JSON(200, response)
		})

		// 更新基础模块
		baseModules.PUT(":id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil || id < 1 {
				respondWithError(c, 400, "id 必须为正整数")
				return
			}

			var baseModule model.BaseModule
			if err := db.First(&baseModule, id).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					respondWithError(c, 404, "基础模块未找到")
				} else {
					respondWithError(c, 500, "查询基础模块失败")
				}
				return
			}

			var updateModule model.BaseModule
			if err := c.ShouldBindJSON(&updateModule); err != nil {
				log.Error().Err(err).Msg("绑定 JSON 失败")
				respondWithError(c, 400, "绑定 JSON 失败："+err.Error())
				return
			}
			log.Info().Interface("更新模块", updateModule).Msg("成功将 JSON 绑定到 BaseModule")

			// 验证模块类型（必填）
			if updateModule.ModuleType == "" {
				respondWithError(c, 400, "ModuleType 不能为空")
				return
			}
			if updateModule.ModuleType != model.ModuleTypeBase && updateModule.ModuleType != model.ModuleTypeAudio {
				log.Error().Str("module_type", updateModule.ModuleType).Msg("检测到无效的模块类型")
				respondWithError(c, 400, "无效的模块类型")
				return
			}

			// 如果是音频模块，验证音频路径
			if updateModule.ModuleType == model.ModuleTypeAudio && updateModule.AudioPath == "" {
				respondWithError(c, 400, "音频模块需要音频路径")
				return
			}

			// 验证 IntervalTimeStart 和 IntervalTimeEnd
			//if updateModule.IntervalTimeStart >= updateModule.IntervalTimeEnd {
			//respondWithError(c, 400, "IntervalTimeStart 必须小于 IntervalTimeEnd")
			//return
			//}

			// 验证 TriggerConditions 不为空
			if len(updateModule.TriggerConditions) == 0 {
				respondWithError(c, 400, "TriggerConditions 不能为空")
				return
			}

			// 验证 TriggerConditions 有效性
			validConditions := make(map[string]bool)
			for _, cond := range model.ValidTriggerConditions() {
				validConditions[cond] = true
			}
			for _, cond := range updateModule.TriggerConditions {
				if !validConditions[cond] {
					respondWithError(c, 400, "无效的 TriggerCondition: "+cond)
					return
				}
			}

			// 验证 ModuleName 非空
			if updateModule.ModuleName == "" {
				respondWithError(c, 400, "ModuleName 不能为空")
				return
			}

			// 验证 ReadStep
			validReadSteps := map[string]bool{"random": true, "sequential": true}
			if updateModule.ReadStep != "" && !validReadSteps[updateModule.ReadStep] {
				respondWithError(c, 400, "无效的 ReadStep: "+updateModule.ReadStep)
				return
			}

			// 验证 RewriteFrequency（仅在 IsModelRewrite 为 true 时要求）
			if updateModule.IsModelRewrite && updateModule.RewriteFrequency <= 0 {
				respondWithError(c, 400, "当 IsModelRewrite 为 true 时，RewriteFrequency 必须大于 0")
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

			// 使用事务更新模块 改
			// 直接更新模块
			if err := db.Save(&baseModule).Error; err != nil {
				log.Error().Err(err).Interface("module", baseModule).Msg("更新基础模块失败")
				respondWithError(c, 500, "更新基础模块失败")
				return
			}
			log.Info().Interface("updated_module", baseModule).Msg("基础模块更新成功")

			response := dto.ToBaseModuleResponse(baseModule)
			log.Info().
				Str("method", "PUT").
				Str("path", "/base-modules/"+strconv.Itoa(id)).
				Interface("response", response).
				Msg("基础模块更新成功")
			c.JSON(200, response)
		})

		// 删除基础模块
		baseModules.DELETE(":id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil || id < 1 {
				respondWithError(c, 400, "id 必须为正整数")
				return
			}

			var baseModule model.BaseModule
			if err := db.First(&baseModule, id).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					respondWithError(c, 404, "基础模块未找到")
				} else {
					respondWithError(c, 500, "查询基础模块失败")
				}
				return
			}

			// 使用事务删除模块 改
			// 直接删除模块
			if err := db.Delete(&baseModule).Error; err != nil {
				log.Error().Err(err).Msg("删除基础模块失败")
				respondWithError(c, 500, "删除基础模块失败")
				return
			}
			log.Info().Msg("基础模块删除成功")

			log.Info().
				Str("method", "DELETE").
				Str("path", "/base-modules/"+strconv.Itoa(id)).
				Msg("基础模块删除成功")
			c.JSON(200, gin.H{"message": "基础模块删除成功"})
		})
	}

	// 音频模块路由
	audioModules := r.Group("/audio-modules")
	{
		// 创建音频模块
		audioModules.POST("", func(c *gin.Context) {
			// 读取原始请求体以供调试
			bodyBytes, err := c.GetRawData()
			if err != nil {
				log.Error().Err(err).Msg("读取请求体失败")
				respondWithError(c, 400, "读取请求体失败："+err.Error())
				return
			}
			log.Info().Str("原始请求体", string(bodyBytes)).Msg("接收到音频模块创建请求")

			// 重新绑定请求体（因为 GetRawData 消耗了 c.Request.Body）
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var baseModule model.BaseModule
			if err := c.ShouldBindJSON(&baseModule); err != nil {
				log.Error().Err(err).Msg("绑定 JSON 失败")
				respondWithError(c, 400, "绑定 JSON 失败："+err.Error())
				return
			}
			log.Info().Interface("绑定模块", baseModule).Msg("成功将 JSON 绑定到 BaseModule")

			// 强制设置模块类型为 audio
			baseModule.ModuleType = model.ModuleTypeAudio

			// 验证音频路径（音频模块必须有音频路径）
			if baseModule.AudioPath == "" {
				respondWithError(c, 400, "音频模块需要音频路径")
				return
			}

			// 验证 TriggerConditions 不为空
			if len(baseModule.TriggerConditions) == 0 {
				respondWithError(c, 400, "TriggerConditions 不能为空")
				return
			}

			// 验证 TriggerConditions 有效性
			validConditions := make(map[string]bool)
			for _, cond := range model.ValidTriggerConditions() {
				validConditions[cond] = true
			}
			for _, cond := range baseModule.TriggerConditions {
				if cond == "" {
					respondWithError(c, 400, "TriggerConditions 不能包含空字符串")
					return
				}
				if !validConditions[cond] {
					respondWithError(c, 400, "无效的 TriggerCondition: "+cond)
					return
				}
			}

			// 验证 IntervalTimeStart 和 IntervalTimeEnd
			if baseModule.IntervalTimeStart > baseModule.IntervalTimeEnd {
				respondWithError(c, 400, "IntervalTimeStart 必须小于或等于 IntervalTimeEnd")
				return
			}

			// 使用事务创建模块 改
			// 直接创建模块
			if err := db.Create(&baseModule).Error; err != nil {
				log.Error().Err(err).Interface("module", baseModule).Msg("创建音频模块失败")
				respondWithError(c, 500, "创建音频模块失败："+err.Error())
				return
			}
			log.Info().Interface("创建的模块", baseModule).Msg("音频模块创建成功")

			response := dto.ToBaseModuleResponse(baseModule)
			log.Info().
				Str("method", "POST").
				Str("path", "/audio-modules").
				Interface("response", response).
				Msg("音频模块创建成功")
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
			c.JSON(200, response)
		})
	}
}
