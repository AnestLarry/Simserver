package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

var (
	ls_open = false
)

func main() {
	p := false
	if len(os.Args) > 1 {
		for _, v := range os.Args {
			switch v {
			case "h", "-h", "help":
				fmt.Println("h  - show this help\nv  - get version\nls  - open ls function.")
				os.Exit(0)
			case "v", "-v", "version":
				fmt.Println("11/30/2019 Sat")
				os.Exit(0)
			case "ls", "-ls":
				p = true
				ls_open = true
			}
		}
		if !p {
			fmt.Println("Do you mean \"-h\" ?")
			os.Exit(0)
		}
	}
	e := echo.New()
	e.GET("/", index)
	e.GET("/download/*", getfile)
	e.GET("/ls/*", getls)
	e.GET("/version", getversion)
	e.Logger.Fatal(e.Start(":5000"))
}

func getversion(c echo.Context) error {
	return c.String(http.StatusOK, "11/30/2019 Sun")
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "It's the file download server.\nYou can use the path to download the file on the machine.")
}

func getfile(c echo.Context) error {
	println(c.Request().URL.Path[10:])
	return c.File(c.Request().URL.Path[10:])
}

func getls(c echo.Context) error {
	if ls_open {
		files := getls_getfiles(c.Request().URL.Path[4:])
		fmt.Println(c.Request().URL.Path[4:])
		return c.String(http.StatusOK, files)
	} else {
		return c.String(http.StatusNotImplemented, "Error 501")
	}

}
func getls_getfiles(path string) string {
	skillfolder := path
	result := "Items:\n\n"
	fs, ds := " File:\n", " Dir:\n"
	files, _ := ioutil.ReadDir(skillfolder)
	for _, file := range files {
		if file.IsDir() {
			ds += "  " + file.Name() + "\n"
		} else {
			fs += "  " + file.Name() + "\n"
		}
	}
	result = result + ds + "\n" + fs
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
