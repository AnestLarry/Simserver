package uploadGroup

import (
	"Simserver/Libs"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strings"
	"time"
)

var (
	Upload_open      = false
	Upload_text_open = false
)

func upload_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pathDict := map[string]bool{"/upload/": Upload_open, "/api/upload/": Upload_open, "/upload/text": Upload_text_open, "/api/upload/text": Upload_text_open}
		v, ok := pathDict[c.FullPath()]
		if !ok || !v {
			c.JSON(501, gin.H{"message": fmt.Sprintf("The server is not supported \"%s\"", c.FullPath())})
			c.Abort()
		}
	}
}

func Upload_routerGroup_init(Uploader_routerGroup *gin.Engine, staticFiles embed.FS) {
	routerPage, routerApi := Uploader_routerGroup.Group("/upload"), Uploader_routerGroup.Group("/api/upload")
	routerPage.Use(upload_middleware())
	routerApi.Use(upload_middleware())
	routerPage.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Stream(func(w io.Writer) bool {
			file, _ := staticFiles.ReadFile("static/upload.html")
			w.Write(file)
			return false
		})
	})
	routerApi.POST("/", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(500, gin.H{"message": "got file error."})
		}
		files := form.File["files"]
		for _, file := range files {
			if !Libs.LibsXExists("upload") {
				os.Mkdir("upload", 0764)
			}
			folder := fmt.Sprintf("upload/from_%s_", strings.ReplaceAll(c.ClientIP(), ".", "_"))
			if !Libs.LibsXExists(folder) {
				os.Mkdir(folder, 0764)
			}
			c.SaveUploadedFile(file, fmt.Sprintf("%s/%s_dat", folder, file.Filename))
		}
		c.JSON(200, gin.H{"message": "OK"})

	})

	routerPage.GET("/text", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Stream(func(w io.Writer) bool {
			file, _ := staticFiles.ReadFile("static/uploadText.html")
			w.Write(file)
			return false
		})
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
