package pkg

import (
	"fmt"
	"testing"
)

func TestGetModelList(t *testing.T) {
	data, err := GetModelList()
	fmt.Println("GetModelList:ERROR: ", err)
	//遍历打印 data 中的每个元素的 Title 字段
	for _, v := range data {
		fmt.Printf("%+v\n", v.ModelName)
	}
	fmt.Println("-------------------")
	// GetRagList() 测试
	ragData, err := GetRagList()
	fmt.Println("GetRagList:ERROR: ", err)
	for _, v := range ragData {
		fmt.Printf("%+v\n", v)
	}
	fmt.Println("-------------------")
	searchRagData, err := SearchRagDocument([]string{"知识库1"}, "111")
	fmt.Println("SearchRagDocument:ERROR: ", err)
	for _, v := range searchRagData {
		fmt.Printf("%+v\n", v)
	}
}

func TestGetModelList2(t *testing.T) {
	// supplierName, model, parameters, title, agentName
	// fmt.Println("-------------------") b49c32b0-27fd-11f0-9f33-8bd9fef57116
	id, err := CreateChatApi("qwen2.5", "7b", "232435674323123")
	fmt.Println(id, err)
	fmt.Println("-------------------")
	str, err := SendChatApi(
		"qwen2.5", //模型名称 qwen3
		"7b",      //模型名称 0.6b
		id,        // 对话框id
		`
# 口播文案		
一号链接刷新去拍，马上封单下架了，已经有一个姐姐抢了五件了，喜欢的姐姐抓紧了一号链接，马上没有了
`, // 用户输入的消息
		nil,  //知识库检索结果
		true, //是否为临时对话
		`
你是一个专业的全球化直播文案优化助手。请按照以下要求处理和改善我提供的口播文案：

1. 语言转换：将原始文案精准改写为目标语种的本地化表达，保持口语化直播风格
2. 格式要求：返回标准的JSON格式，元素包含：
   - original_text: 原始文案
   - translated_text: 改善后的对应语种文案 
3. 文案规范：
   - 保留核心营销点，但根据目标语言文化调整表达方式
   - 添加适合直播场景的语气词和互动话术
4. 特殊处理：
   - 保持专业术语/品牌词不变
   - 数字保留原始格式
   - 本地化度量衡单位
   - 禁止使用不适合直播场景的语言
   - 避免使用不专业的语言
4. 目标语种：
   - 韩语

请返回可直接被JSON.parse解析的规范格式，不要任何额外说明。   

`, // 系统提示词
		"", // 知识库名称S
	)
	fmt.Println(err)
	fmt.Println(str)
	RemoveChatApi(id)
	// 由于 RemoveChatApi 未定义，暂时注释掉该调用
	// RemoveChatApi(id)
}
