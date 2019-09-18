package auth

import (
	"net/http"
	"time"

	"github.com/gaeanetwork/gaea-core/common/glog"
	"github.com/gaeanetwork/gaea-core/protos/user"
	"github.com/gaeanetwork/gaea-core/services"
	"github.com/gin-gonic/gin"
)

var (
	expireAfter = time.Hour
	logger      = glog.MustGetLoggerWithNamed("api")
)

// RegisterAPI register to login related api to gin
func RegisterAPI(apiRG *gin.RouterGroup) {
	apiRG.POST("/register", Register)
	apiRG.POST("/login", Login)
	apiRG.POST("/logout", Login)
}

// Register a x user
func Register(c *gin.Context) {
	var req user.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Sugar().Debugf("Received register request: [%s]", req.String())

	conn, err := services.GetGRPCConnection()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := user.NewUserServiceClient(conn)
	resp, err := client.Register(c, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Sugar().Debugf("Received register response: [%s]", resp.String())

	c.JSON(http.StatusOK, gin.H{"user": resp.User})
}

// Login a x user
func Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Sugar().Debugf("Received login request: [%s]", req.String())

	conn, err := services.GetGRPCConnection()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := user.NewUserServiceClient(conn)
	resp, err := client.Login(c, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Sugar().Debugf("Received login response: [%s]", req.String())

	tokenString, err := SignJWTToken(resp.User.UserName, expireAfter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": resp.User, "access_token": tokenString})
}
