package tee

import (
	"net/http"

	"github.com/gaeanetwork/gaea-core/common/glog"
	"github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/gaeanetwork/gaea-core/services"
	"github.com/gin-gonic/gin"
)

var (
	logger = glog.MustGetLoggerWithNamed("api.tee")
)

// RegisterSharedDataAPI register shared data apis to gin server
func RegisterSharedDataAPI(apiRG *gin.RouterGroup) {
	dataRG := apiRG.Group("/tee/data")
	dataRG.POST("", upload)
}

/**
Upload used to upload shared data for users. After the data is uploaded, once someone else requests to query this data,
the user will be notified and can authorize or reject the request.
// @Param	data			formData 	string		true	"Encrypted ciphertext used to share data, usually encrypted with a private key. Of course, you can also not encrypt, upload data plaintext, such as data addresses."
// @Param	hash			formData 	string		true	"A summary of the data shared by the user. It is generally calculated using SM3/SHA-256/MD5. SM3 encryption is currently recommended."
// @Param	description		formData 	string		true	"A data description of the user's shared data. Often used to explain the basics of data or what it can be used for."
// @Param	owner			formData 	string		true	"Data owner for user shared data. Generally use owner public key."
// @Param	signatures		formData 	[]string	false	"Signature of the data summary by the user's private key."
// @router / [post]
*/
func upload(c *gin.Context) {
	var req service.UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Sugar().Debugf("Received upload shared data request, req: [%s]", req.String())

	conn, err := services.GetGRPCConnection()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := service.NewSharedDataClient(conn)
	resp, err := client.Upload(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Sugar().Debugf("Received upload shared data response, resp: [%s]", resp.String())

	c.JSON(http.StatusOK, gin.H{"data": resp.Data})
}
