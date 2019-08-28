package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gaeanetwork/gaea-core/common/config"
	"github.com/gaeanetwork/gaea-core/tee/server"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	ownerBytes = sha256.Sum256([]byte("buddingleader"))
	owner      = hex.EncodeToString(ownerBytes[:])
)

func Test_UploadData(t *testing.T) {
	c := gin.Default()
	router := c.Group("/testing")
	RegisterSharedDataAPI(router)
	config.GRPCAddr = ":22667"
	teeServer := server.NewTeeServer(config.GRPCAddr)
	go teeServer.Start()

	reqParams := `{
		"content": {
			"data": "data",
			"hash": "dataHash",
			"description": "I'm a good boy.",
			"owner": "%s"
		}
	}`
	reqParams = fmt.Sprintf(reqParams, owner)
	req, _ := http.NewRequest(http.MethodPost, "/testing/data", strings.NewReader(reqParams))
	w := httptest.NewRecorder()
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")

	// Invalid Request
	req, _ = http.NewRequest(http.MethodPost, "/testing/data", nil)
	w = httptest.NewRecorder()
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")

	// Invalid connection
	teeServer.GracefulStop()
	req, _ = http.NewRequest(http.MethodPost, "/testing/data", strings.NewReader(reqParams))
	w = httptest.NewRecorder()
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "connection refused")
}
