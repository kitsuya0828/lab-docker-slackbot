package docker

import (
	"context"
	"testing"
)

func TestGetReccomendation(t *testing.T) {
	ctx := context.Background()
	cli, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	reccomendation, err := GetReccomendation(ctx, cli)
	if err != nil {
		t.Fatal(err)
	}
	for _, image := range reccomendation.Images {
		t.Logf("Image: %v", image)
	}
	for _, container := range reccomendation.Containers {
		t.Logf("Container: %v", container)
	}
}
