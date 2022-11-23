package downloadGroup

import (
	"Simserver/Libs"
	"archive/zip"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Ls_open           = false
	Zip_open          = false
	DownloadCode_open = false
	DownloadCodeMap   = map[string]DownloadCodeItem{}
)

type ItemField struct {
	Name string
	Path string
	Size float32
}

type DownloadCodeItem struct {
	Code  string   `json:"Code"`
	Name  string   `json:"Name"`
	Files []string `json:"Files"`
}

func Downloader_routerGroup_init(Downloader_routerGroup *gin.Engine, staticFiles embed.FS) {
	routerPage, routerApi := Downloader_routerGroup.Group("/dl"), Downloader_routerGroup.Group("/api/dl")
	if DownloadCode_open {
		LoadDownloadCodeJson()
	}

	routerPage.Use(download_middleware())
	routerApi.Use(download_middleware())

	routerApi.GET("/n/*path", func(c *gin.Context) {
		fileName := c.Param("path")[1:]
		if Libs.LibsXIsFile(fileName) {
			c.File(fileName)
		} else {
			c.JSON(404, gin.H{"message": "Not file found."})
		}
	})
	routerApi.GET("/ls/*path", func(c *gin.Context) {
		ls := getFilesLists(c.Param("path")[1:], c.Request.URL.String())
		c.JSON(200, gin.H{"folderList": ls[0], "fileList": ls[1]})
	})
	routerApi.GET("/zip/*path", func(c *gin.Context) {
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
					newPath = strings.ReplaceAll(newPath, "\\", "/")
					if newPath[0] == '/' {
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
	routerApi.GET("/downloadCode/*dCode", func(c *gin.Context) {
		dCode := c.Param("dCode")[1:]
		dCodeItem, ok := DownloadCodeMap[dCode]
		if !ok {
			c.JSON(403, gin.H{"message": "this Code is not support!"})
		} else {
			c.Writer.Header().Set("Content-type", "application/octet-stream")
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
		}
	})
}

func download_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()[:strings.LastIndex(c.FullPath(), "/")]
		pathDict := map[string]bool{"/api/dl/ls": Ls_open, "/api/dl/zip": Zip_open, "/api/dl/downloadCode": DownloadCode_open, "/api/dl/n": true}
		v, ok := pathDict[path]
		if !ok || !v {
			c.JSON(501, gin.H{"message": fmt.Sprintf("The server is not supported \"%s\"", path)})
			c.Abort()
		}
	}
}

func getFilesLists(path, Request_URL_Path string) [][]ItemField {
	res := make([][]ItemField, 2)
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			res[0] = append(res[0], ItemField{file.Name(), fmt.Sprintf("%s/%s", path, file.Name()), 0.0})
		} else {
			res[1] = append(res[1], ItemField{file.Name(), fmt.Sprintf("/api/dl/n/%s/%s", path, file.Name()), float32(file.Size()) / 1048576}) //MB
		}
	}
	for resI := range res {
		sort.Slice(res[resI], func(i, j int) bool {
			return strings.ToLower(res[resI][i].Name) < strings.ToLower(res[resI][j].Name)
		})
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
