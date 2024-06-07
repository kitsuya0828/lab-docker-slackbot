package disk

import (
	"testing"
)

func TestDisk(t *testing.T) {
	usage, err := GetDiskUsage("/")
	if err != nil {
		t.Errorf("Failed to get disk usage: %s", err)
	}
	t.Logf("Disk usage: %v", usage)
}
