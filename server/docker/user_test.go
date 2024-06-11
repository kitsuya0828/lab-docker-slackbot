package docker

import (
	"testing"
)

func TestGetUsernames(t *testing.T) {
	usernames, err := getUsernames()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if len(usernames) == 0 {
		t.Errorf("Expected usernames to be non-empty")
	}
	t.Logf("Usernames: %v", usernames)
}
