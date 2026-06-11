package event

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/google/go-github/v50/github"
	"github.com/pkg/errors"
)

func readEventFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read event file")
	}
	return data, nil
}

func parsePREvent(eventData []byte) (*github.PullRequestEvent, error) {
	var prEvent github.PullRequestEvent
	if err := json.Unmarshal(eventData, &prEvent); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal pull request event")
	}
	return &prEvent, nil
}

func GetPRDetails() (*models.PullRequest, error) {
	eventFilePath := os.Getenv("GITHUB_EVENT_PATH")
	if eventFilePath == "" {
		return nil, errors.New("GITHUB_EVENT_PATH is not set")
	}

	eventData, err := readEventFile(eventFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading event file")
	}

	prEvent, err := parsePREvent(eventData)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing event file")
	}

	return models.FromPREvent(prEvent), nil
}
