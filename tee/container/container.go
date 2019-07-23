package container

import (
	"crypto"

	"github.com/gaeanetwork/gaea-core/tee/container/dev"
	"github.com/gaeanetwork/gaea-core/tee/container/docker"
	"github.com/pkg/errors"
)

// Container to do trusted execution
type Container interface {
	// Get the container public key
	GetPublicKey() crypto.PublicKey

	// Upload algorithm and data to container
	Upload(algorithm []byte, dataList [][]byte) error

	// Verify algorithm and data integrity
	Verify(algorithmHash string, dataHash []string) error

	// Execute the container
	Execute() ([]byte, error)

	// Destroy the container
	Destroy() error
}

// GetContainer get a container by type
func GetContainer(containerType Type) (Container, error) {
	switch containerType {
	case Dev:
		return dev.Create()
	case Docker:
		return docker.Create()
	default:
		return nil, errors.New("Not implemented")
	}
}

// Type for how to use the trusted execution environment
type Type int

// Docker is a folder for using container inside a chaincode container.
//
// Sibling is a sibling docker container for using container inside a chaincode container.
// It needs to update core.yaml and dockercontroller.go to bind docker.sock to mounts.
//
// SGX is a Hardware chip CPU.
const (
	Dev Type = iota
	Docker
	SGX
)
