package docker

import (
	"github.com/selyusize/my-home/internal/config"
	pkgdocker "github.com/selyusize/my-home/pkg/docker"
)

type DockerService struct {
	pkgdocker.Client
}

func NewDockerService() (*DockerService, error) {
	client, err := factory.New(pkgdocker.ProviderEngine, config.Docker())
	if err != nil {
		return nil, err
	}
	return &DockerService{Client: client}, nil
}
