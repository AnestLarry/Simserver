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
	ip := "0.0.0.0"
	port := "5000"
	if len(os.Args) > 1 {
		for i := 0; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "h", "-h", "help":
				fmt.Println(
					"Tips:\n" +
						" h  - show this help\n" +
						" v  - get version\n" +
						"Mode:\n" +
						" ls  - open ls function\n" +
						" dls  - add download links with the ls function's list\n" +
						"Args:\n" +
						" p / port  - use the port\n" +
						" ip  - use the ip.")
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
			case "ip", "-ip":
				p = true
				i++
				ip = os.Args[i]
			case "p", "-p", "port", "-port":
				p = true
				i++
				port = os.Args[i]
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
	e.Logger.Fatal(e.Start(ip + ":" + port))
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
