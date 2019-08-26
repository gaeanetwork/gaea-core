package dev

import (
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gaeanetwork/gaea-core/common"
	"github.com/google/uuid"
)

// Container for development, it will create a temporary folder in your environment,
// calculate it when it's done, it will be destroyed. It is recommended to use a
// development environment.
type Container struct {
	address       string
	program       string
	args          []string
	algorithmHash string
	dataHash      []string
}

// Create a development container
func Create() (*Container, error) {
	address := filepath.Join("/tmp/teetask/container/", uuid.New().String())
	return &Container{args: make([]string, 0), address: address}, os.MkdirAll(address, 0700)
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
	c.algorithmHash = hex.EncodeToString(hash[:])

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
		c.dataHash = append(c.dataHash, hex.EncodeToString(hash[:]))

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
	if algorithmHash != c.algorithmHash {
		return fmt.Errorf("Failed to verify the algorithm hash, task: %s, container: %s", algorithmHash, c.algorithmHash)
	}

	// check data length
	if taskLength, containerLength := len(dataHash), len(c.dataHash); taskLength != containerLength {
		return fmt.Errorf("Failed to verify the data hashes, task length: %d, container length: %d", taskLength, containerLength)
	}

	if str, ok := common.ContainsStringArray(dataHash, c.dataHash); !ok {
		return fmt.Errorf("Failed to verify the data hash, task: %v doesn't includes container: %s", dataHash, str)
	}

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

// GetPublicKey for development
func (c *Container) GetPublicKey() crypto.PublicKey {
	return nil
}
