package gh

import (
	"context"

	"github.com/vimek-go/pr-analizer-action/models"
)

func (c *Client) AddLabels(ctx context.Context, pr *models.PullRequest, labels []string) error {
	if len(labels) == 0 {
		return nil
	}
	_, _, err := c.issuesService.AddLabelsToIssue(ctx, pr.Owner, pr.Repo, pr.Number, labels)
	return err
}
