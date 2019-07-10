package docker

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

// Container for dev
type Container struct {
	address       string
	program       string
	args          []string
	algorithmHash string
	dataHash      []string
}

// New a development container
func New() *Container {
	return &Container{args: make([]string, 0)}
}

// Create for development
func (c *Container) Create() error {
	c.address = filepath.Join("/tmp/teetask/container/", uuid.New().String())
	return os.MkdirAll(c.address, 0755)
}

// Upload for development
func (c *Container) Upload(algorithm []byte, dataList [][]byte) error {
	return nil
}

// Verify for development
func (c *Container) Verify(algorithmHash string, dataHash []string) error {
	return nil
}

// Execute for development
func (c *Container) Execute() ([]byte, error) {
	return exec.Command(c.program, c.args...).CombinedOutput()
}

// Destroy for development
func (c *Container) Destroy() error {
	return os.RemoveAll(c.address)
}
