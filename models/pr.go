package models

import (
	"github.com/google/go-github/v50/github"
)

type PullRequest struct {
	Number     int
	Repo       string
	Owner      string
	HeadBranch string
	BaseBranch string
}

func FromPREvent(event *github.PullRequestEvent) *PullRequest {
	pr := PullRequest{
		Number:     event.GetNumber(),
		Owner:      event.GetRepo().GetOwner().GetLogin(),
		Repo:       event.GetRepo().GetName(),
		HeadBranch: event.GetPullRequest().GetHead().GetRef(),
		BaseBranch: event.GetPullRequest().GetBase().GetRef(),
	}

	return &pr
}
