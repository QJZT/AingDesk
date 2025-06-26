package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

// API基础URL
const aing_url = "http://127.0.0.1:7071"

// 模型响应结构体
type ModelResponse struct {
	Status   int                    `json:"status"`    // 状态码
	Code     int                    `json:"code"`      // 业务码
	Msg      string                 `json:"msg"`       // 消息
	ErrorMsg string                 `json:"error_msg"` // 错误信息
	Message  map[string]interface{} `json:"message"`   // 消息体
}

// 模型信息结构体
type ModelInfo struct {
	Title         string   `json:"title"`         // 模型标题
	SupplierName  string   `json:"supplierName"`  // 供应商名称
	Model         string   `json:"model"`         // 模型名称
	Size          int64    `json:"size"`          // 模型大小(字节)
	ContextLength int      `json:"contextLength"` // 上下文长度
	Capability    []string `json:"capability"`    // 能力列表
}
type ModelList struct {
	ModelInfo []ModelInfo `json:"modelInfo"` // 模型信息列表
	ModelName string      `json:"modelName"` // 模型名称
}

// 获取模型列表
func GetModelList() ([]ModelList, error) {
	url := aing_url + "/chat/get_model_list"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ModelResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var retList = make([]ModelList, 0)
	for k, v := range result.Message {
		var models = make([]ModelInfo, 0)
		if ollama, ok := v.([]interface{}); ok {
			for _, m := range ollama {
				if model, ok := m.(map[string]interface{}); ok {
					var info ModelInfo
					if data, err := json.Marshal(model); err == nil {
						json.Unmarshal(data, &info)
						models = append(models, info)
					}
				}
			}
		}
		retList = append(retList, ModelList{ModelInfo: models, ModelName: k})
	}
	return retList, nil
}

// 知识库信息结构体
type RagInfo struct {
	RagName             string  `json:"ragName"`             // 知识库名称
	RagDesc             string  `json:"ragDesc"`             // 知识库描述
	RagCreateTime       int64   `json:"ragCreateTime"`       // 创建时间(时间戳)
	SupplierName        string  `json:"supplierName"`        // 供应商名称
	EmbeddingModel      string  `json:"embeddingModel"`      // 嵌入模型
	SearchStrategy      int     `json:"searchStrategy"`      // 搜索策略
	MaxRecall           int     `json:"maxRecall"`           // 最大召回数量
	RecallAccuracy      float64 `json:"recallAccuracy"`      // 召回准确率
	ResultReordering    int     `json:"resultReordering"`    // 结果重排序
	RerankModel         string  `json:"rerankModel"`         // 重排序模型
	QueryRewrite        int     `json:"queryRewrite"`        // 查询重写
	VectorWeight        float64 `json:"vectorWeight"`        // 向量权重
	KeywordWeight       float64 `json:"keywordWeight"`       // 关键词权重
	EmbeddingModelExist bool    `json:"embeddingModelExist"` // 嵌入模型是否存在
	ErrorMsg            string  `json:"errorMsg"`            // 错误信息
}

// 知识库列表响应结构体
type RagListResponse struct {
	Status   int       `json:"status"`    // 状态码
	Code     int       `json:"code"`      // 业务码
	Msg      string    `json:"msg"`       // 消息
	ErrorMsg string    `json:"error_msg"` // 错误信息
	Message  []RagInfo `json:"message"`   // 知识库列表
}

// 获取知识库列表
func GetRagList() ([]RagInfo, error) {
	url := aing_url + "/rag/get_rag_list"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RagListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Message, nil
}

// 知识库检索请求结构体
type RagSearchRequest struct {
	RagList   string `json:"ragList"`   // 知识库列表(JSON字符串)
	QueryText string `json:"queryText"` // 查询内容
}

// 知识库检索结果项结构体
type RagSearchResultItem struct {
	ID           string  `json:"id"`           // 文档唯一ID
	Doc          string  `json:"doc"`          // 文档内容
	DocId        string  `json:"docId"`        // 文档ID
	DocName      string  `json:"docName"`      // 文档名称
	DocFile      string  `json:"docFile"`      // 文档文件路径
	Score        float64 `json:"score"`        // 综合评分
	VectorScore  float64 `json:"vectorScore"`  // 向量评分
	KeywordScore float64 `json:"keywordScore"` // 关键词评分
}

// 知识库检索响应结构体
type RagSearchResponse struct {
	Status   int                   `json:"status"`    // 状态码
	Code     int                   `json:"code"`      // 业务码
	Msg      string                `json:"msg"`       // 消息
	ErrorMsg string                `json:"error_msg"` // 错误信息
	Message  []RagSearchResultItem `json:"message"`   // 检索结果列表
}

