package tee

import (
	"context"
	"encoding/json"

	"github.com/gaeanetwork/gaea-core/protos/service"
	pb "github.com/gaeanetwork/gaea-core/protos/tee"
	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/pkg/errors"
)

// Service is used to do some shared data tasks
type Service struct {
	scService smartcontract.Service
}

// NewSharedDataService create a shared data service
func NewSharedDataService() *Service {
	smartcontractService, err := factory.GetSmartContractService(tee.ImplementPlatform)
	if err != nil {
		panic(err)
	}

	return &Service{smartcontractService}
}

// Upload used to upload shared data for users. After the data is uploaded, once someone else
// requests to query this data,
func (s *Service) Upload(ctx context.Context, req *service.UploadRequest) (*service.UploadResponse, error) {
	if req.Hash == "" {
		return nil, errors.New("the hash of upload request is non-empty")
	}
	if req.Description == "" {
		return nil, errors.New("the description of upload request is non-empty")
	}
	if req.Owner == "" {
		return nil, errors.New("the owner of upload request is non-empty")
	}

	// TODO - SIGNATURE NEEDS TO BE RECODED WITH THE CHAINCODE
	// result, err := s.scService.Invoke(tee.ChaincodeName, []string{tee.MethodUpload, req.Data,
	// 	req.Hash, req.Description, req.Owner, req.Hash, req.Signatures.String()})
	result, err := s.scService.Invoke(tee.ChaincodeName, []string{tee.MethodUpload, req.Data,
		req.Hash, req.Description, req.Owner})
	if err != nil {
		return nil, errors.Wrap(err, "failed to invoker chaincode upload function")
	}

	var data pb.SharedData
	if err = json.Unmarshal(result, &data); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal shared data")
	}

	return &service.UploadResponse{Data: &data}, nil
}
