package model

import "gorm.io/gorm"

type AudioBlock struct {
	gorm.Model
	Content   string // 内容
	File      string // 文件路径
	TriggerId int64  // 触发ID
	Priority  int64  `gorm:"autoIncrement;unique:false"` // 自增 优先级 1=最优先
	WinId     string // 窗口ID
	Use       bool   // 是否 已使用
}
