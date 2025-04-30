package dto

import "control-go/model"

// BaseModuleResponse 是返回给前端的 DTO
type BaseModuleResponse struct {
	ID                uint     `json:"id"`
	ModuleType        string   `json:"module_type"`
	OrderNum          int      `json:"order_num"`
	ModuleName        string   `json:"module_name"`
	IntervalTimeStart int      `json:"interval_time_start"`
	IntervalTimeEnd   int      `json:"interval_time_end"`
	TriggerConditions []string `json:"trigger_conditions"`
	ReadStep          string   `json:"read_step"`
	ScriptContent     []string `json:"script_content"` // 修改为字符串数组
	IsModelRewrite    bool     `json:"is_model_rewrite"`
	RewriteFrequency  int      `json:"rewrite_frequency"`
	AudioPath         string   `json:"audio_path"`
}

// ToBaseModuleResponse 转换 BaseModule 为 BaseModuleResponse
func ToBaseModuleResponse(module model.BaseModule) BaseModuleResponse {
	return BaseModuleResponse{
		ID:                module.ID,
		ModuleType:        module.ModuleType,
		OrderNum:          module.OrderNum,
		ModuleName:        module.ModuleName,
		IntervalTimeStart: module.IntervalTimeStart,
		IntervalTimeEnd:   module.IntervalTimeEnd,
		TriggerConditions: module.TriggerConditions,
		ReadStep:          module.ReadStep,
		ScriptContent:     module.ScriptContent, // JSONSlice 类型会自动转换为 []string
		IsModelRewrite:    module.IsModelRewrite,
		RewriteFrequency:  module.RewriteFrequency,
		AudioPath:         module.AudioPath,
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
