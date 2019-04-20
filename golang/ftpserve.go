package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/index", index)
	e.GET("/download/*", getfile)
	e.Logger.Fatal(e.Start(":5000"))
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "It's the file download server.<br>You can use the path to download the file on the machine.")
}
func getfile(c echo.Context) error {
	println(c.Request().URL.Path[10:])
	return c.File(c.Request().URL.Path[10:])
}
