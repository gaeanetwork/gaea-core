package server

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"

	pb "github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/stretchr/testify/assert"
)

var (
	testData = []byte("Hello World!")
)

func Test_TransferFile_Upload(t *testing.T) {
	service := NewTransmissionService()
	uploadReq := &pb.UploadFileRequest{Data: testData}
	resp, err := service.UploadFile(nil, uploadReq)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer os.RemoveAll(service.location)

	resp1, err := service.DownloadFile(nil, &pb.DownloadFileRequest{FileId: resp.FileId})
	assert.NoError(t, err)
	assert.NotNil(t, resp1)
	assert.Equal(t, resp1.Data, testData)

	// Invalid upload - zero data
	uploadReq = &pb.UploadFileRequest{Data: nil}
	_, err = service.UploadFile(nil, uploadReq)
	assert.Error(t, err, "file size is zero")

	// Invalid upload - data overflow
	data := make([]byte, service.maxFileSize+1)
	rand.Read(data)
	uploadReq = &pb.UploadFileRequest{Data: data}
	_, err = service.UploadFile(nil, uploadReq)
	assert.Error(t, err, "file size overflow")

	// Invalid upload - mkdir permission denied
	service.location = "/data"
	uploadReq = &pb.UploadFileRequest{Data: testData}
	_, err = service.UploadFile(nil, uploadReq)
	assert.Error(t, err, "permission denied")
}

func Test_TransferFile_Download(t *testing.T) {
	service := NewTransmissionService()

	// Invalid download - id size
	_, err := service.DownloadFile(nil, &pb.DownloadFileRequest{FileId: "123"})
	assert.Error(t, err, "invalid file id size, should be 64")

	// Invalid download - not exists
	data := make([]byte, 32)
	rand.Read(data)
	_, err = service.DownloadFile(nil, &pb.DownloadFileRequest{FileId: hex.EncodeToString(data)})
	assert.Error(t, err, "no such file or directory")

}
