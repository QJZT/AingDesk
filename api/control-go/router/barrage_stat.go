package router

import (
	"control-go/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupBarrageStatRoutes 注册弹幕统计相关路由
func SetupBarrageStatRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/save_message_stats", func(c *gin.Context) {
		var req struct {
			LiveUrl      string `json:"live_url"`
			CommentEvent uint   `json:"comment_event"`
			GiftEvent    uint   `json:"gift_event"`
			JoinEvent    uint   `json:"join_event"`
			ShareEvent   uint   `json:"share_event"`
			FollowEvent  uint   `json:"follow_event"`
			LikeEvent    uint   `json:"like_event"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		now := time.Now()
		minute := now.Minute() / 5 * 5 // 向下取整到最近的5分钟
		timeSlot := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minute, 0, 0, now.Location())
		liveTime := timeSlot.Format("2006-01-02 15:04")

		// 查找上一条记录，用于计算增量
		var prevStat model.BarrageStat
		err := db.Where("live_url = ?", req.LiveUrl).
			Where("live_date < ?", liveTime).
			Order("live_date DESC").
			Limit(1).
			First(&prevStat).Error

		// 计算增量值
		incrementValues := req
		if err == nil {
			// 如果找到上一条记录，计算增量
			incrementValues.CommentEvent = req.CommentEvent - prevStat.CommentEvent
			incrementValues.GiftEvent = req.GiftEvent - prevStat.GiftEvent
			incrementValues.JoinEvent = req.JoinEvent - prevStat.JoinEvent
			incrementValues.ShareEvent = req.ShareEvent - prevStat.ShareEvent
			incrementValues.FollowEvent = req.FollowEvent - prevStat.FollowEvent
			incrementValues.LikeEvent = req.LikeEvent - prevStat.LikeEvent

			// 防止出现负数（可能是因为重启直播等原因导致计数器重置）
			if incrementValues.CommentEvent > req.CommentEvent {
				incrementValues.CommentEvent = req.CommentEvent
			}
			if incrementValues.GiftEvent > req.GiftEvent {
				incrementValues.GiftEvent = req.GiftEvent
			}
			if incrementValues.JoinEvent > req.JoinEvent {
				incrementValues.JoinEvent = req.JoinEvent
			}
			if incrementValues.ShareEvent > req.ShareEvent {
				incrementValues.ShareEvent = req.ShareEvent
			}
			if incrementValues.FollowEvent > req.FollowEvent {
				incrementValues.FollowEvent = req.FollowEvent
			}
			if incrementValues.LikeEvent > req.LikeEvent {
				incrementValues.LikeEvent = req.LikeEvent
			}
		}

		var stat model.BarrageStat
		// 查找是否已存在该场次
		if err := db.Where("live_url = ? and live_date = ?", req.LiveUrl, liveTime).First(&stat).Error; err == nil {
			// 已存在，更新
			stat.CommentEvent = incrementValues.CommentEvent
			stat.GiftEvent = incrementValues.GiftEvent
			stat.JoinEvent = incrementValues.JoinEvent
			stat.ShareEvent = incrementValues.ShareEvent
			stat.FollowEvent = incrementValues.FollowEvent
			stat.LikeEvent = incrementValues.LikeEvent
			db.Save(&stat)
		} else {
			// 不存在，插入
			stat = model.BarrageStat{
				LiveUrl:      req.LiveUrl,
				LiveDate:     liveTime,
				CommentEvent: incrementValues.CommentEvent,
				GiftEvent:    incrementValues.GiftEvent,
				JoinEvent:    incrementValues.JoinEvent,
				ShareEvent:   incrementValues.ShareEvent,
				FollowEvent:  incrementValues.FollowEvent,
				LikeEvent:    incrementValues.LikeEvent,
			}
			db.Create(&stat)
		}
		c.JSON(http.StatusOK, gin.H{"message": "统计数据保存成功"})
	})

	// 获取url的统计数据
	r.GET("/get_stats", func(c *gin.Context) {
		// 获取查询参数
		liveUrl := c.Query("live_url")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		// 构建查询
		query := db.Model(&model.BarrageStat{})

		// 添加筛选条件
		if liveUrl != "" {
			query = query.Where("live_url = ?", liveUrl)
		}

		// 如果有日期范围，添加日期筛选
		if startDate != "" && endDate != "" {
			query = query.Where("live_date >= ? AND live_date <= ?", startDate, endDate+" 23:59")
		} else if startDate != "" {
			query = query.Where("live_date >= ?", startDate)
		} else if endDate != "" {
			query = query.Where("live_date <= ?", endDate+" 23:59")
		}

		// 按时间排序
		query = query.Order("live_date")

		// 执行查询
		var stats []model.BarrageStat
		if err := query.Find(&stats).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, stats)
	})

	// 获取所有直播链接列表
	r.GET("/get_live_urls", func(c *gin.Context) {
		var urls []string
		if err := db.Model(&model.BarrageStat{}).Distinct().Pluck("live_url", &urls).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取直播链接失败: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, urls)
	})
}
