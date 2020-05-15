package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"libs"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	ls_open     = false
	dls_open    = false
	upload_open = false
	Version     = "May 15,2020 Fri."
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
				upload_open = true
				if libs.LibsXExists("./upload") {
					if !libs.LibsXIsDir("./upload") {
						fmt.Println("upload is not a folder!")
						os.Exit(0)
					}
				} else {
					os.Mkdir("./upload", 0644)
				}
			case "RSUN", "-RSUN":
				if libs.LibsXExists("./upload") {
					if libs.LibsXIsDir("./upload") {
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
				fmt.Println("Do you mean \"-h\" ?")
				os.Exit(0)
			}
		}
	}
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "(${method}, ${status}), time = ${time_rfc3339}, uri = [${uri}], remote ip = <${remote_ip}>\n",
	}))
	e.HideBanner = true
	e.GET("/", index)
	e.GET("/download/*", getfile)
	e.GET("/ls/*", getls)
	e.GET("/dls/*", getdls)
	e.GET("/version", getversion)
	e.GET("/upload", uploadpage)
	e.POST("/upload", upload)
	e.Logger.Fatal(e.Start(ip + ":" + port))
}

func getversion(c echo.Context) error {
	return c.String(http.StatusOK, Version)
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "It's the file download server.\nYou can use the path to download the file on the machine.")
}

func uploadpage(c echo.Context) error {
	return c.HTML(http.StatusOK, `<form action="/upload" method="post" enctype="multipart/form-data">
    Files: <input type="file" name="files" multiple><br><br>
    <input type="submit" value="Submit">
</form>`)
}

func getfile(c echo.Context) error {
	return c.File(c.Request().URL.Path[10:])
}

func getls(c echo.Context) error {
	if ls_open {
		files := getfileslists(c.Request().URL.Path[4:], c)
		return c.HTML(http.StatusOK, files)
	} else {
		return c.String(http.StatusNotImplemented, "Error 501")
	}
}

func getdls(c echo.Context) error {
	if dls_open {
		files := getfileslists(c.Request().URL.Path[5:], c)
		return c.HTML(http.StatusOK, files)
	} else {
		return c.String(http.StatusNotImplemented, "Error 501")
	}
}

func getfileslists(path string, c echo.Context) string {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	skillfolder := path
	result := "<html>Items:<br><br>"
	fs, ds := "<p> File:<br>", "<p> Dir:<br>"
	files, _ := ioutil.ReadDir(skillfolder)
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
				ds += "  <a href='" + c.Request().URL.Path + "/" + file.Name() + "'>" + file.Name() + "</a><br>"
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

func upload(c echo.Context) error {
	if upload_open {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		files := form.File["files"]

		for _, file := range files {
			// Source
			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			// Destination
			dst, err := os.Create(fmt.Sprintf("./upload/%s.dat", file.Filename))
			if err != nil {
				return err
			}
			defer dst.Close()

			// Copy
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}
		}
		return c.String(http.StatusOK, "OK")
	} else {
		return c.String(http.StatusNotImplemented, "Error 501")
	}
}
