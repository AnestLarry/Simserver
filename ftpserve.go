package main

import (
	"Simserver/Libs"
	"Simserver/chatBoard"
	argsConfig "Simserver/config"
	"Simserver/downloadGroup"
	"Simserver/uploadGroup"
	"Simserver/viewGroup"
	"embed"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	log_file_open = false
	https_open    = false
	login         = struct {
		open     bool
		account  string
		password string
	}{}
	Version = "Sum, 2024"
)

var (
	//go:embed view
	viewFiles                    embed.FS
	ip, port, pem_file, key_file = "0.0.0.0", "5000", "", ""
)

func main() {
	parseArgs()
	if uploadGroup.Upload_open {
		if Libs.LibsXExists("./upload") {
			if !Libs.LibsXIsDir("./upload") {
				fmt.Println("upload is not a folder!")
				os.Exit(0)
			}
		} else {
			os.Mkdir("./upload", 0764)
		}
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		r := gin.H{"method": param.Method, "StatusCode": param.StatusCode, "ClientIP": param.ClientIP,
			"TimeStamp": param.TimeStamp.Format(time.RFC1123), "Path": param.Path, "Request.Proto": param.Request.Proto,
			"Latency": param.Latency, "User-Agent": param.Request.UserAgent(), "ErrorMessage": param.ErrorMessage}
		return fmt.Sprintf("%+v\n", r)
	}))
	r.Use(gin.Recovery())
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": Version})
	})
	if viewGroup.View_open {
		r.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusPermanentRedirect, "/view/")
		})
	} else {
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": `It's a file downloadGroup server. You can transfer the file with the machine.`})
		})
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "404 Not Found"})
	})
	routerGroup_init(r)
	fmt.Printf("%s\n%s:%s\n", strings.Repeat("-", 15), ip, port)
	var err error
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
	uploadGroup.Upload_open = acs.Upload
	chatBoard.ChatBoard_open = acs.ChatBoard
	viewGroup.View_open = acs.View
	log_file_open = acs.Log
	login.open = acs.Login.Open
	login.account = acs.Login.Account
	login.password = acs.Login.Password
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
	fmt.Printf("ls:%v, view:%v, zip:%v, downCode:%v\nupload:%v, chatBoard:%v, \nlog:%v, https:%v\n",
		acs.Ls, acs.View, acs.Zip, acs.DownloadCode, acs.Upload, acs.ChatBoard, acs.Log, https_open)
}

func routerGroup_init(r *gin.Engine) {
	// global init
	router_init(r)
	// Uploader routerGroup
	uploadGroup.Upload_routerGroup_init(r)
	// Downloader routerGroup
	downloadGroup.Downloader_routerGroup_init(r)
	// View routerGroup
	viewGroup.View_routerGroup_init(r, viewFiles)
	// ChatBoard routerGroup
	chatBoard.ChatBoard_routerGroup_init(r)
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
	folders, _ := os.ReadDir("./upload")
	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}
		files, _ := os.ReadDir(fmt.Sprintf("./upload/%s", folder.Name()))
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

func parseArgs() {
	httpsArg := ""
	loginArg := ""
	flag.Func("version", "show the version", func(s string) error {
		fmt.Println(Version)
		os.Exit(0)
		return nil
	})
	flag.BoolVar(&downloadGroup.Ls_open, "ls", false, "open ls mode")
	flag.StringVar(&ip, "ip", "0.0.0.0", "set the ip listen")
	flag.StringVar(&port, "port", "5000", "set the port listen")
	flag.BoolVar(&uploadGroup.Upload_open, "upload", false, "open upload mode")
	flag.BoolVar(&downloadGroup.Zip_open, "zip", false, "open zip mode")
	flag.BoolVar(&downloadGroup.DownloadCode_open, "dC", false, "open download_code mode")
	flag.BoolVar(&viewGroup.View_open, "view", false, "open view mode")
	flag.BoolVar(&chatBoard.ChatBoard_open, "chatBoard", false, "open chatBoard mode")
	flag.Func("log", "open log file", func(s string) error {
		log_file_open = true
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
		return nil
	})
	flag.StringVar(&httpsArg, "https", "", "use HTTPS with cer and key\nexample: \"cer.cer key.pvk\"")
	flag.Func("RFN", "restore files' name", func(s string) error {
		restoreFileName()
		os.Exit(0)
		return nil
	})
	flag.Func("config", "use 'config.json' args", func(s string) error {
		loadConfigFromArgsConfigStruct(argsConfig.ArgConfigInit())
		return nil
	})
	flag.StringVar(&loginArg, "login", "", "add account password auth for all resource.\nexample: \"admin:admin\"")
	flag.Parse()
	if httpsArg != "" {
		https := strings.Split(httpsArg, " ")
		https_open = true
		pem_file, key_file = https[0], https[1]
	}
	if loginArg != "" {
		fmt.Printf("login [%s]\n", loginArg)
		loginArgs := strings.Split(loginArg, ":")
		login.open = true
		login.account, login.password = loginArgs[0], loginArgs[1]
	}
}

func router_init(r *gin.Engine) {
	if login.open {
		r.Use(gin.BasicAuth(gin.Accounts{login.account: login.password}))
	}
}
