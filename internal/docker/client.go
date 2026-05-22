package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ContainerInfo holds relevant runtime info for a container.
type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Labels  map[string]string
	Env     []string
	Ports   []string
	Running bool
}

// Client wraps the Docker SDK client.
type Client struct {
	docker *client.Client
}

// NewClient creates a new Docker client using the host environment.
func NewClient() (*Client, error) {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("creating docker client: %w", err)
	}
	return &Client{docker: c}, nil
}

// Close releases the underlying Docker client resources.
func (c *Client) Close() error {
	return c.docker.Close()
}

// ListContainers returns running containers, optionally filtered by compose project label.
func (c *Client) ListContainers(ctx context.Context, projectName string) ([]ContainerInfo, error) {
	f := filters.NewArgs()
	if projectName != "" {
		f.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
	}

	containers, err := c.docker.ContainerList(ctx, types.ContainerListOptions{
		All:     false,
		Filters: f,
	})
	if err != nil {
		return nil, fmt.Errorf("listing containers: %w", err)
	}

	result := make([]ContainerInfo, 0, len(containers))
	for _, ct := range containers {
		name := ""
		if len(ct.Names) > 0 {
			name = ct.Names[0]
		}
		ports := make([]string, 0, len(ct.Ports))
		for _, p := range ct.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
		}
		result = append(result, ContainerInfo{
			ID:      ct.ID[:12],
			Name:    name,
			Image:   ct.Image,
			Labels:  ct.Labels,
			Ports:   ports,
			Running: ct.State == "running",
		})
	}
	return result, nil
}
