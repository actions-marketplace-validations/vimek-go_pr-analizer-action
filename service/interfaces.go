package service

import (
	"context"
	"io"

	"github.com/vimek-go/pr-analizer-action/models"
)

type Client interface {
	GetChangedFiles(ctx context.Context, pr *models.PullRequest) ([]models.ChangedFile, error)
	GetFileContent(ctx context.Context, owner, repo, path, ref string) (io.ReadCloser, error)
	PostOrUpdateComment(ctx context.Context, pr *models.PullRequest, body string) error
	AddLabels(ctx context.Context, pr *models.PullRequest, labels []string) error
}
