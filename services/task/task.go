package task

import (
	"context"
	"encoding/json"

	"github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/gaeanetwork/gaea-core/smartcontract"
	"github.com/gaeanetwork/gaea-core/smartcontract/factory"
	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/gaeanetwork/gaea-core/tee/task"
	"github.com/pkg/errors"
)

// Service is used to do some tee tasks
type Service struct {
	scService smartcontract.Service
}

// NewTeeTaskService create a tee task service
func NewTeeTaskService() *Service {
	smartcontractService, err := factory.GetSmartContractService(tee.ImplementPlatform)
	if err != nil {
		panic(err)
	}

	return &Service{smartcontractService}
}

// Create creates a tee task by sending a create transaction to the blockchain.
func (s *Service) Create(ctx context.Context, req *service.CreateRequest) (*service.CreateResponse, error) {
	if len(req.DataId) <= 0 {
		return nil, errors.Errorf("data id are non-empty")
	}

	ids, err := json.Marshal(req.DataId)
	if err != nil {
		return nil, err
	}

	result, err := s.scService.Invoke(task.ChaincodeName, []string{task.MethodCreate, string(ids), req.AlgorithmId, req.ResultAddress})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create tee task, req: %v", req)
	}

	return &service.CreateResponse{TaskId: string(result)}, nil
}
