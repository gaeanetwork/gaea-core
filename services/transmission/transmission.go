package transmission

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gaeanetwork/gaea-core/common"
	pb "github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/pkg/errors"
)

const (
	// MaxFileIDSize Is the maximum file id size. The File calculated by hash is 32 bits,
	// and the file id is a hex string, so the file id size is 32 * 8 / 4 = 64.
	MaxFileIDSize = 64
)

// Service is used to do some transmission tasks
type Service struct {
	location    string
	maxFileSize int
}

// NewTransmissionService create a transmission service
func NewTransmissionService() *Service {
	// TODO - read in config
	location, maxFileSize := "/tmp/data/files/", 4*1024*1024
	return &Service{location, maxFileSize}
}

// UploadFile is used to process files uploaded by the client.
func (s *Service) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	fileSize := len(req.Data)
	if fileSize == 0 {
		return nil, errors.New("file size is zero")
	}

	if fileSize > s.maxFileSize {
		return nil, errors.Errorf("file size overflow, size: %d", fileSize)
	}

	if !common.FileOrFolderExists(s.location) {
		if err := os.MkdirAll(s.location, 0755); err != nil {
			return nil, errors.Wrapf(err, "failed to mkdir all, location: %s", s.location)
		}
	}

	hash := sha256.Sum256(req.Data)
	fileID := hex.EncodeToString(hash[:])
	filePath := filepath.Join(s.location, fileID)

	return &pb.UploadFileResponse{FileId: fileID}, ioutil.WriteFile(filePath, req.Data, 0755)
}

// DownloadFile is used to process the file requested by the client to download.
func (s *Service) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	idSize := len(req.FileId)
	if idSize != MaxFileIDSize {
		return nil, errors.Errorf("invalid file id size, should be %d, fileID: %s, size: %d",
			MaxFileIDSize, req.FileId, idSize)
	}

	filePath := filepath.Join(s.location, req.FileId)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file in location, fileID: %s", req.FileId)
	}

	return &pb.DownloadFileResponse{Data: data}, nil
}
