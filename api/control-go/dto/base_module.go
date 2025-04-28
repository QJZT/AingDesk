package dto

import "control-go/model"

// BaseModuleResponse 是返回给前端的 DTO，排除 ID、CreatedAt、UpdatedAt、DeletedAt
type BaseModuleResponse struct {
	ModuleConfigID    uint     `json:"ModuleConfigID"`
	OrderNum          int      `json:"OrderNum"`
	ModuleName        string   `json:"ModuleName"`
	IntervalTimeStart int      `json:"IntervalTimeStart"`
	IntervalTimeEnd   int      `json:"IntervalTimeEnd"`
	TriggerConditions []string `json:"TriggerConditions"`
	ReadStep          string   `json:"ReadStep"`
	ScriptContent     string   `json:"ScriptContent"`
	IsModelRewrite    bool     `json:"IsModelRewrite"`
	RewriteFrequency  int      `json:"RewriteFrequency"`
}

// ToBaseModuleResponse 转换 BaseModule 为 BaseModuleResponse
func ToBaseModuleResponse(module model.BaseModule) BaseModuleResponse {
	return BaseModuleResponse{
		ModuleConfigID:    module.ModuleConfigID,
		OrderNum:          module.OrderNum,
		ModuleName:        module.ModuleName,
		IntervalTimeStart: module.IntervalTimeStart,
		IntervalTimeEnd:   module.IntervalTimeEnd,
		TriggerConditions: module.TriggerConditions,
		ReadStep:          module.ReadStep,
		ScriptContent:     module.ScriptContent,
		IsModelRewrite:    module.IsModelRewrite,
		RewriteFrequency:  module.RewriteFrequency,
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
