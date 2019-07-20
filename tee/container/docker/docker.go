package docker

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"strconv"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gaeanetwork/gaea-core/common"
	"github.com/hyperledger/fabric/core/container/util"
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

// Upload for development
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
	c.algorithmHash = common.BytesToHex(hash[:])
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
		c.dataHash = append(c.dataHash, common.BytesToHex(hash[:]))

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

// Destroy for development
func (c *Container) Destroy() error {
	return c.client.RemoveContainer(docker.RemoveContainerOptions{ID: c.id})
}
