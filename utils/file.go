package utils

import (
	"go-web/appError"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// 上传音频或图片
func UploadMedia(c *gin.Context) (string, error) {
	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		return "", appError.NewAppError(500, "文件上传失败", nil)
	}

	// 判断文件类型（根据 Content-Type）
	fileType := c.PostForm("type") // audio 或 image
	var saveDir string
	if fileType == "audio" {
		saveDir = "./uploads/audio/"
	} else if fileType == "image" {
		saveDir = "./uploads/image/"
	} else {
		return "", appError.NewAppError(500, "未知的文件类型", nil)
	}

	// 创建目录
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		log.Println("创建目录失败:", err)
		return "", appError.NewAppError(500, "创建目录失败", nil)
	}

	// 生成保存路径
	filename := time.Now().Format("20060102150405") + "_" + file.Filename
	savePath := filepath.Join(saveDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		log.Println("保存文件失败:", err)
		return "", appError.NewAppError(500, "保存文件失败", nil)
	}

	// 返回文件URL
	fileURL := "/uploads/" + fileType + "/" + filename
	return fileURL, nil
}
