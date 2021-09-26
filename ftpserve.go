package main

import (
	"Simserver/Libs"
	"Simserver/downloadGroup"
	"Simserver/uploadGroup"
	"Simserver/viewGroup"
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	log_file_open = false
	https_open    = false
	Version       = "Aut, 2021"
)

var (
	//go:embed static
	staticFiles embed.FS
	//go:embed view
	viewFiles embed.FS
)

func main() {
	ip, port := "0.0.0.0", "5000"
	pem_file, key_file := "", ""
	var err error
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "h", "-h", "help":
				fmt.Println(`Tips:
 h  - show this help
 v  - get version
Mode:
 ls  - open ls function
 dls  - add downloadGroup links with the ls function's list
 upload  - allow user upload files to host
 uploadText  - allow user fill textarea to save text in txt
 zip  - allow zip dir for downloadGroup (DANGER!)
 https  - use https with crt and key 
 log  - put log in file
 downloadCode  - use downloadGroup code to downloadGroup a group file with setting
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
				downloadGroup.Ls_open = true
			case "dls", "-dls":
				fmt.Println(" -  dls mode on.")
				downloadGroup.Dls_open = true
			case "ip", "-ip":
				i++
				ip = os.Args[i]
			case "p", "-p", "port", "-port":
				i++
				port = os.Args[i]
			case "upload", "-upload":
				fmt.Println(" -  upload mode on.")
				uploadGroup.Upload_open = true
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
				uploadGroup.Upload_text_open = true
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
				downloadGroup.Zip_open = true
			case "downloadCode", "dC", "-downloadCode", "-dC":
				fmt.Println(" -  downloadCode mode on.")
				downloadGroup.LoadDownloadCodeJson()
				downloadGroup.DownloadCode_open = true
			case "view", "-view":
				fmt.Println(" -  view mode on.")
				viewGroup.View_open = true
			case "log", "-log":
				fmt.Println(" -  log file mode on.")
				log_file_open = true
			case "https", "-https":
				https_open = true
				pem_file = os.Args[i+1]
				key_file = os.Args[i+2]
				i = i + 2
				fmt.Println(" -  https mode on.")
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
		c.JSON(200, gin.H{"message": `It's a file downloadGroup server. You can transfer the file with the machine.`})
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "404 Not Found"})
	})
	routerGroup_init(r)
	fmt.Printf("%s\n%s:%s\n", strings.Repeat("-", 15), ip, port)
	if https_open {
		err = r.RunTLS(fmt.Sprintf("%s:%s", ip, port), pem_file, key_file)
	} else {
		err = r.Run(fmt.Sprintf("%s:%s", ip, port))
	}
	if err != nil {
		panic(err)
	}
}

func routerGroup_init(r *gin.Engine) {
	// Uploader routerGroup
	uploadGroup.Upload_routerGroup_init(r.Group("/upload"), staticFiles)
	// Downloader routerGroup
	downloadGroup.Downloader_routerGroup_init(r.Group("/dl"), staticFiles, r)
	// View routerGroup
	viewGroup.View_routerGroup_init(r.Group("/view"), viewFiles)
}
