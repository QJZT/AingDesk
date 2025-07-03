package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"time"

	"control-go/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getCPUInfo 获取CPU信息
func getCPUInfo() string {
	return fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
}

// getUsername 获取当前用户名
func getUsername() string {
	currentUser, err := user.Current()
	if err != nil {
		// 如果获取失败，尝试从环境变量获取
		if username := os.Getenv("USERNAME"); username != "" {
			return username
		}
		if username := os.Getenv("USER"); username != "" {
			return username
		}
		return "unknown"
	}
	return currentUser.Username
}

// ActivationRequest 激活码验证请求结构
type ActivationRequest struct {
	ActivationCode string `json:"activation_code" binding:"required"`
	CPUInfo        string `json:"cpu_info" binding:"required"`
	Username       string `json:"username" binding:"required"`
}

// validateActivationCode 验证激活码
func validateActivationCode(activationCode string) (bool, string, string) {
	// 服务器验证接口URL（请根据实际情况修改）
	serverURL := "http://localhost:7777/api/user/activation"

	// 获取系统信息
	cpuInfo := getCPUInfo()
	username := getUsername()

	// 构建请求数据
	requestData := ActivationRequest{
		ActivationCode: activationCode,
		CPUInfo:        cpuInfo,
		Username:       username,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return false, "请求数据格式错误", ""
	}

	// 创建HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// 发送POST请求
	resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Sprintf("网络请求失败: %v", err), ""
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "读取服务器响应失败", ""
	}

	// 解析响应
	var response struct {
		Success   bool   `json:"success"`
		Message   string `json:"message"`
		ExpiresAt string `json:"expires_at"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return false, "解析服务器响应失败", ""
	}

	return response.Success, response.Message, response.ExpiresAt
}

func SetupKvRoutes(r *gin.Engine, db *gorm.DB) {
	//获取数据map值
	r.POST("/get_kv", func(c *gin.Context) {
		var data = struct {
			Key string `json:"key"`
		}{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		KvStr := model.KvStr{}
		db.Where("key =?", data.Key).Find(&KvStr) //删除指定数据

		//如果没数据就创建个空的
		if KvStr.Key == "" {
			db.Create(&model.KvStr{Key: data.Key, V: make(map[string]interface{})})
			c.JSON(http.StatusOK, gin.H{"kv": make(map[string]interface{})})
			return
		}
		c.JSON(http.StatusOK, gin.H{"kv": KvStr.V})
	})
	//获取数据map值
	r.POST("/set_kv", func(c *gin.Context) {
		var data = struct {
			Key string                 `json:"key"`
			Kv  map[string]interface{} `json:"kv"`
		}{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		KvStr := model.KvStr{}
		db.Where("key =?", data.Key).Find(&KvStr) //删除指定数据

		//如果没数据就创建个空的
		if KvStr.Key == "" {
			db.Create(&model.KvStr{Key: data.Key, V: data.Kv})
			c.JSON(http.StatusOK, gin.H{"kv": KvStr.V})
			return
		}
		for k, v := range data.Kv {
			KvStr.V[k] = v
		}
		db.Save(&KvStr)
		c.JSON(http.StatusOK, gin.H{"kv": KvStr.V})
	})

	//获取activation_code
	r.POST("/client_config", func(c *gin.Context) {
		// 解析请求参数
		var requestData struct {
			ActivationCode string `json:"activation_code"`
		}
		c.ShouldBindJSON(&requestData)

		KvStr := model.KvStr{}
		db.Where("key =?", "client_config").Find(&KvStr) //获取指定数据

		//如果没数据，创建空配置
		if KvStr.Key == "" || KvStr.V == nil {
			db.Create(&model.KvStr{Key: "client_config", V: make(map[string]interface{})})
			KvStr.Key = "client_config"
			KvStr.V = make(map[string]interface{})
		}

		// 安全地获取激活码
		activationCodeInterface, exists := KvStr.V["activation_code"]
		var activationCode string
		if exists {
			activationCode, _ = activationCodeInterface.(string)
		}

		// 如果是唯一值 就不请求服务器确认
		// 如果是特殊激活码，直接验证通过并设置永久有效期
		if requestData.ActivationCode == "OUYIUEWQIYEIUIUDHASIDUHWQIOUHFIABIABWIQODGO" || activationCode == "OUYIUEWQIYEIUIUDHASIDUHWQIOUHFIABIABWIQODGO" {
			// 如果数据库中没有激活码且请求中有激活码，则设置激活码
			if KvStr.V["activation_code"] == nil || KvStr.V["activation_code"] == "" {
				if requestData.ActivationCode != "" {
					KvStr.V["activation_code"] = requestData.ActivationCode
				} else {
					// 如果请求中没有激活码，但数据库中的激活码是特殊激活码，保持不变
					KvStr.V["activation_code"] = "OUYIUEWQIYEIUIUDHASIDUHWQIOUHFIABIABWIQODGO"
				}
			}

			// 设置永久有效期至2099年
			KvStr.V["expires_at"] = "2099-01-01 00:00:00"
			db.Save(&KvStr)

			c.JSON(http.StatusOK, gin.H{
				"client_config":    KvStr.V,
				"activation_valid": true,
				"message":          "vip账号",
			})
			return
		}

		// 如果接口有传参（新激活码），则验证并更新
		if requestData.ActivationCode != "" {
			// 调用服务器验证新激活码
			isValid, message, expiresAt := validateActivationCode(requestData.ActivationCode)

			if isValid {
				// 验证通过，更新数据库中的激活码和有效期
				KvStr.V["activation_code"] = requestData.ActivationCode
				KvStr.V["expires_at"] = expiresAt
				db.Save(&KvStr)

				c.JSON(http.StatusOK, gin.H{
					"client_config":    KvStr.V,
					"activation_valid": true,
					"message":          "激活码验证成功",
				})
				return
			} else {
				// 验证失败
				c.JSON(http.StatusOK, gin.H{
					"client_config":    KvStr.V,
					"activation_valid": false,
					"message":          message,
				})
				return
			}
		}

		// 如果数据库中没有激活码
		if activationCode == "" {
			c.JSON(http.StatusOK, gin.H{
				"client_config":    KvStr.V,
				"activation_valid": false,
				"message":          "请填写激活码",
			})
			return
		}

		// 验证数据库中已有的激活码
		isValid, message, expiresAt := validateActivationCode(activationCode)

		// 如果验证成功且有新的有效期，更新数据库
		if isValid && expiresAt != "" {
			KvStr.V["expires_at"] = expiresAt
			db.Save(&KvStr)
		}

		c.JSON(http.StatusOK, gin.H{
			"client_config":    KvStr.V,
			"activation_valid": isValid,
			"message":          message,
		})
	})
}
