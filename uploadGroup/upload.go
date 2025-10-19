package uploadGroup

import (
	"Simserver/Libs"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// getDecodedFileName extracts and decodes the file name from the request header.
func getDecodedFileName(c *gin.Context) (string, error) {
	x_file_name := c.Request.Header.Get("x-file-name")
	if len(x_file_name) == 0 {
		return "", fmt.Errorf("get file name failed")
	}
	x_file_name_byte, err := base64.StdEncoding.DecodeString(x_file_name)
	if err != nil {
		return "", fmt.Errorf("decode file name failed")
	}
	return string(x_file_name_byte), nil
}

// streamUpload handles the streaming file upload from a multipart reader.
func streamUpload(reader *multipart.Reader, destPath string) error {
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading multipart section: %w", err)
		}

		if part.FormName() == "file" && part.FileName() != "" {
			defer part.Close()
			out, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("failed to create destination file: %w", err)
			}
			defer out.Close()

			_, err = io.Copy(out, part)
			if err != nil {
				return fmt.Errorf("failed to save file content: %w", err)
			}
			return nil // File found and processed
		}
		part.Close()
	}
	return fmt.Errorf("file part 'file' not found in request")
}

func HandleFileUpload(c *gin.Context, secureExt bool) {
	folder := fmt.Sprintf("upload/from_%s_", strings.ReplaceAll(c.ClientIP(), ".", "_"))
	if !Libs.LibsXExists(folder) {
		os.MkdirAll(folder, 0664)
	}

	fileName, err := getDecodedFileName(c)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	destFileName := fmt.Sprintf("%s/%s", folder, fileName)
	if secureExt {
		destFileName = fmt.Sprintf("%s/%s_dat", folder, fileName)
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to get multipart reader: " + err.Error()})
		return
	}

	if err := streamUpload(reader, destFileName); err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Upload success"})
}