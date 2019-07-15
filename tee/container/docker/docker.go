package docker

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gaeanetwork/gaea-core/common"
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
	if len(algorithm) == 0 {
		return fmt.Errorf("algorithm bytes is empty")
	}

	if len(dataList) == 0 {
		return fmt.Errorf("dataList bytes is empty")
	}

	// Calculate algorithm hash
	hash := sha256.Sum256(algorithm)
	c.algorithmHash = common.BytesToHex(hash[:])

	c.program = filepath.Join(c.address, "main")
	err := ioutil.WriteFile(c.program, algorithm, 0755)
	if err != nil {
		return err
	}

	c.dataHash = make([]string, 0)
	for index, data := range dataList {
		// Calculate data hash
		if hash = sha256.Sum256(data); err != nil {
			return err
		}
		c.dataHash = append(c.dataHash, common.BytesToHex(hash[:]))

		arg := filepath.Join(c.address, strconv.Itoa(index))
		if err = ioutil.WriteFile(arg, data, 0755); err != nil {
			return err
		}
		c.args = append(c.args, arg)
	}

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
