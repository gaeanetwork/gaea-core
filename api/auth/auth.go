package auth

import (
	"net/http"
	"time"

	"github.com/gaeanetwork/gaea-core/protos/user"
	"github.com/gaeanetwork/gaea-core/services"
	"github.com/gin-gonic/gin"
)

var (
	expireAfter = time.Hour
)

// RegisterAPI register to login related api to gin
func RegisterAPI(auth *gin.RouterGroup) {
	auth.POST("/register", Register)
	auth.POST("/login", Login)
	auth.POST("/logout", Login)
}

// Register a x user
func Register(c *gin.Context) {
	var req user.RegisterRequest
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

	client := user.NewUserServiceClient(conn)
	resp, err := client.Register(nil, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": resp.User})
}

// Login a x user
func Login(c *gin.Context) {
	var req user.LoginRequest
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

	client := user.NewUserServiceClient(conn)
	resp, err := client.Login(nil, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := SignJWTToken(resp.User.UserName, expireAfter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": resp.User, "access_token": tokenString})
}
