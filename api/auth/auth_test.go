package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	expireAfterForTests = 1 * time.Second
)

func getTestLoginRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/testing/login", strings.NewReader("{\"user_name\":\"AAAAAAAA\", \"password\":\"123\"}"))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func registerTestUser(c *gin.Engine) {
	expireAfter = expireAfterForTests
	req, _ := http.NewRequest(http.MethodPost, "/testing/register", strings.NewReader("{\"user_name\":\"AAAAAAAA\", \"password\":\"123\"}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		panic(w.Body.String())
	}
}
