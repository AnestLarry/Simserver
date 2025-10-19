package uploadGroup

import (
	argsConfig "Simserver/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

var (
	Args = argsConfig.GetConfig().Upload
)

func upload_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pathDict := map[string]bool{"/api/upload/": Args.Enable, "/api/upload/text": Args.UploadText}
		v, ok := pathDict[c.FullPath()]
		if !ok || !v {
			c.JSON(501, gin.H{"message": fmt.Sprintf("The server is not supported \"%s\"", c.FullPath())})
			c.Abort()
		}
	}
}

func Upload_routerGroup_init(Uploader_routerGroup *gin.Engine) {
	routerApi := Uploader_routerGroup.Group("/api/upload")
	routerApi.Use(upload_middleware())

	routerApi.POST("/", func(c *gin.Context) {
		HandleFileUpload(c, Args.SecureExt)
	})
	routerApi.POST("/text", HandleTextUpload)
}