package seed

import (
	"control-go/model"

	"gorm.io/gorm"
)

func SeedBaseModule(db *gorm.DB) error {
	var count int64
	db.Model(&model.BaseModule{}).Count(&count)
	if count == 0 {
		baseModules := []model.BaseModule{
			{
				ModuleType:        model.ModuleTypeBase,
				OrderNum:          0,
				ModuleName:        "控场循环", //22222
				IntervalTimeStart: 0,
				IntervalTimeEnd:   1,
				TriggerConditions: model.JSONSlice{model.TriggerSceneLoop},
				ReadStep:          "sequential",
				ScriptContent:     model.JSONSlice{"升级版电动打地鼠来啦！七彩炫光+魔性音效，沉浸式解压体验！", "内置3种游戏模式！挑战模式、音乐模式、宝宝模式，全家都能玩到停不下来！", "智能感应地鼠！反应速度可调节，从慢速到疯狂模式，锻炼孩子手眼协调！", "超长续航！Type-C快充2小时，畅玩一整天！露营聚餐都能带！", "职场人看这里！工作烦躁时敲敲地鼠，灯光随节奏闪烁，压力秒变多巴胺！", "安全细节拉满！ABS防火材质+硅胶软锤，欧盟CE认证，宝宝啃咬也不怕！", "宝妈实测！熊孩子玩到没空拆家！自动休眠功能，省电更省妈！", "宝妈实测：熊孩子专注玩半小时！旅行聚餐哄娃神器！", "幼儿园同款教具！拓展玩法：DIY纸杯地鼠，创造力UP！", "家人们！看过来！手速慢的真的要拍大腿！这款能让你玩到上头的解压打地鼠反应力按乐终于到货了！看到没？小地鼠 “哒哒哒” 随机冒头，手指要眼疾手快狂按下去，玩一把肾上腺素直接拉满！", "上班族下班解压神器说的就是它！摸鱼来两把瞬间赶走工作压力，学生党课间玩还能锻炼反应力和专注力！", "宝妈宝爸们看过来！是不是总愁孩子在家只知道玩手机、看电视？想陪孩子玩又找不到合适的项目？救命！我挖到宝了！这款全家都能玩到尖叫的打地鼠玩具，能让孩子主动放下平板，拉着你玩亲子对战！有位妈妈和我说，现在每天晚饭后的 “打地鼠时间”，成了家里最温馨欢乐的时刻！别犹豫，幸福亲子时光从现在开始！"},
			},
			{
				ModuleType:        model.ModuleTypeBase,
				OrderNum:          1,
				ModuleName:        "定时播报",
				IntervalTimeStart: 60,
				IntervalTimeEnd:   65,
				TriggerConditions: model.JSONSlice{model.TriggerIntervalLoop},
				ReadStep:          "random",
				ScriptContent:     model.JSONSlice{"现在时间{现在时间(时分)}，不要再犹豫了，快下单", "再过五分钟，{再过五分钟(时分)}了"},
			},
			{
				ModuleType:        model.ModuleTypeBase,
				OrderNum:          2,
				ModuleName:        "弹幕回答",
				IntervalTimeStart: 5,
				IntervalTimeEnd:   6,
				TriggerConditions: model.JSONSlice{model.TriggerBarrageComment},
				ReadStep:          "sequential",
				ScriptContent:     model.JSONSlice{"# 角色定位\n你是一个抖音的专业的卖货主播，为自己的店铺售卖产品。\n\n# 售卖产品信息\n产品名：蒜鸟玩具儿童解压神器记忆思维闯关益智减压游戏机打地鼠\n价格:：20美元\n规格：红色，蓝色，绿色\n\n# 行为回复\n控制再15个字左右去回答用户的问题\n拟人的语法，不要让人感觉出你是AI\n友好的与观众互动\n禁止回答抖音违规词语"},
			},
			{
				ModuleType:        model.ModuleTypeBase,
				OrderNum:          3,
				ModuleName:        "谢送礼",
				IntervalTimeStart: 15,
				IntervalTimeEnd:   20,
				TriggerConditions: model.JSONSlice{model.TriggerSendGift},
				ReadStep:          "sequential",
				ScriptContent:     model.JSONSlice{"感谢{送礼用户名},送我的{礼物名称}"},
				IsModelRewrite:    false,
				RewriteFrequency:  0,
			},
			{
				ModuleType:        model.ModuleTypeBase,
				OrderNum:          4,
				ModuleName:        "谢点赞",
				IntervalTimeStart: 20,
				IntervalTimeEnd:   25,
				TriggerConditions: model.JSONSlice{model.TriggerLike},
				ReadStep:          "sequential",
				ScriptContent:     model.JSONSlice{"感谢{点赞用户名}的点赞，谢谢支持", "{点赞用户名}谢谢哟，谢谢你的点赞"},
				IsModelRewrite:    false,
				RewriteFrequency:  0,
			},
			{
				ModuleType:        model.ModuleTypeBase,
				OrderNum:          5,
				ModuleName:        "进入直播间触发",
				IntervalTimeStart: 7,
				IntervalTimeEnd:   15,
				TriggerConditions: model.JSONSlice{model.TriggerEnterLiveRoom},
				ReadStep:          "sequential",
				ScriptContent:     model.JSONSlice{"{进入直播间用户名} 你来拉，看看有没有喜欢的", "{进入直播间用户名}看到你了，快选选看有没有喜欢的"},
			},
			{
				ModuleType:        model.ModuleTypeBase,
				OrderNum:          6,
				ModuleName:        "谢关注",
				IntervalTimeStart: 5,
				IntervalTimeEnd:   10,
				TriggerConditions: model.JSONSlice{model.TriggerFollowRoom},
				ReadStep:          "sequential",
				ScriptContent:     model.JSONSlice{"{关注直播用户名}感谢关注"},
				IsModelRewrite:    false,
				RewriteFrequency:  0,
			},
		}

		for _, module := range baseModules {
			if err := db.Create(&module).Error; err != nil {
				return err
			}
		}
	}
	var count2 int64
	db.Model(&model.KvStr{}).Count(&count2)
	if count2 == 0 {
		kvStrs := []model.KvStr{
			{Key: "prompt_", V: map[string]interface{}{"data": "你是一个专业的全球化直播文案优化助手。请按照以下要求处理和改善我提供的口播文案：\n\n1. 语言转换：将原始文案精准改写为目标语种的本地化表达，保持口语化直播风格\n2. 格式要求：返回标准的JSON格式，元素包含：\n   - original_text: 原始文案\n   - translated_text: 改善后的对应语种文案 \n3. 文案规范：\n   - 保留核心营销点，但根据目标语言文化调整表达方式\n   - 添加适合直播场景的语气词和互动话术\n4. 特殊处理：\n   - 保持专业术语/品牌词不变\n   - 数字保留原始格式\n   - 本地化度量衡单位\n   - 禁止使用不适合直播场景的语言\n   - 避免使用不专业的语言\n4. 目标语种：\n   - {语种}\n\n请返回可直接被JSON.parse解析的规范格式，不要任何额外说明。  "}},
			{Key: "prompt_rewrite", V: map[string]interface{}{"data": "你是一个专业的全球化直播文案优化助手。请按照以下要求处理和改善我提供的口播文案：\n\n1. 语言转换：将原始文案精准改写为目标语种的本地化表达，保持口语化直播风格\n2. 格式要求：直接给出改善后的对应语种文案 \n3. 文案规范：\n   - 保留核心营销点，但根据目标语言文化调整表达方式\n   - 添加适合直播场景的语气词和互动话术\n4. 特殊处理：\n   - 保持专业术语/品牌词不变\n   - 数字保留原始格式\n   - 本地化度量衡单位\n   - 禁止使用不适合直播场景的语言\n   - 避免使用不专业的语言\n4. 目标语种：\n   - {语种}\n\n请直接返回出改善后的对应语种文案，不要任何特殊符号，不要任何额外说明。"}},
		}
		for _, data := range kvStrs {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
