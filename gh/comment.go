package gh

import (
	"context"
	"log"
	"strings"

	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/google/go-github/v50/github"
	"github.com/pkg/errors"
)

const (
	hiddenMarker = "<!-- github.com/vimek-go/pr-analizer-action-id: 5f9a2b8 -->"
	//nolint:lll
	coffieeSection = `[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?&logo=buy-me-a-coffee&logoColor=black)](https://www.buymeacoffee.com/vimekgo)`
)

func (c *Client) PostOrUpdateComment(ctx context.Context, pr *models.PullRequest, body string) error {
	opts := &github.IssueListCommentsOptions{ListOptions: github.ListOptions{PerPage: 100}}
	var existingComment *github.IssueComment
	for {
		comments, resp, err := c.issuesService.ListComments(ctx, pr.Owner, pr.Repo, pr.Number, opts)
		if err != nil {
			return errors.Wrap(err, "listing comments")
		}

		for _, comment := range comments {
			if strings.Contains(comment.GetBody(), hiddenMarker) {
				existingComment = comment
				break
			}
		}

		if existingComment != nil || resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	body += "\n" + hiddenMarker
	body += "\n" + coffieeSection

	if existingComment != nil {
		log.Printf("Updating existing comment on PR #%d (ID: %d)", pr.Number, existingComment.GetID())
		_, _, err := c.issuesService.EditComment(
			ctx,
			pr.Owner,
			pr.Repo,
			existingComment.GetID(),
			&github.IssueComment{Body: new(body)},
		)
		if err != nil {
			return errors.Wrap(err, "editing comment")
		}
	} else {
		log.Printf("Creating new comment on PR #%d", pr.Number)
		_, _, err := c.issuesService.CreateComment(
			ctx,
			pr.Owner,
			pr.Repo,
			pr.Number,
			&github.IssueComment{Body: new(body)},
		)
		if err != nil {
			return errors.Wrap(err, "creating comment")
		}
	}

	return nil
}
