package transmission

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/gaeanetwork/gaea-core/common/config"
	pb "github.com/gaeanetwork/gaea-core/protos/service"
	"github.com/pkg/errors"
)

const (
	// StandardIDSize Is the standard id size. The data calculated by hash is 32 bits,
	// and the id is a hex string, so the id size is 32 * 8 / 4 = 64.
	StandardIDSize = 64

	// DefaultLocation is the default file storage location.
	// It should be re-read in the configuration.
	DefaultLocation = "/tmp/data/files/"
)

// Service is used to do some transmission tasks
type Service struct {
	location    string
	maxFileSize int
}

// NewTransmissionService create a transmission service
func NewTransmissionService() *Service {
	// TODO - read in config
	location, maxFileSize := DefaultLocation, config.MaxSendMsgSize
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

	location, err := handleUserID(s.location, req.UserId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to handle user id, userID: %s", req.UserId)
	}
	if !common.FileOrFolderExists(location) {
		if err := os.MkdirAll(location, 0755); err != nil {
			return nil, errors.Wrapf(err, "failed to mkdir all, location: %s", location)
		}
	}

	hash := sha256.Sum256(req.Data)
	fileID := hex.EncodeToString(hash[:])
	filePath := filepath.Join(location, fileID)

	return &pb.UploadFileResponse{FileId: fileID}, ioutil.WriteFile(filePath, req.Data, 0755)
}

// DownloadFile is used to process the file requested by the client to download.
func (s *Service) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	if fileIDSize := len(req.FileId); fileIDSize != StandardIDSize {
		return nil, errors.Errorf("invalid file id size, should be %d, fileID: %s, size: %d",
			StandardIDSize, req.FileId, fileIDSize)
	}

	location, err := handleUserID(s.location, req.UserId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to handle user id, userID: %s", req.UserId)
	}

	filePath := filepath.Join(location, req.FileId)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file in location, fileID: %s", req.FileId)
	}

	return &pb.DownloadFileResponse{Data: data}, nil
}

func handleUserID(location, userID string) (string, error) {
	userIDSize := len(userID)
	switch userIDSize {
	case 0:
		return location, nil
	case StandardIDSize:
		return filepath.Join(location, userID), nil
	default:
		return "", errors.Errorf("invalid user id size, should be %d, userID: %s, size: %d",
			StandardIDSize, userID, userIDSize)
	}
}
