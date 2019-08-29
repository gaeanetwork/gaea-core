package task

import (
	"testing"

	"github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	sdService := NewTeeTaskService()
	req := &service.CreateRequest{
		DataId: []string{
			"1",
			"2",
			"3",
		},
		AlgorithmId:   "4",
		ResultAddress: "5",
	}

	resp, err := sdService.Create(nil, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp.TaskId)

	// Invalid request - dataId empty
	req.DataId = nil
	_, err = sdService.Create(nil, req)
	assert.Error(t, err, "data id are non-empty")
}
