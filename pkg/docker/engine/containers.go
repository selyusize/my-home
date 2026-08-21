package engine

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/selyusize/my-home/pkg/docker"

	"github.com/moby/moby/api/types/container"
	moby "github.com/moby/moby/client"
)

// ListContainers returns containers on the engine.
func (c *Client) ListContainers(ctx context.Context, opts docker.ContainerListOptions) ([]docker.Container, error) {
	client, err := c.api()
	if err != nil {
		return nil, err
	}

	result, err := client.ContainerList(ctx, moby.ContainerListOptions{
		All:   opts.All,
		Limit: opts.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("docker: list containers: %w", err)
	}

	out := make([]docker.Container, 0, len(result.Items))
	for _, item := range result.Items {
		out = append(out, mapContainer(item))
	}
	return out, nil
}

// InspectContainer returns detailed information about a container.
func (c *Client) InspectContainer(ctx context.Context, id string) (*docker.ContainerDetails, error) {
	client, err := c.api()
	if err != nil {
		return nil, err
	}
	if id == "" {
		return nil, fmt.Errorf("docker: container id is required")
	}

	result, err := client.ContainerInspect(ctx, id, moby.ContainerInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("docker: inspect container %q: %w", id, err)
	}
	return mapContainerDetails(result.Container), nil
}

// StartContainer starts a container.
func (c *Client) StartContainer(ctx context.Context, id string) error {
	client, err := c.api()
	if err != nil {
		return err
	}
	if id == "" {
		return fmt.Errorf("docker: container id is required")
	}
	if _, err := client.ContainerStart(ctx, id, moby.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("docker: start container %q: %w", id, err)
	}
	return nil
}

// StopContainer stops a container.
func (c *Client) StopContainer(ctx context.Context, id string, opts docker.StopOptions) error {
	client, err := c.api()
	if err != nil {
		return err
	}
	if id == "" {
		return fmt.Errorf("docker: container id is required")
	}
	if _, err := client.ContainerStop(ctx, id, moby.ContainerStopOptions{Timeout: opts.Timeout}); err != nil {
		return fmt.Errorf("docker: stop container %q: %w", id, err)
	}
	return nil
}

// RestartContainer restarts a container.
func (c *Client) RestartContainer(ctx context.Context, id string, opts docker.StopOptions) error {
	client, err := c.api()
	if err != nil {
		return err
	}
	if id == "" {
		return fmt.Errorf("docker: container id is required")
	}
	if _, err := client.ContainerRestart(ctx, id, moby.ContainerRestartOptions{Timeout: opts.Timeout}); err != nil {
		return fmt.Errorf("docker: restart container %q: %w", id, err)
	}
	return nil
}

// RemoveContainer removes a container.
func (c *Client) RemoveContainer(ctx context.Context, id string, opts docker.ContainerRemoveOptions) error {
	client, err := c.api()
	if err != nil {
		return err
	}
	if id == "" {
		return fmt.Errorf("docker: container id is required")
	}
	if _, err := client.ContainerRemove(ctx, id, moby.ContainerRemoveOptions{
		Force:         opts.Force,
		RemoveVolumes: opts.RemoveVolumes,
	}); err != nil {
		return fmt.Errorf("docker: remove container %q: %w", id, err)
	}
	return nil
}

func mapContainer(item container.Summary) docker.Container {
	names := make([]string, 0, len(item.Names))
	for _, name := range item.Names {
		names = append(names, strings.TrimPrefix(name, "/"))
	}

	ports := make([]docker.Port, 0, len(item.Ports))
	for _, port := range item.Ports {
		ports = append(ports, docker.Port{
			IP:          addrString(port.IP),
			PrivatePort: port.PrivatePort,
			PublicPort:  port.PublicPort,
			Type:        port.Type,
		})
	}

	return docker.Container{
		ID:      item.ID,
		Name:    firstName(names),
		Names:   names,
		Image:   item.Image,
		ImageID: item.ImageID,
		Command: item.Command,
		State:   string(item.State),
		Status:  item.Status,
		Ports:   ports,
		Labels:  item.Labels,
		Created: unixOrZero(item.Created),
	}
}

func mapContainerDetails(item container.InspectResponse) *docker.ContainerDetails {
	names := []string{strings.TrimPrefix(item.Name, "/")}
	details := &docker.ContainerDetails{
		Container: docker.Container{
			ID:      item.ID,
			Name:    firstName(names),
			Names:   names,
			Image:   item.Image,
			Command: strings.TrimSpace(strings.Join(append([]string{item.Path}, item.Args...), " ")),
			Created: parseTime(item.Created),
		},
		Path:         item.Path,
		Args:         item.Args,
		Platform:     item.Platform,
		RestartCount: item.RestartCount,
	}
	if item.State != nil {
		details.State = string(item.State.Status)
		details.PID = item.State.Pid
		details.ExitCode = item.State.ExitCode
		details.StartedAt = parseTime(item.State.StartedAt)
		details.FinishedAt = parseTime(item.State.FinishedAt)
		if item.State.Status != "" {
			details.Status = string(item.State.Status)
		}
	}
	return details
}

func firstName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func unixOrZero(sec int64) time.Time {
	if sec <= 0 {
		return time.Time{}
	}
	return time.Unix(sec, 0).UTC()
}

func parseTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "0001-01-01") {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		if t, err := time.Parse(layout, value); err == nil {
			return t
		}
	}
	return time.Time{}
}
