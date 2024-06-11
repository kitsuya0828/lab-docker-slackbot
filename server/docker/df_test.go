package docker

import (
	"context"
	"testing"
)

func TestGetDiskUsage(t *testing.T) {
	ctx := context.Background()
	cli, err := NewClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}
	du, err := GetDiskUsage(ctx, cli)
	if err != nil {
		t.Errorf("Error getting disk usage: %v", err)
	}
	t.Logf("Disk usage: %v", du)
}
