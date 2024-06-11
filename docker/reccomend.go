package docker

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

const (
	numReccomendation = 5
)

type ReccomendationItem struct {
	Id   string
	User string
	Name string
	Size int64
}

type Reccomendation struct {
	Images     []*ReccomendationItem
	Containers []*ReccomendationItem
}

func GetReccomendation(ctx context.Context, cli *client.Client) (*Reccomendation, error) {
	images, err := cli.ImageList(ctx, image.ListOptions{SharedSize: true})
	if err != nil {
		return nil, err
	}
	slog.Info("Images", "len", len(images))

	candidateImages := make([]*image.Summary, 0)
	for _, image := range images {
		if image.Containers > 0 {
			continue
		}
		candidateImages = append(candidateImages, &image)
	}
	sort.Slice(candidateImages, func(i, j int) bool {
		return candidateImages[i].Size-candidateImages[i].SharedSize > candidateImages[j].Size-candidateImages[j].SharedSize
	})
	reccomendedImages := make([]*ReccomendationItem, 0)
	for _, image := range candidateImages[:min(len(candidateImages), numReccomendation)] {
		reccomendedImages = append(reccomendedImages, &ReccomendationItem{
			Id:   image.ID,
			User: image.Labels["maintainer"],
			Name: image.RepoTags[0],
			Size: image.Size - image.SharedSize,
		})
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true, Size: true})
	if err != nil {
		return nil, err
	}
	slog.Info("Containers", "len", len(containers))

	candidateContainers := make([]*types.Container, 0)
	for _, container := range containers {
		if strings.Contains(container.State, "running") ||
			strings.Contains(container.State, "paused") ||
			strings.Contains(container.State, "restarting") {
			continue
		}
		candidateContainers = append(candidateContainers, &container)
	}
	sort.Slice(candidateContainers, func(i, j int) bool {
		return candidateContainers[i].SizeRw > candidateContainers[j].SizeRw
	})
	reccomendedContainers := make([]*ReccomendationItem, 0)
	for _, container := range candidateContainers[:min(len(candidateContainers), numReccomendation)] {
		reccomendedContainers = append(reccomendedContainers, &ReccomendationItem{
			Id:   container.ID,
			User: container.Labels["maintainer"],
			Name: container.Names[0],
			Size: container.SizeRw,
		})
	}
	return &Reccomendation{
		Images:     reccomendedImages,
		Containers: reccomendedContainers,
	}, nil
}
