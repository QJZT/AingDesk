package model

import (
	"gorm.io/gorm"
)

type Trigger struct {
	gorm.Model
	TriggerName string // 名称
	TriggerDesc string // 描述
	TriggerType string // 类型(基础模块/音频模块)

	UseLlm     bool  // 是否使用LLM改写
	UseLlmRate int64 // 改写频率(秒数)

	Conditions string // 触发条件(JSON格式存储多选条件) 循环执行，弹幕评论，谢礼物，点赞，进入直播间，感谢关注。
	IsActive   bool   // 是否激活使用

	ContentRand bool   // 内容是否随机顺序
	Content     string // 全部内容
}
