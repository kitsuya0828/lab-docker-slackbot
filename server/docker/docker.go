package docker

import (
	"github.com/docker/docker/client"
)

func NewClient() (*client.Client, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return client, nil
}
