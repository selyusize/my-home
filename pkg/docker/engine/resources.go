package engine

import (
	"context"
	"fmt"

	"github.com/selyusize/my-home/pkg/docker"

	moby "github.com/moby/moby/client"
)

// ListVolumes returns named volumes on the engine.
func (c *Client) ListVolumes(ctx context.Context) ([]docker.Volume, error) {
	client, err := c.api()
	if err != nil {
		return nil, err
	}

	result, err := client.VolumeList(ctx, moby.VolumeListOptions{})
	if err != nil {
		return nil, fmt.Errorf("docker: list volumes: %w", err)
	}

	out := make([]docker.Volume, 0, len(result.Items))
	for _, item := range result.Items {
		out = append(out, docker.Volume{
			Name:       item.Name,
			Driver:     item.Driver,
			Mountpoint: item.Mountpoint,
			Scope:      item.Scope,
			Labels:     item.Labels,
			CreatedAt:  parseTime(item.CreatedAt),
		})
	}
	return out, nil
}

// ListNetworks returns networks on the engine.
func (c *Client) ListNetworks(ctx context.Context) ([]docker.Network, error) {
	client, err := c.api()
	if err != nil {
		return nil, err
	}

	result, err := client.NetworkList(ctx, moby.NetworkListOptions{})
	if err != nil {
		return nil, fmt.Errorf("docker: list networks: %w", err)
	}

	out := make([]docker.Network, 0, len(result.Items))
	for _, item := range result.Items {
		out = append(out, docker.Network{
			ID:         item.ID,
			Name:       item.Name,
			Driver:     item.Driver,
			Scope:      item.Scope,
			IsInternal: item.Internal,
			Created:    item.Created,
		})
	}
	return out, nil
}
