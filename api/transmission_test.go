package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_UploadFile(t *testing.T) {
	c := gin.Default()
	router := c.Group("/testing")
	RegisterAPI(router)

	// write tmp file and make request body
	filePath, data := filepath.Join(os.TempDir(), "hello"), []byte("world!")
	ioutil.WriteFile(filePath, data, 0755)
	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("file", filepath.Base(filePath))
	io.Copy(part, file)
	writer.Close()

	// Upload
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/testing", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "'hello' uploaded!", w.Body.String())

	// Invalid request - error method
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/testing", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Invalid request - error url
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/testings", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_DownloadFile(t *testing.T) {
	c := gin.Default()
	router := c.Group("/testing")
	RegisterAPI(router)

	// download
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/testing/hello", nil)
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "'hello' downloaded!", w.Body.String())

	// Invalid request - error method
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/testing/hello", nil)
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Invalid request - error url
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/testing/", nil)
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
