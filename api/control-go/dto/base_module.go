package dto

import "control-go/model"

// BaseModuleResponse 是返回给前端的 DTO
type BaseModuleResponse struct {
	ID                uint     `json:"id"`                  // 模块ID
	ModuleType        string   `json:"module_type"`         // 模块类型：base(基础模块)或audio(音频模块)
	OrderNum          int      `json:"order_num"`           // 排序编号，值越小越靠前
	ModuleName        string   `json:"module_name"`         // 模块名称
	IntervalTimeStart int      `json:"interval_time_start"` // 间隔时间起始（秒）
	IntervalTimeEnd   int      `json:"interval_time_end"`   // 间隔时间结束（秒）
	TriggerConditions []string `json:"trigger_conditions"`  // 触发条件列表
	ReadStep          string   `json:"read_step"`           // 读取步骤：random(随机)或sequential(顺序)
	ScriptContent     []string `json:"script_content"`      // 话术文案列表
	IsModelRewrite    bool     `json:"is_model_rewrite"`    // 是否启用模型改写
	RewriteFrequency  int      `json:"rewrite_frequency"`   // 改写频率（秒）
	AudioName         string   `json:"audio_name"`          // 音频文件名称，仅音频模块使用
	AudioPath         string   `json:"audio_path"`          // 音频文件路径，仅音频模块使用

}

// ToBaseModuleResponse 转换 BaseModule 为 BaseModuleResponse
// 该函数将数据库模型转换为前端响应DTO
// 参数:
//   - module: 数据库中的基础模块模型
// 返回:
//   - BaseModuleResponse: 转换后的前端响应对象
func ToBaseModuleResponse(module model.BaseModule) BaseModuleResponse {
	return BaseModuleResponse{
		ID:                module.ID,                // 模块ID
		ModuleType:        module.ModuleType,        // 模块类型
		OrderNum:          module.OrderNum,          // 排序编号
		ModuleName:        module.ModuleName,        // 模块名称
		IntervalTimeStart: module.IntervalTimeStart, // 间隔时间起始
		IntervalTimeEnd:   module.IntervalTimeEnd,   // 间隔时间结束
		TriggerConditions: module.TriggerConditions, // 触发条件列表
		ReadStep:          module.ReadStep,          // 读取步骤
		ScriptContent:     module.ScriptContent,     // 话术文案列表
		IsModelRewrite:    module.IsModelRewrite,    // 是否启用模型改写
		RewriteFrequency:  module.RewriteFrequency,  // 改写频率
		AudioName:         module.AudioName,         // 音频文件名称
		AudioPath:         module.AudioPath,         // 音频文件路径
	}
}

// ToBaseModuleResponseSlice 转换 BaseModule 切片为 BaseModuleResponse 切片
func ToBaseModuleResponseSlice(modules []model.BaseModule) []BaseModuleResponse {
	result := make([]BaseModuleResponse, len(modules))
	for i, module := range modules {
		result[i] = ToBaseModuleResponse(module)
	}
	return result
}
