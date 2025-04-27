package router

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"control-go/global"
	"control-go/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Name struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var names = []Name{}
var lastID = 0

func SetupNameRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/names", func(c *gin.Context) {
		var names []Name
		db.Find(&names)
		c.JSON(http.StatusOK, names)
	})
	r.POST("/set_name", func(c *gin.Context) {

		//读取请求的name string {"name": "123"}
		var data = struct {
			Name string `json:"name"`
		}{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//写入数据库
		db.Create(&Name{Name: data.Name})
		//返回成功
		c.JSON(http.StatusOK, gin.H{"message": "name set successfully"})
	})
	r.POST("/del_name", func(c *gin.Context) {
		var data = struct {
			Id int `json:"id"`
		}{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Where("id = ?", data.Id).Delete(&Name{}) //删除指定数据
		c.JSON(http.StatusOK, gin.H{"message": "name delete successfully"})
	})

	//获取指定文件夹下的音频文件列表
	r.POST("/get_name_list", func(c *gin.Context) {
		var data = struct {
			Id int `json:"id"`
		}{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		list := []model.FileData{}
		db.Where("associated_id =? and data_type = ?", data.Id, "人声音色").Find(&list) //删除指定数据
		c.JSON(http.StatusOK, gin.H{"message": "name delete successfully", "list": list})
	})

	r.POST("/upload", func(c *gin.Context) {

		// 获取关联的ID
		id := c.PostForm("id")
		fmt.Println("c.PostForm(id)：", id)
		// 处理文件上传
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get file",
			})
			fmt.Println("Failed to get file ", err)
			return
		}

		// 检查文件类型
		allowedTypes := []string{
			"audio/mpeg",      // MP3
			"audio/wav",       // WAV
			"audio/ogg",       // OGG
			"audio/aac",       // AAC
			"audio/mp4",       // MP4
			"audio/webm",      // WebM
			"audio/flac",      // FLAC
			"audio/x-m4a",     // M4A
			"audio/x-ms-wma",  // WMA
			"audio/x-aiff",    // AIFF
			"audio/x-wav",     // WAV (alternative)
			"audio/x-mpegurl", // M3U
			"audio/x-scpls",   // PLS
		}
		valid := false
		for _, t := range allowedTypes {
			if file.Header.Get("Content-Type") == t {
				valid = true
				break
			}
		}
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Only audio files are allowed (MP3, WAV, OGG)",
			})
			fmt.Println("Only audio files are allowed (MP3, WAV, OGG)")
			return
		}

		// 创建保存路径
		uploadPath := global.ToneFilePath + "/" + id
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create upload directory",
			})
			return
		}

		// 保存文件
		filePath := filepath.Join(uploadPath, file.Filename)
		// 检查文件是否已存在
		if _, err := os.Stat(filePath); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "文件已存在: " + file.Filename,
			})
			return
		}
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save file",
			})
			return
		}
		//id 转成 int类型
		idInt, _ := strconv.Atoi(id)

		//保存表数据 FileData
		fileData := model.FileData{
			FileType:     filepath.Ext(file.Filename), // 文件后缀名
			FileName:     file.Filename,               // 文件名称
			DataType:     "人声音色",                      // 假设数据类型为音频
			AssociatedID: idInt,                       // 关联ID
			FilePath:     filePath,                    // 文件绝对路径
			Description:  "人声音色",                      // 文件说明
			FileSize:     file.Size,                   // 文件大小（字节）
			Metadata:     "{}",                        // 文件其他信息（JSON格式）
		}
		if err := db.Create(&fileData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save file data",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded successfully",
			"path":    filePath,
		})
	})

	// 新增文件下载接口
	r.GET("/download/:id/:filename", func(c *gin.Context) {
		id := c.Param("id")
		filename := c.Param("filename")
		// 安全校验文件名
		if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件名"})
			return
		}
		// 构建文件路径
		filePath := filepath.Join(global.ToneFilePath, id, filename)
		filePath = filepath.Clean(filePath) // 防止路径遍历
		fmt.Println("filePath:", filePath)
		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
			return
		}
		// 设置响应头
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Expires", "0")
		c.Header("Cache-Control", "must-revalidate")
		c.Header("Pragma", "public")
		// 发送文件
		c.File(filePath)
	})

	// 删除文件接口
	r.GET("/delete_file/:id/:filename", func(c *gin.Context) {
		id := c.Param("id")
		filename := c.Param("filename")

		// 安全校验文件名
		if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件名"})
			return
		}

		// 构建文件路径
		filePath := filepath.Join(global.ToneFilePath, id, filename)
		filePath = filepath.Clean(filePath)

		// 删除物理文件
		os.Remove(filePath)
		// if err := os.Remove(filePath); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "文件删除失败"})
		// 	return
		// }

		// 删除数据库记录
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID格式"})
			return
		}

		// 删除对应表数据
		db.Where("associated_id = ? AND file_name = ?", idInt, filename).Delete(&model.FileData{})

		c.JSON(http.StatusOK, gin.H{
			"message": "文件删除成功",
		})
	})
}
