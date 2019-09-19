package tee

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/gaeanetwork/gaea-core/protos/tee"
	"github.com/gaeanetwork/gaea-core/tee/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	owner              = sha256.Sum256([]byte("buddingleader"))
	sharedDataForTests = &tee.SharedData{
		Id:          "data_id",
		Ciphertext:  "data",
		Hash:        "dataHash",
		Description: "I'm a good boy.",
		Owner:       hex.EncodeToString(owner[:]),
	}
)

func Test_Upload(t *testing.T) {
	sdService := NewSharedDataService()
	data, _ := json.Marshal(sharedDataForTests)
	sdService.scService = &mock.TeeChaincodeService{Result: data}

	uploadReq := &service.UploadRequest{
		Hash:        sharedDataForTests.Hash,
		Description: sharedDataForTests.Description,
		Owner:       sharedDataForTests.Owner,
	}

	uploadResp, err := sdService.Upload(nil, uploadReq)
	assert.NoError(t, err)
	assert.Equal(t, sharedDataForTests.Hash, uploadResp.Data.Hash)
	assert.Equal(t, sharedDataForTests.Description, uploadResp.Data.Description)
	assert.Equal(t, sharedDataForTests.Owner, uploadResp.Data.Owner)

	// Invalid request - No hash
	invalidReq := &service.UploadRequest{
		// Hash:        sharedDataForTests.Hash,
		Description: sharedDataForTests.Description,
		Owner:       sharedDataForTests.Owner,
	}
	_, err = sdService.Upload(nil, invalidReq)
	assert.Contains(t, fmt.Sprintf("%v", err), "the hash of upload request is non-empty")

	// Invalid request - No Description
	invalidReq = &service.UploadRequest{
		Hash: sharedDataForTests.Hash,
		// Description: sharedDataForTests.Description,
		Owner: sharedDataForTests.Owner,
	}
	_, err = sdService.Upload(nil, invalidReq)
	assert.Contains(t, fmt.Sprintf("%v", err), "the description of upload request is non-empty")

	// Invalid request - No Owner
	invalidReq = &service.UploadRequest{
		Hash:        sharedDataForTests.Hash,
		Description: sharedDataForTests.Description,
		// Owner:       sharedDataForTests.Owner,
	}
	_, err = sdService.Upload(nil, invalidReq)
	assert.Contains(t, fmt.Sprintf("%v", err), "the owner of upload request is non-empty")

	// Chaincode Error
	sdService.scService = &mock.TeeChaincodeService{Error: errors.New("occured error")}
	_, err = sdService.Upload(nil, uploadReq)
	assert.Error(t, err, "occured error")

	// Chaincode data Error
	sdService.scService = &mock.TeeChaincodeService{Result: []byte("abc")}
	_, err = sdService.Upload(nil, uploadReq)
	assert.Error(t, err, "failed to unmarshal shared data")
}
