package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// 触发条件常量
const (
	TriggerSceneLoop     = "SceneLoop"     // 控场循环
	TriggerIntervalLoop  = "IntervalLoop"  // 间隔循环：按照指定间隔时间循环执行
	TriggerBarrageComment = "BarrageComment" // 弹幕评论：是否自动发送弹幕评论
	TriggerSendGift      = "SendGift"      // 送礼物：是否自动发送礼物
	TriggerLike          = "Like"          // 点赞：是否自动点赞
	TriggerEnterLiveRoom = "EnterLiveRoom" // 进入直播间：是否自动进入直播间
	TriggerShareRoom     = "ShareRoom"     // 分享直播间
	TriggerFollowRoom    = "FollowRoom"    // 关注直播间
)

// ValidTriggerConditions 返回所有有效的触发条件
func ValidTriggerConditions() []string {
	return []string{
		TriggerSceneLoop,
		TriggerIntervalLoop,
		TriggerBarrageComment,
		TriggerSendGift,
		TriggerLike,
		TriggerEnterLiveRoom,
		TriggerShareRoom,
		TriggerFollowRoom,
	}
}

// JSONSlice 是一个自定义类型，用于处理 JSON 序列化的字符串切片
type JSONSlice []string

// Value 实现 driver.Valuer 接口，将 JSONSlice 序列化为 JSON 字符串
func (j JSONSlice) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口，从数据库字段反序列化为 JSONSlice
func (j *JSONSlice) Scan(value interface{}) error {
	if value == nil {
		*j = JSONSlice{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("unsupported type for JSONSlice")
	}
	return json.Unmarshal(bytes, j)
}

// 模块类型常量
const (
	ModuleTypeBase  = "base"  // 基础模块
	ModuleTypeAudio = "audio" // 音频模块
)

// BaseModule 基础模块模型
type BaseModule struct {
	gorm.Model
	ModuleType        string    `gorm:"not null" json:"module_type"`         // 模块类型 区分用途
	OrderNum          int       `gorm:"index" json:"order_num"`              // 排序编号，值越小越靠前
	ModuleName        string    `json:"module_name"`                         // 模块名称
	IntervalTimeStart int       `json:"interval_time_start"`                 // 间隔时间起始（秒）
	IntervalTimeEnd   int       `json:"interval_time_end"`                   // 间隔时间结束（秒）
	TriggerConditions JSONSlice `gorm:"type:text" json:"trigger_conditions"` // 触发条件（存储为 JSON 字符串）
	ReadStep          string    `json:"read_step"`                           // 读取步骤（random或sequential）
	ScriptContent     JSONSlice `gorm:"type:text" json:"script_content"`     // 话术文案列表
	IsModelRewrite    bool      `json:"isModel_rewrite"`                     // 是否模型改写
	RewriteFrequency  int       `json:"rewrite_frequency"`                   // 改写频率（秒）
	AudioName         string    `gorm:"type:text" json:"audio_name"`         // 音频文件名称，仅音频模块使用
	AudioPath         string    `gorm:"type:text" json:"audio_path"`         // 音频文件路径，仅音频模块使用
}
