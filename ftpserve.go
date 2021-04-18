package main

import (
	"Libs/Libs"
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ls_open           = false
	dls_open          = false
	upload_open       = false
	zip_open          = false
	downloadCode_open = false
	Version           = "Apr 17,2021 Sa."
)

var (
	downloadCodeMap = map[string]downloadCodeItem{}
)

type downloadCodeItem struct {
	Code  string
	Name  string
	Files []string
}

func main() {
	ip := "0.0.0.0"
	port := "5000"
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "h", "-h", "help":
				fmt.Println(
					"Tips:\n" +
						" h  - show this help\n" +
						" v  - get version\n" +
						"Mode:\n" +
						" ls  - open ls function\n" +
						" dls  - add download links with the ls function's list\n" +
						" upload  - allow user upload files to host\n" +
						" zip  - allow zip dir for download (DANGER!)\n" +
						" downloadCode  - use download code to download a group file with setting\n" +
						"Args:\n" +
						" p / port  - use the port\n" +
						" ip  - use the ip.\n" +
						"Task:\n" +
						" RSUN  - reset files which in upload folder to origin's name")
				os.Exit(0)
			case "v", "-v", "version":
				fmt.Println(Version)
				os.Exit(0)
			case "ls", "-ls":
				fmt.Println(" -  ls mode on.")
				ls_open = true
			case "dls", "-dls":
				fmt.Println(" -  dls mode on.")
				dls_open = true
			case "ip", "-ip":
				i++
				ip = os.Args[i]
			case "p", "-p", "port", "-port":
				i++
				port = os.Args[i]
			case "upload", "-upload":
				fmt.Println(" -  upload mode on.")
				upload_open = true
				if Libs.LibsXExists("./upload") {
					if !Libs.LibsXIsDir("./upload") {
						fmt.Println("upload is not a folder!")
						os.Exit(0)
					}
				} else {
					os.Mkdir("./upload", 0764)
				}
			case "zip", "-zip":
				fmt.Println(" -  zip mode on.")
				zip_open = true
			case "downloadCode", "dC", "-downloadCode", "-dC":
				fmt.Println(" -  downloadCode mode on.")
				loadDownloadCodeJson()
				downloadCode_open = true
			case "RSUN", "-RSUN":
				if Libs.LibsXExists("./upload") {
					if Libs.LibsXIsDir("./upload") {
						fmt.Println("WARMING :  It may be fail or rewrite the same name file.\nkeyin \"y\" to continue or other to exit")
						{
							temp := ""
							fmt.Scanf("%s", &temp)
							if temp != "y" {
								os.Exit(0)
							}
						}
						files, _ := ioutil.ReadDir("./upload")
						for _, file := range files {
							temp := len(file.Name()) - 4
							if !file.IsDir() && file.Name()[temp:] == ".dat" {
								os.Rename(fmt.Sprintf("./upload/%s", file.Name()), fmt.Sprintf("./upload/%s", file.Name()[:temp]))
								fmt.Println("[", file.Name(), "] reset to [", file.Name()[:temp], "]")
							}
						}
						os.Exit(0)
					}
				}
			default:
				fmt.Println("There is arg(s) cannot parse. Do you need \"-h\" ?")
				os.Exit(0)
			}
		}
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		r := gin.H{"method": param.Method, "StatusCode": param.StatusCode, "ClientIP": param.ClientIP,
			"TimeStamp": param.TimeStamp.Format(time.RFC1123), "Path": param.Path, "Request.Proto": param.Request.Proto,
			"Latency": param.Latency, "User-Agent": param.Request.UserAgent(), "ErrorMessage": param.ErrorMessage}
		return fmt.Sprintf("%v\n", r)
	}))
	r.Use(gin.Recovery())
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": Version})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "It's the file download server." +
			"You can use the path to download the file on the machine."})
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "404 Not Found"})
	})
	Uploader_routerGroup := r.Group("/upload")
	Uploader_routerGroup.Use(upload_middleware())
	Uploader_routerGroup.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<form action="/upload/" method="post" enctype="multipart/form-data">
    Files: <input type="file" name="files" multiple><br><br>
    <input type="submit" value="Submit">
</form>`)
	})
	Uploader_routerGroup.POST("/", func(c *gin.Context) {
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

	Downloader_routerGroup := r.Group("/dl")
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
		path := c.Param("path")[1:]
		c.Header("Content-Type", "text/html; charset=utf-8")
		files := getFilesLists(path, c, "ls")
		c.String(200, files)
	})
	Downloader_routerGroup.GET("/dls/*path", func(c *gin.Context) {
		path := c.Param("path")[1:]
		c.Header("Content-Type", "text/html; charset=utf-8")
		files := getFilesLists(path, c, "dls")
		c.String(200, files)
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
		dCodeItem, ok := downloadCodeMap[dCode]
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
	fmt.Println(strings.Repeat("-", 15) + "\n" + fmt.Sprintf("%s:%s", ip, port))
	r.Run(fmt.Sprintf("%s:%s", ip, port))
}

func upload_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !upload_open {
			c.JSON(501, gin.H{"message": "The server is not supported \"upload\""})
			c.Abort()
		}
	}
}

func download_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()[4:]
		path = path[:strings.Index(path, "/")]
		if path == "ls" {
			if !ls_open {
				c.JSON(501, gin.H{"message": "The server is not supported \"list files\""})
				c.Abort()
			}
		} else if path == "dls" {
			if !dls_open {
				c.JSON(501, gin.H{"message": "The server is not supported \"list download files\""})
				c.Abort()
			}
		} else if path == "zip" {
			if !zip_open {
				c.JSON(501, gin.H{"message": "The server is not supported \"zip\""})
				c.Abort()
			}
		} else if path == "downloadCode" {
			if !downloadCode_open {
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

func loadDownloadCodeJson() {
	if !Libs.LibsXExists("./downloadCodes.json") {
		fmt.Println("  downloadCodes.json is not exist.")
		os.Exit(-1)
	} else {
		var downloadCodeJson []downloadCodeItem
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
			downloadCodeMap[v.Code] = v
		}
	}
}

func getFilesLists(path string, c *gin.Context, listType string) string {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	skillFolder := path
	result := "<html>Items:<br><br>"
	fs, ds := "<p> File:<br>", "<p> Dir:<br>"
	files, _ := ioutil.ReadDir(skillFolder)
	if listType == "ls" {
		for _, file := range files {
			if file.IsDir() {
				ds += "  " + file.Name() + "<br>"
			} else {
				fs += "  " + file.Name() + "<br>"
			}
		}
	} else if listType == "dls" {
		for _, file := range files {
			if file.IsDir() {
				ds += "  <a href='" + c.Request.URL.Path + "/" + file.Name() + "'>" + file.Name() + "</a><br>"
			} else {
				fs += "  <a href='/dl/n/" + path + file.Name() + "'>" + file.Name() + "</a><br>"
			}
		}
	}
	ds += "</p>"
	fs += "</p>"
	result = result + ds + "<br>" + fs + "</html>"
	return result
}
