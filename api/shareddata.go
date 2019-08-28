package api

import (
	"net/http"

	"github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/gaeanetwork/gaea-core/services"
	"github.com/gin-gonic/gin"
)

// RegisterSharedDataAPI register shared data apis to gin server
func RegisterSharedDataAPI(rg *gin.RouterGroup) {
	data := rg.Group("/data")
	data.POST("", upload)
}

/**
upload a shared data to server

// @Param	ciphertext		formData 	string		true	"Encrypted ciphertext used to share data, usually encrypted with a private key. Of course, you can also not encrypt, upload data plaintext, such as data addresses."
// @Param	summary			formData 	string		true	"A summary of the data shared by the user. It is generally calculated using SM3/SHA-256/MD5. SM3 encryption is currently recommended."
// @Param	description		formData 	string		true	"A data description of the user's shared data. Often used to explain the basics of data or what it can be used for."
// @Param	owner			formData 	string		true	"Data owner for user shared data. Generally use owner public key."
// @Param	hash			formData 	string		false	"All parameters except the signature are sequentially connected to obtain a hash."
// @Param	signatures		formData 	[]string	false	"Signature of the data summary by the user's private key."
// @router / [post]
*/
func upload(c *gin.Context) {
	var req service.UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"data": resp.Data})
}
