package seed

import (
	"control-go/model"
	"fmt"

	"gorm.io/gorm"
)

func SeedBaseModule(db *gorm.DB) error {
	var count int64
	db.Model(&model.BaseModule{}).Count(&count)
	if count > 0 {
		fmt.Println("BaseModule 数据已存在，跳过填充")
		return nil
	}

	baseModules := []model.BaseModule{
		{
			ModuleType:        model.ModuleTypeBase,
			OrderNum:          0,
			ModuleName:        "基础模块", //22222
			IntervalTimeStart: 5,
			IntervalTimeEnd:   10,
			TriggerConditions: model.JSONSlice{model.TriggerExecuteLoop},
			ReadStep:          "sequential",
			ScriptContent: model.JSONSlice{"一号链接刷新去拍，马上封单下架了，已经有一个姐姐抢了五件了，喜欢的姐姐抓紧了一号链接，马上没有了，小助理，一号链接还有几个库存，有多少个拍下没有付款的，好的，姐妹们，一定要秒拍秒付不要卡库存",
				"姐姐们点击下方一号链接直接拍下，卡码往大一码拍，价格已开，下单的的姐妹直接付款，别占库存库存有限，先到先得一号链接已经开价。拼手速了姐妹，刷新一号链接。小助理，现在几点了，再过五分钟{再过五分钟}，就将一号链接下架",
				"姐妹一号链接价格已开，别犹豫徘徊，{现在时间}犹豫徘徊等于白来，马上就没有库存了。抓紧时间，秒拍秒付。刷新链接点击一号，根据尺码去拍，秒拍秒付姐姐们，能拍下就是有库存，拍不了就没货了。",
				"来姐姐们点击下方小黄车一号链接，今天给你们的价格是你意想不到的价格，一定会让你惊喜。里面有尺码介绍，脚厚脚肥卡码的拍大一码。",
				"追求高级感追求品质的姐妹一定要拍，下方小黄色一号链接，二十这个价格一定秒没，手速要快，别纠结颜色，你能抢到再说。小助理看下还有几个库存小黄车的一号链接二十个名额，追求高货的姐妹一定别错过这款式，它真的是让你收到货有无限惊喜，它会告诉你什么叫看一眼就心动",
				"来话不多说，下方小黄车一号链接，如果你拍回家不止这个三五百你给我退回来。今天的这个价格给到二十个名额，是做品宣的，明天我们就要涨价了。",
			},
			IsModelRewrite:   false,
			RewriteFrequency: 0,
		},
		{
			ModuleType:        model.ModuleTypeBase,
			OrderNum:          1,
			ModuleName:        "欢迎模块",
			IntervalTimeStart: 20,
			IntervalTimeEnd:   40,
			TriggerConditions: model.JSONSlice{model.TriggerEnterLiveRoom},
			ReadStep:          "random",
			ScriptContent:     model.JSONSlice{"欢迎{用户名称}的到来，您你的到来让直播间蓬革", "{用户名称}来了呀，欢迎你，谢谢你来我直播间", "日照香炉生紫烟，感谢{用户名称}来到直播间"},
			IsModelRewrite:    true,
			RewriteFrequency:  60,
		},
		{
			ModuleType:        model.ModuleTypeAudio,
			OrderNum:          2,
			ModuleName:        "欢迎模块",
			IntervalTimeStart: 30,
			IntervalTimeEnd:   50,
			TriggerConditions: model.JSONSlice{model.TriggerEnterLiveRoom},
			ReadStep:          "sequential",
			ScriptContent:     model.JSONSlice{"{用户名称}我的老板，你来了呀", "欢迎{用户名称}老板进来捧场", "{用户名称}老板，你好呀", "掌声欢迎欢迎{用户名称}老板进入直播间，祝您财运滚滚来", "撒花欢迎{用户名称}老总进入直播间，愿你事业蒸蒸日上", "欢迎{用户名称}老板进入直播间，谢谢！祝你发大财。"},
			IsModelRewrite:    false,
			RewriteFrequency:  0,
		},
		{
			ModuleType:        model.ModuleTypeAudio,
			OrderNum:          3,
			ModuleName:        "点赞模块",
			IntervalTimeStart: 30,
			IntervalTimeEnd:   50,
			TriggerConditions: model.JSONSlice{model.TriggerEnterLiveRoom},
			ReadStep:          "sequential",
			ScriptContent:     model.JSONSlice{"感谢{用户名称}的点赞，谢谢老板", "感谢{用户名称}老板的点赞，太感谢了", "感谢{用户名称}的点赞，谢谢", "感谢亲爱的{用户名称}老板的点赞，谢谢你了", "谢谢{用户名称}老板给主播点赞，爱你哟"},
			IsModelRewrite:    false,
			RewriteFrequency:  0,
		},
		{
			ModuleType:        model.ModuleTypeAudio,
			OrderNum:          4,
			ModuleName:        "礼物模块",
			IntervalTimeStart: 5,
			IntervalTimeEnd:   10,
			TriggerConditions: model.JSONSlice{model.TriggerSendGift},
			ReadStep:          "sequential",
			ScriptContent: model.JSONSlice{"感谢{用户名称}老板送的{礼物名称}", "感谢{用户名称}大老板给我送的{礼物名称}，祝您吉祥如意笑口常开：身体健康，天天开心。",
				"感谢{用户名称}老板给我送的{礼物名称}，祝你好事都成双", "感谢{用户名称}给我送的{礼物名称}，祝你福气东来，鸿运通天",
				"感谢{用户名称}给我送的{礼物名称}，祝你愿你合家欢乐", "感谢{用户名称}老板给我送的{礼物名称}，祝你生意兴隆，财源广进"},
			IsModelRewrite:   false,
			RewriteFrequency: 0,
		},
		{
			ModuleType:        model.ModuleTypeAudio,
			OrderNum:          5,
			ModuleName:        "礼物模块",
			IntervalTimeStart: 5,
			IntervalTimeEnd:   10,
			TriggerConditions: model.JSONSlice{model.TriggerSendGift},
			ReadStep:          "sequential",
			ScriptContent: model.JSONSlice{"感谢{用户名称}老板送的{礼物名称}", "感谢{用户名称}大老板给我送的{礼物名称}，祝您吉祥如意笑口常开：身体健康，天天开心。",
				"感谢{用户名称}老板给我送的{礼物名称}，祝你好事都成双", "感谢{用户名称}给我送的{礼物名称}，祝你福气东来，鸿运通天",
				"感谢{用户名称}给我送的{礼物名称}，祝你愿你合家欢乐", "感谢{用户名称}老板给我送的{礼物名称}，祝你生意兴隆，财源广进"},
			IsModelRewrite:   false,
			RewriteFrequency: 0,
		},
		{
			ModuleType:        model.ModuleTypeBase,
			OrderNum:          6,
			ModuleName:        "基础模块",
			IntervalTimeStart: 5,
			IntervalTimeEnd:   10,
			TriggerConditions: model.JSONSlice{model.TriggerExecuteLoop},
			ReadStep:          "sequential",
			ScriptContent: model.JSONSlice{"我看下现在几点了，已经{现在时间}了",
				"今天是{现在日期},已经{现在时间}了",
				"直播间有多少人，我看下，现在直播间{在线人数}人了",
				"谢谢直播间家人们的支持，我先去休息换新的助理和大家讲解我下去给大家打包发货了，换新的助手来给大家互动",
				"不知不觉时间已经{现在时间}，下面换新的小助手，我休息去了直播间{在线人数}位宝宝，我们要换人咯，马上换新的小助手来",
			},
			IsModelRewrite:   false,
			RewriteFrequency: 0,
		},
		{
			ModuleType:        model.ModuleTypeAudio,
			OrderNum:          7,
			ModuleName:        "音频模块",
			IntervalTimeStart: 5,
			IntervalTimeEnd:   10,
			TriggerConditions: model.JSONSlice{model.TriggerExecuteLoop},
			AudioName:         "background.mp3",
			AudioPath:         "/path/to/background.mp3",
		},
	}

	for _, module := range baseModules {
		if err := db.Create(&module).Error; err != nil {
			return err
		}
	}
	return nil
}
