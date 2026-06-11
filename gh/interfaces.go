package gh

import (
	"context"
	"io"

	"github.com/google/go-github/v50/github"
)

type issuesInterface interface {
	ListComments(
		ctx context.Context,
		owner, repo string,
		number int,
		opts *github.IssueListCommentsOptions,
	) ([]*github.IssueComment, *github.Response, error)
	EditComment(
		ctx context.Context,
		owner, repo string,
		id int64,
		comment *github.IssueComment,
	) (*github.IssueComment, *github.Response, error)
	CreateComment(
		ctx context.Context,
		owner, repo string,
		number int,
		comment *github.IssueComment,
	) (*github.IssueComment, *github.Response, error)
	AddLabelsToIssue(
		ctx context.Context,
		owner, repo string,
		number int,
		labels []string,
	) ([]*github.Label, *github.Response, error)
}

type pullRequestsInterface interface {
	ListFiles(
		ctx context.Context,
		owner, repo string,
		number int,
		opts *github.ListOptions,
	) ([]*github.CommitFile, *github.Response, error)
}

type repositoriesInterface interface {
	DownloadContents(
		ctx context.Context,
		owner, repo, filepath string,
		opts *github.RepositoryContentGetOptions,
	) (io.ReadCloser, *github.Response, error)
}
