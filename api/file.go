package api

import (
	"go-web/appError"
	"go-web/utils"

	"github.com/gin-gonic/gin"
)

// UploadMedia 文件上传接口
// @Summary 上传文件（图片/音频）
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param type formData string true "文件类型" Enums(audio, image)
// @Success 200 {object} map[string]string "返回文件URL"
// @Router /api/v1/upload/media [post]
func UploadMedia() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *string) error {
			*req = ""
			return nil
		},
		func(c *gin.Context, req *string) (any, *appError.AppError) {
			fileURL, err := utils.UploadMedia(c)
			if err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return map[string]string{"url": fileURL}, nil
		},
	)
}
