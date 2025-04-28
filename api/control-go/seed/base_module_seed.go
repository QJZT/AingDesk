package seed

import (
	"control-go/model"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// SeedBaseModule 填充 BaseModule 和 ModuleConfig 数据
func SeedBaseModule(db *gorm.DB) error {
	// 检查是否已经存在 ModuleConfig 数据
	var moduleConfigCount int64
	db.Model(&model.ModuleConfig{}).Count(&moduleConfigCount)
	if moduleConfigCount > 0 {
		fmt.Println("ModuleConfig 数据已存在，跳过填充")
		return nil
	}

	// 定义 BaseModule 假数据（3 份，按 OrderNum 排序）
	baseModules := []model.BaseModule{
		{
			ModuleConfigID:    1, // 稍后会更新为实际的 ModuleConfig ID
			OrderNum:          1,
			ModuleName:        "欢迎模块",
			IntervalTimeStart: 20,
			IntervalTimeEnd:   40,
			TriggerConditions: model.JSONSlice{model.TriggerExecuteLoop, model.TriggerLike, model.TriggerBarrageComment},
			ReadStep:          "random",
			ScriptContent:     "欢迎大家来到直播间，今天有精彩活动！{现在星期}{送礼物用户名}",
			IsModelRewrite:    true,
			RewriteFrequency:  60,
		},
		{
			ModuleConfigID:    1,
			OrderNum:          2,
			ModuleName:        "规则提醒模块",
			IntervalTimeStart: 50,
			IntervalTimeEnd:   70,
			TriggerConditions: model.JSONSlice{model.TriggerSendGift, model.TriggerEnterLiveRoom, model.TriggerWarningTip},
			ReadStep:          "sequential",
			ScriptContent:     "请大家注意直播间规则，感谢配合！",
			IsModelRewrite:    false,
			RewriteFrequency:  0,
		},
		{
			ModuleConfigID:    1,
			OrderNum:          3,
			ModuleName:        "福利模块",
			IntervalTimeStart: 30,
			IntervalTimeEnd:   60,
			TriggerConditions: model.JSONSlice{model.TriggerSendGift, model.TriggerLike, model.TriggerBarrageComment},
			ReadStep:          "random",
			ScriptContent:     "直播间福利来袭，快来参与吧！{现在时间}",
			IsModelRewrite:    true,
			RewriteFrequency:  120,
		},
	}

	// 创建 ModuleConfig（为 user_id=1 创建一个 ModuleConfig）
	firstModule := baseModules[0]
	moduleConfigName := "基础模块-" + firstModule.ScriptContent
	if len(moduleConfigName) > 20 {
		moduleConfigName = moduleConfigName[:20]
	}
	moduleConfigDescription := fmt.Sprintf("%s (触发条件: %s)", firstModule.ScriptContent, strings.Join(firstModule.TriggerConditions, ", "))
	if len(moduleConfigDescription) > 255 {
		moduleConfigDescription = moduleConfigDescription[:255]
	}

	moduleConfig := model.ModuleConfig{
		UserID:      1, // 关联到用户 ID 为 1
		Name:        moduleConfigName,
		Description: moduleConfigDescription,
	}

	// 保存 ModuleConfig
	if err := db.Create(&moduleConfig).Error; err != nil {
		return fmt.Errorf("创建 ModuleConfig 失败: %v", err)
	}

	// 更新 baseModules 的 ModuleConfigID
	for i := range baseModules {
		baseModules[i].ModuleConfigID = moduleConfig.ID
	}

	// 保存 BaseModule 数据
	for _, baseModule := range baseModules {
		if err := db.Create(&baseModule).Error; err != nil {
			return fmt.Errorf("创建 BaseModule 失败: %v", err)
		}
	}

	fmt.Println("BaseModule 和 ModuleConfig 数据填充成功")
	return nil
}
