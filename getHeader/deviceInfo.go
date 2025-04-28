package getheader

import "github.com/gin-gonic/gin"

func GetDeviceInfoHeader(c *gin.Context) string {

	deviceInfo := c.GetHeader("DeviceInfo")
	if deviceInfo != "" {
		deviceInfo = c.GetHeader("deviceInfo")
	}
	return deviceInfo
}
