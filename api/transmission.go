package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gaeanetwork/gaea-core/common/config"
	"github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/gaeanetwork/gaea-core/services"
	"github.com/gaeanetwork/gaea-core/services/transmission"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data bytes.Buffer
	if _, err = io.Copy(&data, file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithInsecure())
	dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(config.MaxRecvMsgSize),
		grpc.MaxCallSendMsgSize(config.MaxSendMsgSize)))

	conn, err := grpc.Dial(config.GRPCAddr, dialOpts...)
	// conn, err := services.GetGRPCConnection()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := service.NewTransmissionClient(conn)
	resp, err := client.UploadFile(c, &service.UploadFileRequest{Data: data.Bytes()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file_id": resp.GetFileId()})
}

/**
download a file from server

path: /files/{file_id} [GET]
*/
func downloadFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if idSize := len(fileID); idSize != transmission.StandardIDSize {
		errmsg := fmt.Sprintf("invalid file id size, should be %d, file_id: %s, size: %d",
			transmission.StandardIDSize, fileID, idSize)
		c.JSON(http.StatusBadRequest, gin.H{"error": errmsg})
		return
	}

	conn, err := services.GetGRPCConnection()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := service.NewTransmissionClient(conn)
	resp, err := client.DownloadFile(c, &service.DownloadFileRequest{FileId: fileID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	returnFile(c, url.QueryEscape(fileID), resp.Data)
}

func returnFile(c *gin.Context, filename string, fileData []byte) {
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", fileData))
	c.Writer.Write(fileData)
}
