package pkg

import (
	"bytes"
	"encoding/json"
	"net/http"
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
