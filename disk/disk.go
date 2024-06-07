package disk

import (
	"github.com/shirou/gopsutil/v4/disk"
)

type DiskUsage struct {
	Total uint64
	Free  uint64
	Used  uint64
}

func GetDiskUsage(path string) (*DiskUsage, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}
	return &DiskUsage{
		Total: usage.Total,
		Free:  usage.Free,
		Used:  usage.Used,
	}, nil
}
