package viewGroup

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

var (
	View_open = false
)

func view_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !View_open {
			c.JSON(500, gin.H{"message": "The server is not supported \"view\""})
			c.Abort()
		}
	}
}
func View_routerGroup_init(View_routerGroup *gin.Engine, viewFiles embed.FS) {
	routerPage := View_routerGroup.Group("/view")
	routerPage.Use(view_middleware())
	views := []string{}
	de, _ := viewFiles.ReadDir("view")
	for _, e := range de {
		if e.IsDir() {
			views = append(views, e.Name())
		}
	}
	for _, plugin := range views {
		view, err := fs.Sub(viewFiles, fmt.Sprintf("view/%s", plugin))
		if err != nil {
			panic(err)
		}
		routerPage.StaticFS(fmt.Sprintf("/%s", plugin), http.FS(view))
	}
}
