package docker

import "context"

// Client is the common API implemented by every container engine provider.
type Client interface {
	Provider() Provider
	Close() error

	Ping(ctx context.Context) (*Ping, error)
	Info(ctx context.Context) (*Info, error)

	ListContainers(ctx context.Context, opts ContainerListOptions) ([]Container, error)
	InspectContainer(ctx context.Context, id string) (*ContainerDetails, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string, opts StopOptions) error
	RestartContainer(ctx context.Context, id string, opts StopOptions) error
	RemoveContainer(ctx context.Context, id string, opts ContainerRemoveOptions) error

	ListImages(ctx context.Context, all bool) ([]Image, error)
	InspectImage(ctx context.Context, id string) (*ImageDetails, error)
	RemoveImage(ctx context.Context, id string, force bool) error

	ListVolumes(ctx context.Context) ([]Volume, error)
	ListNetworks(ctx context.Context) ([]Network, error)
}

// Constructor creates a Client for a registered provider.
type Constructor func(cfg Config) (Client, error)
