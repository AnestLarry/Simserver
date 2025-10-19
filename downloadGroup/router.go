package downloadGroup

import (
	argsConfig "Simserver/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	Args            = argsConfig.GetConfig().Download
	DownloadCodeMap = map[string]DownloadCodeItem{}
)

func Downloader_routerGroup_init(Downloader_routerGroup *gin.Engine) {
	routerPage, routerApi := Downloader_routerGroup.Group("/dl"), Downloader_routerGroup.Group("/api/dl")
	if Args.DownloadCode {
		LoadDownloadCodeJson()
	}

	routerPage.Use(download_middleware())
	routerApi.Use(download_middleware())

	routerApi.GET("/n/*path", n_handler)
	routerApi.GET("/ls/*path", ls_handler)
	routerApi.GET("/zip/*path", zip_handler)
	routerApi.GET("/downloadCode/*dCode", downloadCode_handler)
}

func download_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()[:strings.LastIndex(c.FullPath(), "/")]
		pathDict := map[string]bool{"/api/dl/ls": Args.Ls, "/api/dl/zip": Args.Zip, "/api/dl/downloadCode": Args.DownloadCode, "/api/dl/n": true}
		v, ok := pathDict[path]
		if !ok || !v {
			c.JSON(501, gin.H{"message": fmt.Sprintf("The server is not supported \"%s\"", path)})
			c.Abort()
		}
	}
}