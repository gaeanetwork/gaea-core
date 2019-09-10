package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/user"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gaeanetwork/gaea-core/common/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	userName = "Budding Leader"
)

func init() {
	config.Initialize()
}

func Test_SignJWTToken(t *testing.T) {
	expireAfter := 3 * time.Second
	tokenString, err := SignJWTToken(userName, expireAfter)
	assert.NoError(t, err)
	assert.NotNil(t, tokenString)

	// Special user name
	tokenString1, err1 := SignJWTToken("@!#@%#@%$&)_(*&^%$userName", expireAfter)
	assert.NoError(t, err1)
	assert.NotNil(t, tokenString1)

	// Verify the token string
	token, err2 := jwt.Parse(tokenString, VerifyJWTToken)
	assert.NoError(t, err2)
	assert.True(t, token.Valid)

	// Expire token
	time.Sleep(expireAfter + 1*time.Second)
	token1, err3 := jwt.Parse(tokenString, VerifyJWTToken)
	assert.Equal(t, err3.Error(), "Token is expired")
	assert.False(t, token1.Valid)
}

func Test_JWTAuth(t *testing.T) {
	c := gin.Default()
	router := c.Group("/testing")
	RegisterAPI(router)

	auth := router.Group("/", JWTAuth())
	auth.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Login to get token
	go registerTestUser(c)
	w := httptest.NewRecorder()
	req := getTestLoginRequest()
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result loginResult
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	// No token can't be accessed
	w1 := httptest.NewRecorder()
	pingReqWithoutToken, _ := http.NewRequest(http.MethodGet, "/testing/ping", nil)
	c.ServeHTTP(w1, pingReqWithoutToken)
	assert.Equal(t, http.StatusUnauthorized, w1.Code)

	// Use token access
	w2 := httptest.NewRecorder()
	pingReqWithToken, _ := http.NewRequest(http.MethodGet, "/testing/ping", nil)
	pingReqWithToken.Header.Set("Authorization", "Bearer "+result.Token)
	c.ServeHTTP(w2, pingReqWithToken)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "pong", w2.Body.String())

	// Use expire token access
	time.Sleep(expireAfterForTests)
	w3 := httptest.NewRecorder()
	pingReqWithExpireToken, _ := http.NewRequest(http.MethodGet, "/testing/ping", nil)
	pingReqWithToken.Header.Set("Authorization", "Bearer "+result.Token)
	c.ServeHTTP(w3, pingReqWithExpireToken)
	assert.Equal(t, http.StatusUnauthorized, w3.Code)
}

type loginResult struct {
	Token string     `json:"access_token"`
	User  *user.User `json:"user"`
}
