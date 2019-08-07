package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterAPI register apis to gin server
func RegisterAPI(rg *gin.RouterGroup) {
	rg.POST("", uploadFile)
	rg.GET("/:file_id", downloadFile)
}

/**
upload a file to server

path: /files [POST]
*/
func uploadFile(c *gin.Context) {
	fmt.Println("upload")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO - upload file to server
	// upload(file.content)
	fmt.Println("upload")
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

/**
download a file from server

path: /files/{file_id} [GET]
*/
func downloadFile(c *gin.Context) {
	fileID := c.Param("file_id")
	// TODO - download file from server
	// download(fileID)
	c.String(http.StatusOK, fmt.Sprintf("'%s' downloaded!", fileID))
}
