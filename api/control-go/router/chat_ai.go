package router

import (
	"control-go/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ChatAiRoutes(r *gin.Engine) {
	//一次性调用API聊天 (不保存对话记录)
	r.POST("/disposable_send", func(c *gin.Context) {
		var data = struct {
			Model        string `json:"model"`         //qwen2.5
			Parameters   string `json:"parameters"`    //7b
			UserContent  string `json:"user_content"`  // 用户输入的消息
			SystemPrompt string `json:"system_prompt"` // 系统提示词
		}{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := pkg.CreateChatApi(data.Model, data.Parameters, uuid.New().String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		str, err := pkg.SendChatApi(
			data.Model,        //模型名称 qwen3
			data.Parameters,   //模型名称 0.6b
			id,                // 对话框id
			data.UserContent,  // 用户输入的消息
			nil,               //知识库检索结果
			true,              //是否为临时对话
			data.SystemPrompt, // 系统提示词
			"",                // 知识库名称S
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		pkg.RemoveChatApi(id)
		c.JSON(http.StatusOK, gin.H{"data": str})
	})
}
