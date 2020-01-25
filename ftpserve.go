package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	ls_open  = false
	dls_open = false
	Version  = "Jan 24,2020 Fri."
)

func main() {
	p := false
	if len(os.Args) > 1 {
		for _, v := range os.Args {
			switch v {
			case "h", "-h", "help":
				fmt.Println(" h  - show this help\n v  - get version\n ls  - open ls function\n dls  - add download links with the ls function's list.")
				os.Exit(0)
			case "v", "-v", "version":
				fmt.Println(Version)
				os.Exit(0)
			case "ls", "-ls":
				fmt.Println(" -  ls mode on.")
				p = true
				ls_open = true
			case "dls", "-dls":
				fmt.Println(" -  dls mode on.")
				p = true
				dls_open = true
			}
		}
		if !p {
			fmt.Println("Do you mean \"-h\" ?")
			os.Exit(0)
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
	e.Logger.Fatal(e.Start(":5000"))
}

func getversion(c echo.Context) error {
	return c.String(http.StatusOK, Version)
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "It's the file download server.\nYou can use the path to download the file on the machine.")
}

func getfile(c echo.Context) error {
	return c.File(c.Request().URL.Path[10:])
}

func getls(c echo.Context) error {
	if ls_open {
		files := getfileslists(c.Request().URL.Path[4:])
		return c.HTML(http.StatusOK, files)
	} else {
		return c.String(http.StatusNotImplemented, "Error 501")
	}
}

func getdls(c echo.Context) error {
	if dls_open {
		files := getfileslists(c.Request().URL.Path[5:])
		return c.HTML(http.StatusOK, files)
	} else {
		return c.String(http.StatusNotImplemented, "Error 501")
	}
}

func getfileslists(path string) string {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	skillfolder := path
	result := "<html>Items:<br><br>"
	fs, ds := "<p> File:<br>", "<p> Dir:\n"
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
				ds += "  " + file.Name() + "<br>"
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

// func Exists(path string) bool {
// 	_, err := os.Stat(path)
// 	if err != nil {
// 		if os.IsExist(err) {
// 			return true
// 		}
// 		return false
// 	}
// 	return true
// }
