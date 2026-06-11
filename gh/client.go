package gh

import (
	"context"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type Client struct {
	issuesService       issuesInterface
	pullRequestsService pullRequestsInterface
	repositoriesService repositoriesInterface
}

func NewClient(ctx context.Context, token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	gc := github.NewClient(tc)

	return &Client{
		issuesService:       gc.Issues,
		pullRequestsService: gc.PullRequests,
		repositoriesService: gc.Repositories,
	}
}
