package engine

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/selyusize/my-home/pkg/docker"

	moby "github.com/moby/moby/client"
)

var _ docker.Client = (*Client)(nil)

// Client is a Docker Engine API wrapper around the official Moby client.
type Client struct {
	mu     sync.Mutex
	cfg    docker.Config
	client *moby.Client
}

// Register registers the Docker Engine constructor on a docker.Factory.
func Register() docker.Option {
	return docker.WithProvider(docker.ProviderEngine, New)
}

// New creates a Docker Engine client. The daemon connection is established lazily.
func New(cfg docker.Config) (docker.Client, error) {
	return &Client{cfg: cfg}, nil
}

// Provider returns the Docker Engine provider identifier.
func (c *Client) Provider() docker.Provider {
	return docker.ProviderEngine
}

// Close closes the underlying engine connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.client == nil {
		return nil
	}
	err := c.client.Close()
	c.client = nil
	return err
}

func (c *Client) api() (*moby.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.client != nil {
		return c.client, nil
	}

	opts := []moby.Opt{moby.FromEnv}
	if c.cfg.Host != "" {
		opts = append(opts, moby.WithHost(c.cfg.Host))
	}
	if c.cfg.APIVersion != "" {
		opts = append(opts, moby.WithAPIVersion(c.cfg.APIVersion))
	}
	if c.cfg.CertPath != "" {
		opts = append(opts, moby.WithTLSClientConfig(
			filepath.Join(c.cfg.CertPath, "ca.pem"),
			filepath.Join(c.cfg.CertPath, "cert.pem"),
			filepath.Join(c.cfg.CertPath, "key.pem"),
		))
	}

	client, err := moby.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("docker: create engine client: %w", err)
	}
	c.client = client
	return c.client, nil
}
