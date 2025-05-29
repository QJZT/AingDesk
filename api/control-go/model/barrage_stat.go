package model

import (
	"gorm.io/gorm"
)

// BarrageStat 用于统计每种消息类型的数量
// BarrageStat 用于统计每种消息类型的数量
type BarrageStat struct {
	gorm.Model
	LiveDate     string `json:"live_date"`     // 直播时间，格式如 2025-05-22
	LiveUrl      string `json:"live_url"`      // 直播链接
	CommentEvent uint   `json:"comment_event"` // 评论消息数量
	GiftEvent    uint   `json:"gift_event"`    // 礼物消息数量
	JoinEvent    uint   `json:"join_event"`    // 进入直播间消息数量
	ShareEvent   uint   `json:"share_event"`   // 分享消息数量
	FollowEvent  uint   `json:"follow_event"`  // 关注消息数量
	LikeEvent    uint   `json:"like_event"`    // 点赞消息数量
}
