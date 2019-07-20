package docker

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
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

// Destroy for development
func (c *Container) Destroy() error {
	return c.client.RemoveContainer(docker.RemoveContainerOptions{ID: c.id})
}
