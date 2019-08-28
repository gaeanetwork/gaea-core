package shareddata

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/gaeanetwork/gaea-core/protos/tee"
	"github.com/gaeanetwork/gaea-core/tee/mock"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	owner           = sha256.Sum256([]byte("buddingleader"))
	contentForTests = &service.Content{
		Data:        "data",
		Hash:        "dataHash",
		Description: "I'm a good boy.",
		Owner:       hex.EncodeToString(owner[:]),
	}
	sharedDataForTests = &tee.SharedData{
		Id:          "data_id",
		Ciphertext:  contentForTests.Data,
		Hash:        contentForTests.Hash,
		Description: contentForTests.Description,
		Owner:       contentForTests.Owner,
	}
)

func Test_Upload(t *testing.T) {
	sdService := NewSharedDataService()
	uploadReq := &service.UploadRequest{
		Content: contentForTests,
	}

	data, err := proto.Marshal(sharedDataForTests)
	if err != nil {
		t.Fatal(err)
	}
	sdService.scService = &mock.TeeChaincodeService{Result: data}
	uploadResp, err := sdService.Upload(nil, uploadReq)
	assert.NoError(t, err)
	assert.Equal(t, sharedDataForTests.Owner, uploadResp.Data.Owner)

	// Invalid request - No content
	invalidReq := &service.UploadRequest{}
	_, err = sdService.Upload(nil, invalidReq)
	assert.Error(t, err, "the content of upload request is non-empty")

	// Invalid request - Owner size
	content1 := *contentForTests
	content1.Owner = "buddingleader"
	invalidReq = &service.UploadRequest{
		Content: &content1,
	}
	_, err = sdService.Upload(nil, invalidReq)
	assert.Error(t, err, "invalid owner size, should be 64")

	// Chaincode Error
	sdService.scService = &mock.TeeChaincodeService{Error: errors.New("occured error")}
	_, err = sdService.Upload(nil, uploadReq)
	assert.Error(t, err, "occured error")

	// Chaincode data Error
	sdService.scService = &mock.TeeChaincodeService{Result: []byte("abc")}
	_, err = sdService.Upload(nil, uploadReq)
	assert.Error(t, err, "failed to unmarshal shared data")
}
