package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gaeanetwork/gaea-core/common/config"
	"github.com/gaeanetwork/gaea-core/tee/server"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_TransferFile(t *testing.T) {
	c := gin.Default()
	router := c.Group("/testing")
	RegisterTransmissionAPI(router)
	go server.NewTeeServer(config.GRPCAddr).Start()

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
	result := w.Body.String()
	assert.Contains(t, result, "file_id")

	// Download
	fileID := getFileID(result)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/testing/"+fileID, nil)
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "world!", w.Body.String())
}

func Test_UploadFile_Error(t *testing.T) {
	c := gin.Default()
	router := c.Group("/testing")
	RegisterTransmissionAPI(router)
	config.GRPCAddr = "localhost:1"

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

	// Invalid upload request - error method
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/testing", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Invalid upload request - error url
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/testings", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Invalid upload request - error file field
	var body1 bytes.Buffer
	writer1 := multipart.NewWriter(&body1)
	part1, _ := writer1.CreateFormFile("file1", filepath.Base(filePath))
	io.Copy(part1, file)
	writer1.Close()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/testing", &body1)
	req.Header.Set("Content-Type", writer1.FormDataContentType())
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "http: no such file")

	// Internal Error - server crash
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/testing", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "connect: connection refused")
}

func Test_DownloadFile_Error(t *testing.T) {
	c := gin.Default()
	router := c.Group("/testing")
	RegisterTransmissionAPI(router)
	config.GRPCAddr = "localhost:1"

	// Invalid download request - error method
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/testing/hello", nil)
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Invalid download request - error url
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/testing/", nil)
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Invalid download request - error file_id field
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/testing/asdf", nil)
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid file id size")

	// Internal Error - server crash
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/testing/711e9609339e92b03ddc0a211827dba421f38f9ed8b9d806e1ffdd8c15ffa03d", nil)
	c.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "connect: connection refused")
}

// result := "{"file_id":"711e9609339e92b03ddc0a211827dba421f38f9ed8b9d806e1ffdd8c15ffa03d"}"
func getFileID(result string) string {
	res := struct {
		FileID string `json:"file_id"`
	}{}
	if err := json.Unmarshal([]byte(result), &res); err != nil {
		panic(err)
	}

	return res.FileID
}
