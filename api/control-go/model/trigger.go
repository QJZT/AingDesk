package model

import (
	"gorm.io/gorm"
)

type Trigger struct {
	gorm.Model
	TriggerType string // 类型(基础模块/音频模块)

	UseLlm         bool  // 是否使用LLM改写
	UseLlmRate     int64 // 改写频率(秒数)
	UseLlmRateNext int64 // 下次改写时间

	Conditions string // 触发条件(JSON格式存储多选条件)
	IsActive   bool   // 是否激活使用

	LastTriggeredUinx int64  // 上次触发时间
	NextTriggerUinx   int64  // 下次计划触发时间
	Content           string // 全部内容
	ContentRand       bool   // 内容是否随机顺序
	NextContent       string // 触发内容
}
