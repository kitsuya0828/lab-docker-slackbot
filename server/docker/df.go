package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type Stat struct {
	Active      uint64
	Size        float64
	Reclaimable float64
	TotalCount  uint64
}

type DiskUsage struct {
	Images       Stat
	Containers   Stat
	LocalVolumes Stat
	BuildCache   Stat
}

func GetDiskUsage(ctx context.Context, cli *client.Client) (*DiskUsage, error) {
	d, err := cli.DiskUsage(ctx, types.DiskUsageOptions{})
	if err != nil {
		return nil, err
	}

	var bsz int64
	for _, bc := range d.BuildCache {
		if !bc.Shared {
			bsz += bc.Size
		}
	}

	du := DiskUsage{
		Images:       getImagesStat(d.Images, d.LayersSize),
		Containers:   getContainersStat(d.Containers),
		LocalVolumes: getLocalVolumesStat(d.Volumes),
		BuildCache:   getBuildCacheStat(d.BuildCache, bsz),
	}

	return &du, nil
}

func getImagesStat(images []*image.Summary, layersSize int64) Stat {
	used := uint64(0)
	usedSize := int64(0)
	for _, image := range images {
		if image.Containers > 0 {
			used++
		}
		if image.Containers != 0 {
			if image.Size == -1 || image.SharedSize == -1 {
				continue
			}
			usedSize += image.Size - image.SharedSize
		}
	}
	s := Stat{
		Active:      used,
		Size:        float64(layersSize),
		Reclaimable: float64(layersSize - usedSize),
		TotalCount:  uint64(len(images)),
	}
	return s
}

func getContainersStat(containers []*types.Container) Stat {
	used := uint64(0)
	reclaimable := int64(0)
	totalSize := int64(0)
	for _, container := range containers {
		if strings.Contains(container.State, "running") ||
			strings.Contains(container.State, "paused") ||
			strings.Contains(container.State, "restarting") {
			used++
		} else {
			reclaimable += container.SizeRw
		}
		totalSize += container.SizeRw
	}
	s := Stat{
		Active:      used,
		Size:        float64(totalSize),
		Reclaimable: float64(reclaimable),
		TotalCount:  uint64(len(containers)),
	}
	return s
}

func getLocalVolumesStat(volumes []*volume.Volume) Stat {
	used := uint64(0)
	reclaimable := int64(0)
	totalSize := int64(0)
	for _, volume := range volumes {
		if volume.UsageData.RefCount > 0 {
			used++
		}
		if volume.UsageData.Size != -1 {
			if volume.UsageData.RefCount == 0 {
				reclaimable += volume.UsageData.Size
			}
			totalSize += volume.UsageData.Size
		}
	}
	s := Stat{
		Active:      used,
		Size:        float64(totalSize),
		Reclaimable: float64(reclaimable),
		TotalCount:  uint64(len(volumes)),
	}
	return s
}

func getBuildCacheStat(buildCache []*types.BuildCache, bsz int64) Stat {
	numActive := uint64(0)
	inUseBytes := int64(0)
	usedSize := int64(0)
	for _, bc := range buildCache {
		if bc.InUse {
			numActive++
			if !bc.Shared {
				inUseBytes += bc.Size
			}
		}
	}
	s := Stat{
		Active:      numActive,
		Size:        float64(bsz),
		Reclaimable: float64(bsz - usedSize),
		TotalCount:  uint64(len(buildCache)),
	}
	return s
}
