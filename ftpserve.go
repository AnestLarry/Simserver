package main

import (
	"Simserver/Libs"
	argsConfig "Simserver/config"
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
	Version       = "Sum, 2022"
)

var (
	//go:embed static
	staticFiles embed.FS
	//go:embed view
	viewFiles                    embed.FS
	ip, port, pem_file, key_file = "0.0.0.0", "5000", "", ""
)

func main() {
	var err error
	if len(os.Args) > 1 {
	ArgsFor:
		for i := 1; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "h", "-h", "help":
				fmt.Println(`Tips:
 h  - show this help
 v  - get version
Mode:
 ls  - open ls function
 upload  - allow user upload files to host
 uploadText  - allow user fill textarea to save text in txt
 zip  - allow zip dir for downloadGroup (DANGER!)
 https  - use https with crt and key 
 log  - put log in file
 downloadCode  - use downloadGroup code to downloadGroup a group file with setting
 view  - use view in running
Args:
 p / port  - use the port
 ip  - use the ip
 config  - use 'config.json' args
Task:
 RFN  - restore files' name`)
				os.Exit(0)
			case "v", "-v", "version":
				fmt.Println(Version)
				os.Exit(0)
			case "ls", "-ls":
				fmt.Println(" -  ls mode on.")
				downloadGroup.Ls_open = true
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
			case "RFN", "-RFN":
				restoreFileName()
			case "config", "-config":
				fmt.Println("start parse 'config.json'")
				loadConfigFromArgsConfigStruct(argsConfig.ArgConfigInit())
				break ArgsFor
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
			f, err = os.OpenFile("ftps.log", os.O_APPEND, os.ModePerm)
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

func loadConfigFromArgsConfigStruct(acs argsConfig.ArgConfigStruct) {
	downloadGroup.Ls_open = acs.Ls
	downloadGroup.Zip_open = acs.Zip
	downloadGroup.DownloadCode_open = acs.DownloadCode
	uploadGroup.Upload_text_open = acs.UploadText
	uploadGroup.Upload_open = acs.Upload
	viewGroup.View_open = acs.View
	log_file_open = acs.Log
	if acs.Ip != "" {
		ip = acs.Ip
	}
	if acs.Port != "" {
		port = acs.Port
	}
	if len(acs.Https) > 0 {
		if len(acs.Https) == 2 {
			https_open = true
			pem_file = acs.Https[0]
			key_file = acs.Https[1]
		} else {
			fmt.Println("config File:\nhttps args nums error.")
		}
	}
	fmt.Printf("ls:%v, view:%v, zip:%v, downCode:%v\nupload:%v, uploadText:%v\nlog:%v, https:%v\n",
		acs.Ls, acs.View, acs.Zip, acs.DownloadCode, acs.Upload, acs.UploadText, acs.Log, https_open)
}

func routerGroup_init(r *gin.Engine) {
	// Uploader routerGroup
	uploadGroup.Upload_routerGroup_init(r, staticFiles)
	// Downloader routerGroup
	downloadGroup.Downloader_routerGroup_init(r, staticFiles)
	// View routerGroup
	viewGroup.View_routerGroup_init(r, viewFiles)
}

func restoreFileName() {
	if !(Libs.LibsXExists("./upload") && Libs.LibsXIsDir("./upload")) {
		return
	}
	fmt.Println("WARMING :  It may be fail or rewrite the same name file.\nkeyin \"y\" to continue or other to exit")
	{
		temp := ""
		fmt.Scanf("%s", &temp)
		if temp != "y" {
			os.Exit(0)
		}
	}
	folders, _ := ioutil.ReadDir("./upload")
	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}
		files, _ := ioutil.ReadDir(fmt.Sprintf("./upload/%s", folder.Name()))
		for _, file := range files {
			path := fmt.Sprintf("./upload/%s/%s", folder.Name(), file.Name())
			if !file.IsDir() && path[len(path)-4:] == "_dat" {
				os.Rename(path, path[:len(path)-4])
				fmt.Printf("[%s] reset to [%s]\n", path, path[:len(path)-4])
			}
		}
	}
	os.Exit(0)
}
