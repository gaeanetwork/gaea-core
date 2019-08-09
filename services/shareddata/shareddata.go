package shareddata

import (
	"context"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/gaeanetwork/gaea-core/protos/service"
	pb "github.com/gaeanetwork/gaea-core/protos/tee"
	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// Service is used to do some shared data tasks
type Service struct {
}

// Upload used to upload shared data for users. After the data is uploaded, once someone else
// requests to query this data,
func (s *Service) Upload(ctx context.Context, req *service.UploadRequest) (*service.UploadResponse, error) {
	if req.Content == nil {
		return nil, errors.New("the content of upload request is non-empty")
	}

	if ownerSize := len(req.Content.Owner); ownerSize != common.StandardIDSize {
		return nil, errors.Errorf("invalid owner size, should be %d, owner: %s, size: %d",
			common.StandardIDSize, req.Content.Owner, ownerSize)
	}

	scservice, err := factory.GetSmartContractService(tee.ImplementPlatform)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get smart contract service, platform: %s", tee.ImplementPlatform)
	}

	result, err := scservice.Query(tee.ChaincodeName, []string{tee.MethodUpload, req.Content.Data,
		req.Content.Hash, req.Content.Description, req.Content.Owner, req.Hash, req.Signature.String()})
	if err != nil {
		return nil, err
	}

	var data pb.SharedData
	if err = proto.Unmarshal(result, &data); err != nil {
		return nil, err
	}

	return &service.UploadResponse{Data: &data}, nil
}
