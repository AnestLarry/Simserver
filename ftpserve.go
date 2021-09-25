package main

import (
	"Libs/Libs"
	"archive/zip"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
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
	upload_text_open  = false
	zip_open          = false
	downloadCode_open = false
	view_open         = false
	log_file_open     = false
	Version           = "Aut, 2021"
)

var (
	downloadCodeMap = map[string]downloadCodeItem{}
	//go:embed static
	staticFiles embed.FS
	//go:embed view
	viewFiles embed.FS
)

type downloadCodeItem struct {
	Code  string
	Name  string
	Files []string
}

type itemField struct {
	Name string
	Path string
}

func main() {
	ip := "0.0.0.0"
	port := "5000"
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "h", "-h", "help":
				fmt.Println(`Tips:
 h  - show this help
 v  - get version
Mode:
 ls  - open ls function
 dls  - add download links with the ls function's list
 upload  - allow user upload files to host
 uploadText  - allow user fill textarea to save text in txt
 zip  - allow zip dir for download (DANGER!)
 log  - put log in file
 downloadCode  - use download code to download a group file with setting
Args:
 p / port  - use the port
 ip  - use the ip.
Task:
 RSUN  - reset files which in upload folder to origin's name`)
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
			case "uploadText", "-uploadText", "uT", "-uT":
				fmt.Println(" -  upload text mode on.")
				upload_text_open = true
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
			case "view", "-view":
				fmt.Println(" -  view mode on.")
				view_open = true
			case "log", "-log":
				fmt.Println(" -  log file mode on.")
				log_file_open = true
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
								fmt.Printf("[%s] reset to [%s]\n", file.Name(), file.Name()[:temp])
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
	if log_file_open {
		gin.DisableConsoleColor()
		var f *os.File
		var err error
		if !Libs.LibsXExists("ftps.log") {
			f, err = os.Create("ftps.log")
			if err != nil {
				panic(err)
			}
		} else if Libs.LibsXIsFile("ftps.log") {
			f, err = os.OpenFile("ftps.log", 0666, os.ModeAppend)
			if err != nil {
				panic(err)
			}
		}
		gin.DefaultWriter = io.MultiWriter(f)
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
		c.JSON(200, gin.H{"message": `It's a file download server. You can transfer the file with the machine.`})
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "404 Not Found"})
	})
	routerGroup_init(r)
	fmt.Printf("%s\n%s:%s\n", strings.Repeat("-", 15), ip, port)
	r.Run(fmt.Sprintf("%s:%s", ip, port))
}

func routerGroup_init(r *gin.Engine) {
	// Uploader routerGroup
	func(Uploader_routerGroup *gin.RouterGroup) {
		Uploader_routerGroup.Use(upload_middleware())
		Uploader_routerGroup.GET("/", func(c *gin.Context) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Stream(func(w io.Writer) bool {
				file, _ := staticFiles.ReadFile("static/upload.html")
				w.Write(file)
				return false
			})
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
		Uploader_routerGroup.GET("/text", func(c *gin.Context) {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Stream(func(w io.Writer) bool {
				file, _ := staticFiles.ReadFile("static/uploadText.html")
				w.Write(file)
				return false
			})
		})
		Uploader_routerGroup.POST("/text", func(c *gin.Context) {
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
	}(r.Group("/upload"))
	// Downloader routerGroup
	func(Downloader_routerGroup *gin.RouterGroup) {
		t, _ := template.ParseFS(staticFiles, "static/lists.html")
		r.SetHTMLTemplate(t)
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
	}(r.Group("/dl"))
	// View routerGroup
	func(View_routerGroup *gin.RouterGroup) {
		View_routerGroup.Use(view_middleware())
		views, err := fs.Sub(viewFiles, "view/h5player")
		if err != nil {
			panic(err)
		}
		View_routerGroup.StaticFS("/h5player", http.FS(views))

	}(r.Group("/view"))
}

func upload_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()[7:]
		if path == "/" {
			if !upload_open {
				c.JSON(501, gin.H{"message": "The server is not supported \"upload\""})
				c.Abort()
			}
		} else if path == "/text" {
			if !upload_text_open {
				c.JSON(501, gin.H{"message": "The server is not supported \"uploadText\""})
				c.Abort()
			}
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

func view_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !view_open {
			c.JSON(500, gin.H{"message": "The server is not supported \"view\""})
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

func getFilesLists(path, Request_URL_Path string) [][]itemField {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	res := make([][]itemField, 2)
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			res[0] = append(res[0], itemField{file.Name(), Request_URL_Path + "/" + file.Name()})
		} else {
			res[1] = append(res[1], itemField{file.Name(), "/dl/n/" + path + file.Name()})
		}
	}
	return res
}
