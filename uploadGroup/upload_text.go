package uploadGroup

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleTextUpload(c *gin.Context) {
	text := c.PostForm("text")
	f, err := os.Create(fmt.Sprintf("./upload/%s.txt", time.Now().Format("2006-01-02--15-04-05")))
	if err != nil || text == "" {
		fmt.Println("fail to get text.")
		c.JSON(500, gin.H{"message": "fail to get text."})
		if f != nil {
			f.Close()
		}
		return
	}
	defer f.Close()
	_, err = f.WriteString(text)
	if err != nil {
		fmt.Println("fail to save text file.")
		c.JSON(500, gin.H{"message": "fail to save text file."})
		return
	}
	c.JSON(200, gin.H{"message": "OK"})
}