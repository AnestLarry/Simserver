package uploadGroup

import (
	"Simserver/Libs"
	argsConfig "Simserver/config"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Enable = false
	acs    = argsConfig.ArgConfigInit()
)

func upload_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pathDict := map[string]bool{"/api/upload/": Enable, "/api/upload/text": Enable}
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
		folder := fmt.Sprintf("upload/from_%s_", strings.ReplaceAll(c.ClientIP(), ".", "_"))
		if !Libs.LibsXExists(folder) {
			os.MkdirAll(folder, 0664)
		}

		x_file_name := c.Request.Header.Get("x-file-name")
		if len(x_file_name) == 0 {
			c.JSON(500, gin.H{"message": "Get file name failed"})
			return
		}
		x_file_name_byte, err := base64.StdEncoding.DecodeString(x_file_name)
		x_file_name = string(x_file_name_byte)
		if err != nil {
			c.JSON(500, gin.H{"message": "Decode file name failed"})
			return
		}
		destFileName := fmt.Sprintf("%s/%s", folder, x_file_name)
		if acs.Upload.SecureExt {
			destFileName = fmt.Sprintf("%s/%s_dat", folder, x_file_name)
		}
		errMsg, err := streamToFile(&c.Request.Body, destFileName)
		if err != nil {
			c.JSON(500, errMsg)
			return
		}
		c.JSON(200, gin.H{"message": "Upload success"})
	})

	routerApi.POST("/text", func(c *gin.Context) {
		text := c.PostForm("text")
		f, err := os.Create(fmt.Sprintf("./upload/%s.txt", time.Now().Format("2006-01-02--15-04-05")))
		if err != nil || text == "" {
			fmt.Println("fail to get text.")
			c.JSON(500, gin.H{"message": "fail to get text."})
			f.Close()
			return
		}
		_, err = f.WriteString(text)
		f.Close()
		if err != nil {
			fmt.Println("fail to save text file.")
			c.JSON(500, gin.H{"message": "fail to save text file."})
			return
		}
		c.JSON(200, gin.H{"message": "OK"})
	})
}

func streamToFile(fileHeader *io.ReadCloser, destPath string) (string, error) {

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Sprintf("Created target file error: %v", err), err
	}
	defer out.Close()

	_, err = io.Copy(out, *fileHeader)
	if err != nil {
		return fmt.Sprintf("Writed target file error: %v", err), err
	}

	return "", nil
}
