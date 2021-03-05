package main

import (
	"Libs/Libs"
	"archive/zip"
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
	ls_open     = false
	dls_open    = false
	upload_open = false
	zip_open    = false
	Version     = "Sep 27,2020 Su."
)

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
					os.Mkdir("./upload", 0644)
				}
			case "zip", "-zip":
				fmt.Println(" -  zip mode on.")
				zip_open = true
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
		return fmt.Sprintf("[%s %d %s] \n{%s} < %s >\n\"%s %s \"%s\" %s\"\n\n",
			param.Method,
			param.StatusCode,
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Path,
			param.Request.Proto,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": Version})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "It's the file download server." +
			"You can use the path to download the file on the machine."})
	})
	r.GET("/upload", func(c *gin.Context) {
		if upload_open {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, `<form action="/upload" method="post" enctype="multipart/form-data">
    Files: <input type="file" name="files" multiple><br><br>
    <input type="submit" value="Submit">
</form>`)
		} else {
			c.JSON(200, gin.H{"message": "The server is not supported \"upload\""})
		}
	})

	r.POST("/upload", func(c *gin.Context) {
		if upload_open == true {
			form, err := c.MultipartForm()
			if err != nil {
				c.JSON(500, gin.H{"message": "got file error."})
			}
			files := form.File["files"]
			for _, file := range files {
				if !Libs.LibsXExists("upload") {
					os.Mkdir("upload", 0644)
				}
				folder := fmt.Sprintf("upload/from[%s]", c.ClientIP())
				if !Libs.LibsXExists(folder) {
					os.Mkdir(folder, 0644)
				}
				c.SaveUploadedFile(file, fmt.Sprintf("%s/%s.dat", folder, file.Filename))
			}
			c.JSON(200, gin.H{"message": "OK"})
		} else {
			c.JSON(501, gin.H{"message": "The server is not supported \"upload\""})
		}
	})
	r.GET("/download/*path", func(c *gin.Context) {
		c.File(c.Param("path")[1:])
	})
	r.GET("/ls/*path", func(c *gin.Context) {
		if ls_open {
			path := c.Param("path")[1:]
			c.Header("Content-Type", "text/html; charset=utf-8")
			files := getFilesLists(path, c)
			c.String(200, files)
		} else {
			c.JSON(501, gin.H{"message": "The server is not supported \"list files\""})
		}
	})
	r.GET("/dls/*path", func(c *gin.Context) {
		if dls_open {
			path := c.Param("path")[1:]
			c.Header("Content-Type", "text/html; charset=utf-8")
			files := getFilesLists(path, c)
			c.String(200, files)
		} else {
			c.JSON(501, gin.H{"message": "The server is not supported \"list download files\""})
		}
	})
	r.GET("/zip/*path", func(c *gin.Context) {
		if zip_open {
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
		} else {
			c.JSON(501, gin.H{"message": "The server is not supported."})
		}
	})
	r.Run(fmt.Sprintf("%s:%s", ip, port))
}

func getFilesLists(path string, c *gin.Context) string {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	skillFolder := path
	result := "<html>Items:<br><br>"
	fs, ds := "<p> File:<br>", "<p> Dir:<br>"
	files, _ := ioutil.ReadDir(skillFolder)
	if ls_open && !dls_open {
		for _, file := range files {
			if file.IsDir() {
				ds += "  " + file.Name() + "<br>"
			} else {
				fs += "  " + file.Name() + "<br>"
			}
		}
	} else {
		for _, file := range files {
			if file.IsDir() {
				ds += "  <a href='" + c.Request.URL.Path + "/" + file.Name() + "'>" + file.Name() + "</a><br>"
			} else {
				fs += "  <a href='/download/" + path + file.Name() + "'>" + file.Name() + "</a><br>"
			}
		}
	}
	ds += "</p>"
	fs += "</p>"
	result = result + ds + "<br>" + fs + "</html>"
	return result
}
