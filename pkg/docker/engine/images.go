package engine

import (
	"context"
	"fmt"

	"github.com/selyusize/my-home/pkg/docker"

	"github.com/moby/moby/api/types/image"
	moby "github.com/moby/moby/client"
)

// ListImages returns images on the engine.
func (c *Client) ListImages(ctx context.Context, all bool) ([]docker.Image, error) {
	client, err := c.api()
	if err != nil {
		return nil, err
	}

	result, err := client.ImageList(ctx, moby.ImageListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("docker: list images: %w", err)
	}

	out := make([]docker.Image, 0, len(result.Items))
	for _, item := range result.Items {
		out = append(out, mapImage(item))
	}
	return out, nil
}

// InspectImage returns detailed information about an image.
func (c *Client) InspectImage(ctx context.Context, id string) (*docker.ImageDetails, error) {
	client, err := c.api()
	if err != nil {
		return nil, err
	}
	if id == "" {
		return nil, fmt.Errorf("docker: image id is required")
	}

	result, err := client.ImageInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("docker: inspect image %q: %w", id, err)
	}

	return &docker.ImageDetails{
		Image: docker.Image{
			ID:          result.ID,
			RepoTags:    result.RepoTags,
			RepoDigests: result.RepoDigests,
			Size:        result.Size,
			Created:     parseTime(result.Created),
		},
		Author:       result.Author,
		Comment:      result.Comment,
		Architecture: result.Architecture,
		OS:           result.Os,
	}, nil
}

// RemoveImage deletes an image from the engine.
func (c *Client) RemoveImage(ctx context.Context, id string, force bool) error {
	client, err := c.api()
	if err != nil {
		return err
	}
	if id == "" {
		return fmt.Errorf("docker: image id is required")
	}
	if _, err := client.ImageRemove(ctx, id, moby.ImageRemoveOptions{
		Force:         force,
		PruneChildren: true,
	}); err != nil {
		return fmt.Errorf("docker: remove image %q: %w", id, err)
	}
	return nil
}

func mapImage(item image.Summary) docker.Image {
	return docker.Image{
		ID:          item.ID,
		RepoTags:    item.RepoTags,
		RepoDigests: item.RepoDigests,
		Labels:      item.Labels,
		Size:        item.Size,
		Containers:  item.Containers,
		Created:     unixOrZero(item.Created),
	}
}
