package docker

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strconv"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gaeanetwork/gaea-core/common"
	"github.com/hyperledger/fabric/core/container/util"
	"github.com/pkg/errors"
)

// Container for docker is created by the docker container, which constructs a docker
// container to perform trusted execution as a trusted execution environment.
type Container struct {
	id            string
	client        *docker.Client
	address       string
	cmd           string
	algorithmHash string
	dataHash      []string
	startFunc     startFunc
	uploadFunc    uploadFunc
}

type startFunc func(cmd string) (*docker.Container, error)
type uploadFunc func(id string) error

// Create a docker container
func Create() (*Container, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, fmt.Errorf("Error getting docker client, error: %v", err)
	}

	startFunc := func(cmd string) (*docker.Container, error) {
		return client.CreateContainer(docker.CreateContainerOptions{
			Config: &docker.Config{
				Image:        util.ParseDockerfileTemplate("$(DOCKER_NS)/fabric-ccenv:latest"),
				Cmd:          []string{"/bin/sh", "-c", cmd},
				AttachStdout: true,
				AttachStderr: true,
			},
		})
	}

	return &Container{
		address:   "/tmp/teetask/container/",
		client:    client,
		startFunc: startFunc,
	}, nil
}

// Upload for docker
func (c *Container) Upload(algorithm []byte, dataList [][]byte) error {
	if len(algorithm) == 0 {
		return fmt.Errorf("algorithm bytes is empty")
	}

	if len(dataList) == 0 {
		return fmt.Errorf("dataList bytes is empty")
	}

	payload := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(payload)
	tw := tar.NewWriter(gw)

	// Calculate algorithm hash
	hash := sha256.Sum256(algorithm)
	c.algorithmHash = hex.EncodeToString(hash[:])
	c.cmd = filepath.Join(c.address, "main")
	err := util.WriteBytesToPackage(c.cmd, algorithm, tw)
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
		if err = util.WriteBytesToPackage(arg, data, tw); err != nil {
			return err
		}

		c.cmd += " " + arg
	}

	// Write the tar file out
	if err := tw.Close(); err != nil {
		return fmt.Errorf("Error writing files to upload to Docker instance into a temporary tar blob: %s", err)
	}

	gw.Close()

	c.uploadFunc = func(containerID string) error {
		return c.client.UploadToContainer(containerID, docker.UploadToContainerOptions{
			InputStream:          bytes.NewReader(payload.Bytes()),
			Path:                 "/",
			NoOverwriteDirNonDir: false,
		})
	}
	return nil
}

// Verify for docker
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

// Execute for docker
func (c *Container) Execute() ([]byte, error) {
	cmd := fmt.Sprintf("chmod +x %s/* && %s", c.address, c.cmd)
	container, err := c.startFunc(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create an ephemeral docker container")
	}
	c.id = container.ID

	if err = c.uploadFunc(c.id); err != nil {
		return nil, errors.Wrapf(err, "failed to upload payload to the ephemeral docker container, containerID: %v", c.id)
	}

	stdout := bytes.NewBuffer(nil)
	cw, err := c.client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    c.id,
		OutputStream: stdout,
		ErrorStream:  stdout,
		Logs:         true,
		Stdout:       true,
		Stderr:       true,
		Stream:       true,
	})
	if err != nil {
		return nil, fmt.Errorf("Error attaching to container: %s", err)
	}

	if err := c.client.StartContainer(c.id, nil); err != nil {
		cw.Close()
		return nil, errors.Wrapf(err, "failed to realize the Cmd specified at container creation, cmd: %v", cmd)
	}

	retval, err := c.client.WaitContainer(c.id)
	if err != nil {
		cw.Close()
		return nil, fmt.Errorf("Error waiting for container to complete: %s", err)
	}

	// Wait for stream copying to complete before accessing stdout.
	cw.Close()
	if err := cw.Wait(); err != nil {
		return nil, fmt.Errorf("attach wait failed: %s", err)
	}

	if retval > 0 {
		return nil, fmt.Errorf("Error returned from build: %d \"%s\"", retval, stdout.String())
	}

	return stdout.Bytes(), nil
}

// Destroy for docker
func (c *Container) Destroy() error {
	return c.client.RemoveContainer(docker.RemoveContainerOptions{ID: c.id})
}

// GetPublicKey for docker
func (c *Container) GetPublicKey() crypto.PublicKey {
	return nil
}
