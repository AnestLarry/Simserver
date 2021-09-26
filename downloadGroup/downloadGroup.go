package downloadGroup

import (
	"Simserver/Libs"
	"archive/zip"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Ls_open           = false
	Dls_open          = false
	Zip_open          = false
	DownloadCode_open = false
	DownloadCodeMap   = map[string]DownloadCodeItem{}
)

type ItemField struct {
	Name string
	Path string
}

type DownloadCodeItem struct {
	Code  string
	Name  string
	Files []string
}

func Downloader_routerGroup_init(Downloader_routerGroup *gin.RouterGroup, staticFiles embed.FS, r *gin.Engine) {
	t, _ := template.ParseFS(staticFiles, "static/lists.html")
	r.SetHTMLTemplate(t)
	//Downloader_routerGroup.Use(download_middleware())
	Downloader_routerGroup.Use(download_middleware())
	Downloader_routerGroup.GET("/n/*path", func(c *gin.Context) {
		fileName := c.Param("path")[1:]
		if Libs.LibsXIsFile(fileName) {
			c.File(fileName)
		} else {
			c.JSON(404, gin.H{"message": "Not file found."})
		}
	})
	Downloader_routerGroup.GET("/ls/*path", func(c *gin.Context) {
		ls := getFilesLists(c.Param("path")[1:], c.Request.URL.String())
		c.HTML(200, "lists.html", gin.H{"type": "ls", "folderList": ls[0], "fileList": ls[1]})
	})
	Downloader_routerGroup.GET("/dls/*path", func(c *gin.Context) {
		dls := getFilesLists(c.Param("path")[1:], c.Request.URL.String())
		c.HTML(200, "lists.html", gin.H{"type": "dls", "folderList": dls[0], "fileList": dls[1]})
	})
	Downloader_routerGroup.GET("/zip/*path", func(c *gin.Context) {
		c.Writer.Header().Set("Content-type", "application/octet-stream")
		path := c.Param("path")[1:]
		path = strings.ReplaceAll(path, "/", "\\")
		c.Stream(func(w io.Writer) bool {
			ar := zip.NewWriter(w)
			c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip",
				time.Now().Format("2006-01-02--15-04-05")))
			filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
				if f.IsDir() {
					return nil
				} else {
					newPath := strings.ReplaceAll(p, path, "")
					if newPath[0] == '\\' {
						newPath = newPath[1:]
					}
					file, _ := os.Open(p)
					f, _ := ar.Create(newPath)
					io.Copy(f, file)
				}
				return nil
			})
			ar.Close()
			return false
		})
	})
	Downloader_routerGroup.GET("/downloadCode/*dCode", func(c *gin.Context) {
		c.Writer.Header().Set("Content-type", "application/octet-stream")
		dCode := c.Param("dCode")[1:]
		dCodeItem, ok := DownloadCodeMap[dCode]
		if !ok {
			c.JSON(403, gin.H{"message": "this Code is not support!"})
		}
		downloadCodeFiles := dCodeItem.Files
		c.Stream(func(w io.Writer) bool {
			ar := zip.NewWriter(w)
			if dCodeItem.Name != "" {
				c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip",
					dCodeItem.Name))
			} else {
				c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip",
					time.Now().Format("2006-01-02--15-04-05")))
			}
			for i := 0; i < len(downloadCodeFiles); i++ {
				v := downloadCodeFiles[i]
				file, _ := os.Open(v)
				newPath := strings.ReplaceAll(v, "\\", "/")
				f, _ := ar.Create(newPath[strings.LastIndex(newPath, "/")+1:])
				io.Copy(f, file)
			}
			ar.Close()
			return false
		})
	})
}

func download_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()[4:]
		path = path[:strings.Index(path, "/")]
		if path == "ls" {
			if !Ls_open {
				c.JSON(501, gin.H{"message": "The server is not supported \"list files\""})
				c.Abort()
			}
		} else if path == "dls" {
			if !Dls_open {
				c.JSON(501, gin.H{"message": "The server is not supported \"list downloadGroup files\""})
				c.Abort()
			}
		} else if path == "zip" {
			if !Zip_open {
				c.JSON(501, gin.H{"message": "The server is not supported \"zip\""})
				c.Abort()
			}
		} else if path == "downloadCode" {
			if !DownloadCode_open {
				c.JSON(501, gin.H{"message": "The server is not supported \"downloadCode\""})
				c.Abort()
			}
		} else if path == "n" {
		} else {
			c.JSON(501, gin.H{"message": "undefined."})
			c.Abort()
		}
	}
}

func getFilesLists(path, Request_URL_Path string) [][]ItemField {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	res := make([][]ItemField, 2)
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			res[0] = append(res[0], ItemField{file.Name(), Request_URL_Path + "/" + file.Name()})
		} else {
			res[1] = append(res[1], ItemField{file.Name(), "/dl/n/" + path + file.Name()})
		}
	}
	return res
}
func LoadDownloadCodeJson() {
	if !Libs.LibsXExists("./downloadCodes.json") {
		fmt.Println("  downloadCodes.json is not exist.")
		os.Exit(-1)
	} else {
		var downloadCodeJson []DownloadCodeItem
		downloadCodeFile, err := ioutil.ReadFile("./downloadCodes.json")
		if err != nil {
			fmt.Println("  open downloadCodes.json fail.")
			os.Exit(-1)
		}
		err = json.Unmarshal(downloadCodeFile, &downloadCodeJson)
		if err != nil {
			fmt.Printf("  downloadCodeFile fail to parse json.\n%v\n", err)
			os.Exit(-1)
		}
		for _, v := range downloadCodeJson {
			DownloadCodeMap[v.Code] = v
		}
	}
}