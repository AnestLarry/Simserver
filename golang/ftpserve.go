package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/index", index)
	e.GET("/download/*", getfile)
	e.GET("/ls/*", getls)
	e.GET("/version", getversion)
	e.Logger.Fatal(e.Start(":5000"))
}

func getversion(c echo.Context) error {
	return c.String(http.StatusOK, "06/08/2019 Sat")
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "It's the file download server.<br>You can use the path to download the file on the machine.")
}

func getfile(c echo.Context) error {
	println(c.Request().URL.Path[10:])
	return c.File(c.Request().URL.Path[10:])
}

func getls(c echo.Context) error {
	files := getls_getfiles(c.Request().URL.Path[4:])
	fmt.Println(c.Request().URL.Path[4:])
	return c.String(http.StatusOK, files)
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
