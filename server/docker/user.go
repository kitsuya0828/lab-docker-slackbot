package docker

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/docker/docker/client"
)

// getUsernames returns a list of usernames on the system
func getUsernames() ([]string, error) {
	usernames := make([]string, 0)
	switch runtime.GOOS {
	case "linux":
		out, err := exec.Command("getent", "passwd").Output()
		if err != nil {
			return nil, fmt.Errorf("failed to list users: %v", err)
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			parts := strings.Split(line, ":")
			if len(parts) == 0 {
				continue
			}
			usernames = append(usernames, parts[0])
		}
		return usernames, nil
	case "darwin":
		out, err := exec.Command("dscl", ".", "-list", "/Users").Output()
		if err != nil {
			return nil, fmt.Errorf("failed to list users: %v", err)
		}
		parts := strings.Split(string(out), "\n")
		for _, username := range parts {
			if strings.TrimSpace(username) == "" || strings.HasPrefix(username, "_") {
				continue
			}
			usernames = append(usernames, username)
		}
		return usernames, nil
	default:
		return nil, fmt.Errorf("Unsupported OS: %v", runtime.GOOS)
	}
}

func getUsernameFromImageHistory(cli *client.Client, imageID string) (string, error) {
	usernames, err := getUsernames()
	if err != nil {
		return "", fmt.Errorf("Error getting usernames: %v", err)
	}
	history, err := cli.ImageHistory(context.Background(), imageID)
	if err != nil {
		return "", fmt.Errorf("Error getting history for image %s: %v", imageID, err)
	}

	r := regexp.MustCompile(`DOCKER_USER=(\w+)`)
	itemCreatedBys := make([]string, 0)
	for _, item := range history {
		if matches := r.FindStringSubmatch(item.CreatedBy); matches != nil {
			return matches[1], nil
		}
		itemCreatedBys = append(itemCreatedBys, item.CreatedBy)
	}

	itemCreatedByStr := strings.Join(itemCreatedBys, " ")
	maxCount := 0
	mostFrequentUsername := ""
	for _, username := range usernames {
		count := strings.Count(itemCreatedByStr, username)
		if count > 0 && count >= maxCount {
			maxCount = count
			mostFrequentUsername = username
		}
	}

	if mostFrequentUsername != "" {
		return mostFrequentUsername, nil
	}

	return "", fmt.Errorf("No user found in history for image %s", imageID)
}
