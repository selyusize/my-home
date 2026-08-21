package engine

import (
	"context"
	"fmt"

	"github.com/selyusize/my-home/pkg/docker"

	moby "github.com/moby/moby/client"
)

// Ping checks that the engine is reachable.
func (c *Client) Ping(ctx context.Context) (*docker.Ping, error) {
	client, err := c.api()
	if err != nil {
		return nil, err
	}

	result, err := client.Ping(ctx, moby.PingOptions{NegotiateAPIVersion: true})
	if err != nil {
		return nil, fmt.Errorf("docker: ping: %w", err)
	}
	return &docker.Ping{
		APIVersion:     result.APIVersion,
		OSType:         result.OSType,
		IsExperimental: result.Experimental,
	}, nil
}

// Info returns a summary of the docker daemon.
func (c *Client) Info(ctx context.Context) (*docker.Info, error) {
	client, err := c.api()
	if err != nil {
		return nil, err
	}

	result, err := client.Info(ctx, moby.InfoOptions{})
	if err != nil {
		return nil, fmt.Errorf("docker: info: %w", err)
	}

	info := result.Info
	return &docker.Info{
		ID:                info.ID,
		Name:              info.Name,
		ServerVersion:     info.ServerVersion,
		OperatingSystem:   info.OperatingSystem,
		OSVersion:         info.OSVersion,
		OSType:            info.OSType,
		Architecture:      info.Architecture,
		NCPU:              info.NCPU,
		MemTotal:          info.MemTotal,
		Driver:            info.Driver,
		RootDir:           info.DockerRootDir,
		Containers:        info.Containers,
		ContainersRunning: info.ContainersRunning,
		ContainersPaused:  info.ContainersPaused,
		ContainersStopped: info.ContainersStopped,
		Images:            info.Images,
	}, nil
}
