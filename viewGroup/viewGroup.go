package viewGroup

import (
	"embed"
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
func View_routerGroup_init(View_routerGroup *gin.RouterGroup, viewFiles embed.FS) {
	View_routerGroup.Use(view_middleware())
	views, err := fs.Sub(viewFiles, "view/h5player")
	if err != nil {
		panic(err)
	}
	View_routerGroup.StaticFS("/h5player", http.FS(views))

}
