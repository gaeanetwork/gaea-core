package container

import "crypto"

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
