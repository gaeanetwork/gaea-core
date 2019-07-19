package docker

import docker "github.com/fsouza/go-dockerclient"

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

// New a docker container
func New() *Container {
	return &Container{}
}
