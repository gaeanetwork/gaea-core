package container

// Container to do trusted execution
type Container interface {
	// Create a container
	Create() error

	// Upload algorithm and data to container
	Upload(algorithm []byte, dataList [][]byte) error

	// Verify algorithm and data integrity
	Verify(algorithmHash string, dataHash []string) error

	// Execute the container
	Execute() ([]byte, error)

	// Destroy the container
	Destroy() error
}
