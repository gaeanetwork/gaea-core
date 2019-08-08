package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	pb "github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/gaeanetwork/gaea-core/services/transmission"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	// TODO - read in config
	serverAddr = "localhost:12315"
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
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data bytes.Buffer
	if _, err = io.Copy(&data, file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO - packet it in another file
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := pb.NewTransmissionClient(conn)
	resp, err := client.UploadFile(c, &pb.UploadFileRequest{Data: data.Bytes()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded! file_id: %s", header.Filename, resp.GetFileId()))
}

/**
download a file from server

path: /files/{file_id} [GET]
*/
func downloadFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if idSize := len(fileID); idSize != transmission.MaxFileIDSize {
		errmsg := fmt.Sprintf("invalid file id size, should be %d, file_id: %s, size: %d",
			transmission.MaxFileIDSize, fileID, idSize)
		c.JSON(http.StatusBadRequest, gin.H{"error": errmsg})
		return
	}

	// TODO - packet it in another file
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := pb.NewTransmissionClient(conn)
	resp, err := client.DownloadFile(c, &pb.DownloadFileRequest{FileId: fileID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileID))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(resp.Data)))
	c.Writer.Write(resp.Data)
}
