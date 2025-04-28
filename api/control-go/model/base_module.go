package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// 触发条件常量
const (
	TriggerExecuteLoop    = "ExecuteLoop"    // 执行循环：是否自动循环执行任务
	TriggerBarrageComment = "BarrageComment" // 弹幕评论：是否自动发送弹幕评论
	TriggerSendGift       = "SendGift"       // 送礼物：是否自动发送礼物
	TriggerLike           = "Like"           // 点赞：是否自动点赞
	TriggerEnterLiveRoom  = "EnterLiveRoom"  // 进入直播间：是否自动进入直播间
	TriggerWarningTip     = "WarningTip"     // 警告提示：是否显示警告提示
)

// ValidTriggerConditions 返回所有有效的触发条件
func ValidTriggerConditions() []string {
	return []string{
		TriggerExecuteLoop,
		TriggerBarrageComment,
		TriggerSendGift,
		TriggerLike,
		TriggerEnterLiveRoom,
		TriggerWarningTip,
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

// ModuleConfig 模块配置模型
type ModuleConfig struct {
	gorm.Model
	UserID      uint         `gorm:"index"`    // 用户 ID，绑定到用户
	Name        string       `gorm:"not null"` // 配置名称
	Description string       // 配置描述
	BaseModules []BaseModule `gorm:"foreignKey:ModuleConfigID"` // 一对多关系：一个 ModuleConfig 有多个 BaseModule
}

// BaseModule 基础模块模型
type BaseModule struct {
	gorm.Model
	ModuleConfigID    uint      `gorm:"index"` // 外键，关联 ModuleConfig
	OrderNum          int       `gorm:"index"` // 排序编号，值越小越靠前
	ModuleName        string    // 模块名称
	IntervalTimeStart int       // 间隔时间起始（秒）
	IntervalTimeEnd   int       // 间隔时间结束（秒）
	TriggerConditions JSONSlice `gorm:"type:text"` // 触发条件（存储为 JSON 字符串）
	ReadStep          string    // 读取步骤（random或sequential）
	ScriptContent     string    // 话术文案
	IsModelRewrite    bool      // 是否模型改写
	RewriteFrequency  int       // 改写频率（秒）
}