// 检索知识库
func SearchRagDocument(ragList []string, queryText string) ([]RagSearchResultItem, error) {
	url := aing_url + "/rag/search_document"
	ragListStr := "["
	for i, v := range ragList {
		if i > 0 {
			ragListStr += ","
		}
		ragListStr += `"` + v + `"`
	}
	ragListStr += "]"
	requestBody := RagSearchRequest{
		RagList:   ragListStr,
		QueryText: queryText,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RagSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Message, nil
}

// 创建对话请求结构体
type CreateChatRequest struct {
	Model      string `json:"model"`      // 模型名称
	Parameters string `json:"parameters"` // 参数规模版本
	Title      string `json:"title"`      // 对话标题
}

// 创建对话响应结构体
type CreateChatResponse struct {
	Status   int    `json:"status"`    // 状态码
	Code     int    `json:"code"`      // 业务码
	Msg      string `json:"msg"`       // 消息
	ErrorMsg string `json:"error_msg"` // 错误信息
	Message  struct {
		ContextId string `json:"context_id"` // 对话ID
	} `json:"message"` // 对话ID
}

// // 创建新对话
func CreateChatApi(model, parameters, title string) (string, error) {
	url := aing_url + "/chat/create_chat"
	jsonBody, err := json.Marshal(CreateChatRequest{
		Model:      model,
		Parameters: parameters,
		Title:      title,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	retBytes, _ := ioutil.ReadAll(resp.Body)
	// {"status":0,"code":200,"msg":"对话创建成功","error_msg":"","message":{"model":"0.6b","title":"test0032132101","parameters":"qwen3","contextPath":"C:\\Users\\qjzt\\Desktop\\无人直播\\data\\context\\3c163cd0-27f0-11f0-a60e-2739721944a6","context_id":"3c163cd0-27f0-11f0-a60e-2739721944a6","agent_name":"","create_time":1746257320}}
	context_id := gjson.Get(string(retBytes), "message.context_id").String()
	if context_id == "" {
		return "", fmt.Errorf("context_id is empty")
	}
	return context_id, nil
}

// 发送向AI消息
/*
	{
	"model": "qwen3",
	"parameters": "0.6b",
	"context_id": "d8904650-27f5-11f0-90c2-69aa63d3138c",
	"search": "",
	"rag_list": "[]",
	"supplierName": "ollama",
	"temp_chat": "false",
	"mcp_servers": [],
	"user_content": "678",
	"images": "",
	"doc_files": ""
	}
*/
func SendChatApi(
	model, //模型名称 qwen3
	parameters, //模型名称 0.6b
	context_id, // 对话框id
	user_content string, // 用户输入的消息
	rag_results interface{}, //知识库检索结果
	temp_chat bool, //是否为临时对话
	system_prompt string, // 系统提示词
	regenerate_id string, // 重生成对话id
	supplierName string, //供应商名称,
) (string, error) {
	url := aing_url + "/chat/chat"
	mapData := map[string]interface{}{}
	//模型名称 qwen3
	mapData["model"] = model
	//模型名称 0.6b
	mapData["parameters"] = parameters
	mapData["supplierName"] = supplierName
	//对话框id
	mapData["context_id"] = context_id
	// 用户输入的消息
	mapData["user_content"] = user_content
	//知识库检索结果
	//[ { "id": "bc78dd90-fa75-11ef-a96a-4b17421e6c04", "doc": "其它SSL证书\n\n1. 什么是其它 SSL 证书？\n\n当您自行在其他厂商申请 SSL 证书（如腾讯云、阿里云）并部署到宝塔面板后，面板会将该证书归类到其他证书分类，您可在面板左侧导航栏SSL模块-其他证书中查看管理，支持删除、部署、下载、上传到云端等。\n\n图片\n\n》", "docId": "bc78dd90-fa75-11ef-a96a-4b17421e6c04", "docName": "其它SSL证书.docx", "docFile": "F:\desk-top-ai\data\rag\堡塔\markdown\其它SSL证书.docx.md", "score": 0.23035039901733398, "vectorScore": 0.3290719985961914, "keywordScore": 0 }, ...]
	mapData["rag_results"] = rag_results
	//是否为临时对话：true = 是, false = 否请传字符串如果temp_chat=true，
	// 不再传递历史消息给模型，每一个问题都独立回答作用：解决ollama上下文长度不足的情况下，
	// 大模型在第二次之后的回答容易产生幻觉的问题缺点：开启后，大模型不知道你之前问过什么，
	// 也不知道之前回复过什么
	mapData["temp_chat"] = temp_chat
	// 系统提示词
	mapData["system_prompt"] = system_prompt
	// 系统提示词
	mapData["mcp_servers"] = make([]string, 0)

	mapData["regenerate_id"] = regenerate_id

	jsonBody, err := json.Marshal(mapData)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	retBytes, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(retBytes))
	// fmt.Println("-----------------------")
	//-----------------------
	// req2, _ := http.NewRequest("POST", aing_url+"/chat/get_last_chat_history",
	// 	bytes.NewBuffer(jsonBody))
	// req2.Header.Set("Content-Type", "application/json")
	// client2 := &http.Client{}
	// resp2, err := client2.Do(req2)
	// if err != nil {
	// 	return "", err
	// }
	// defer resp2.Body.Close()
	// retBytes2, _ := ioutil.ReadAll(resp2.Body)
	// fmt.Println(string(retBytes2))
	return string(retBytes), nil
}

// /chat/remove_chat
func RemoveChatApi(context_id string) error {
	url := aing_url + "/chat/remove_chat"
	mapData := map[string]interface{}{}
	mapData["context_id"] = context_id
	jsonBody, err := json.Marshal(mapData)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", url,
		bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// POST /chat/get_last_chat_history
