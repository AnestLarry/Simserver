package viewGroup

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

var (
	Enable = false
)

func view_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !Enable {
			c.JSON(500, gin.H{"message": "The server is not supported \"view\""})
			c.Abort()
		}
	}
}
func View_routerGroup_init(View_routerGroup *gin.Engine, viewFiles embed.FS) {
	routerPage := View_routerGroup.Group("/view")
	routerApi := View_routerGroup.Group("api/view")
	routerPage.Use(view_middleware())
	routerApi.Use(view_middleware())
	views := []string{}
	de, _ := viewFiles.ReadDir("view")
	for _, e := range de {
		if e.IsDir() {
			views = append(views, e.Name())
		}
	}
	routerApi.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Ok", "views": views})
	})

	viewPanelByte, _ := viewFiles.ReadFile("view/index.html")
	routerPage.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html", viewPanelByte)
	})

	for _, plugin := range views {
		view, err := fs.Sub(viewFiles, fmt.Sprintf("view/%s", plugin))
		if err != nil {
			panic(err)
		}
		routerPage.StaticFS(fmt.Sprintf("/%s", plugin), http.FS(view))
	}
}
