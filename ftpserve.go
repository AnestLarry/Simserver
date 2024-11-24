package main

import (
	"Simserver/Libs"
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
	Args    = argsConfig.GetConfig()
	Version = "Aut, 2024"
)

var (
	//go:embed view
	viewFiles embed.FS
	// ip, port, pem_file, key_file = "0.0.0.0", "5000", "", ""
)

func main() {
	parseArgs()
	if uploadGroup.Args.Enable {
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
		if strings.HasPrefix(param.Path, "/view/") {
			return ""
		}
		r := gin.H{"method": param.Method, "StatusCode": param.StatusCode, "ClientIP": param.ClientIP,
			"TimeStamp": param.TimeStamp.Format(time.RFC1123), "Path": param.Path, "Request.Proto": param.Request.Proto,
			"Latency": param.Latency, "User-Agent": param.Request.UserAgent(), "ErrorMessage": param.ErrorMessage}
		return fmt.Sprintf("%+v\n", r)
	}))
	r.Use(gin.Recovery())
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": Version})
	})
	if viewGroup.Args.Enable {
		r.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusPermanentRedirect, "/view/")
		})
	} else {
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": `It's a file downloadGroup server. You can transfer the file with the machine.`})
		})
	}
	if Args.Security.Log {
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
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "404 Not Found"})
	})
	routerGroup_init(r)
	fmt.Printf("%s\n%s:%s\n", strings.Repeat("-", 15), Args.Ip, Args.Port)
	var err error
	if Args.Security.Https.Enable {
		err = r.RunTLS(fmt.Sprintf("%s:%s", Args.Ip, Args.Port), Args.Security.Https.KeyFiles[0], Args.Security.Https.KeyFiles[1])
	} else {
		err = r.Run(fmt.Sprintf("%s:%s", Args.Ip, Args.Port))
	}
	if err != nil {
		panic(err)
	}
}

func loadConfigFromArgsConfigStruct(acs argsConfig.ArgConfigStruct) {
	fmt.Printf("%+v\n\n", acs)
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
	flag.BoolVar(&downloadGroup.Args.Ls, "ls", false, "enable ls mode")
	flag.StringVar(&Args.Ip, "ip", "0.0.0.0", "set the ip listen")
	flag.StringVar(&Args.Port, "port", "5000", "set the port listen")
	flag.BoolVar(&uploadGroup.Args.Enable, "upload", false, "enable upload mode")
	flag.BoolVar(&uploadGroup.Args.UploadText, "uT", false, "enable upload text mode")
	flag.BoolVar(&uploadGroup.Args.SecureExt, "secureExt", false, "set secureExt mode(default: true)")
	flag.BoolVar(&downloadGroup.Args.Zip, "zip", false, "enable zip mode")
	flag.BoolVar(&downloadGroup.Args.DownloadCode, "dC", false, "enable download_code mode")
	flag.BoolVar(&viewGroup.Args.Enable, "view", false, "enable view mode")
	flag.BoolVar(&viewGroup.Args.ChatBoard, "chatBoard", false, "enable chatBoard mode")
	flag.BoolVar(&Args.Security.Log, "log", false, "enable log writing to file")
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
		Args.Security.Https.Enable = true
		Args.Security.Https.KeyFiles[0], Args.Security.Https.KeyFiles[1] = https[0], https[1]
	}
	if loginArg != "" {
		fmt.Printf("login [%s]\n", loginArg)
		loginArgs := strings.Split(loginArg, ":")
		Args.Security.Login.Enable = true
		Args.Security.Login.Account, Args.Security.Login.Password = loginArgs[0], loginArgs[1]
	}
}

func router_init(r *gin.Engine) {
	if Args.Security.Login.Enable {
		r.Use(gin.BasicAuth(gin.Accounts{Args.Security.Login.Account: Args.Security.Login.Password}))
	}
}
