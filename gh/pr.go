package gh

import (
	"context"
	"io"

	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/google/go-github/v50/github"
)

func (c *Client) GetChangedFiles(ctx context.Context, pr *models.PullRequest) ([]models.ChangedFile, error) {
	var allFiles []models.ChangedFile
	opts := &github.ListOptions{PerPage: 100}

	for {
		files, resp, err := c.pullRequestsService.ListFiles(ctx, pr.Owner, pr.Repo, pr.Number, opts)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			if f.GetFilename() == "" {
				continue
			}
			allFiles = append(allFiles, models.ChangedFile{
				Filename: f.GetFilename(),
				Status:   f.GetStatus(),
				Patch:    f.GetPatch(),
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allFiles, nil
}

func (c *Client) GetFileContent(ctx context.Context, owner, repo, path, ref string) (io.ReadCloser, error) {
	opts := &github.RepositoryContentGetOptions{Ref: ref}
	content, _, err := c.repositoriesService.DownloadContents(ctx, owner, repo, path, opts)
	if err != nil {
		return nil, err
	}
	return content, nil
}
