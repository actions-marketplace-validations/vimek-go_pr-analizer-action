package gh_test

import (
	"context"
	"testing"

	"github.com/vimek-go/pr-analizer-action/gh"
	"github.com/vimek-go/pr-analizer-action/gh/mocks"
	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_PostOrUpdateComment(t *testing.T) {
	t.Parallel()
	const (
		prNumber   = 12
		owner      = "test-owner"
		repo       = "test-repo"
		existingID = int64(42)
	)
	opts := &github.IssueListCommentsOptions{ListOptions: github.ListOptions{PerPage: 100}}

	pr := &models.PullRequest{Owner: owner, Repo: repo, Number: prNumber}
	markedBody := "old body\n" + gh.HiddenMarker

	tests := []struct {
		name        string
		setup       func(t *testing.T) *mocks.MockIssuesInterface
		expectedErr error
	}{
		{
			name: "no existing comment - creates new",
			setup: func(t *testing.T) *mocks.MockIssuesInterface {
				t.Helper()
				m := mocks.NewMockIssuesInterface(t)
				m.On("ListComments", context.Background(), owner, repo, prNumber, opts).
					Return([]*github.IssueComment{}, &github.Response{NextPage: 0}, nil).Once()
				m.On("CreateComment", context.Background(), owner, repo, prNumber, mock.Anything).
					Return(&github.IssueComment{}, &github.Response{}, nil).Once()
				return m
			},
		},
		{
			name: "existing comment on first page - updates it",
			setup: func(t *testing.T) *mocks.MockIssuesInterface {
				t.Helper()
				m := mocks.NewMockIssuesInterface(t)
				m.On("ListComments", context.Background(), owner, repo, prNumber, opts).
					Return([]*github.IssueComment{
						{ID: new(existingID), Body: new(markedBody)},
					}, &github.Response{NextPage: 0}, nil).Once()
				m.On("EditComment", context.Background(), owner, repo, existingID, mock.Anything).
					Return(&github.IssueComment{}, &github.Response{}, nil).Once()
				return m
			},
		},
		{
			name: "existing comment found on second page - paginates then updates",
			setup: func(t *testing.T) *mocks.MockIssuesInterface {
				t.Helper()
				m := mocks.NewMockIssuesInterface(t)
				m.On("ListComments", context.Background(), owner, repo, prNumber, opts).
					Return([]*github.IssueComment{
						{ID: new(int64(99)), Body: new("unrelated comment")},
					}, &github.Response{NextPage: 2}, nil).Once()
				m.On(
					"ListComments",
					context.Background(),
					owner,
					repo,
					prNumber,
					&github.IssueListCommentsOptions{ListOptions: github.ListOptions{PerPage: 100, Page: 2}},
				).
					Return(
						[]*github.IssueComment{
							{ID: new(existingID), Body: new(markedBody)},
						},
						&github.Response{NextPage: 0},
						nil,
					).
					Once()
				m.On("EditComment", context.Background(), owner, repo, existingID, mock.Anything).
					Return(&github.IssueComment{}, &github.Response{}, nil).Once()
				return m
			},
		},
		{
			name: "no existing comment across multiple pages - creates new",
			setup: func(t *testing.T) *mocks.MockIssuesInterface {
				t.Helper()
				m := mocks.NewMockIssuesInterface(t)
				m.On("ListComments", context.Background(), owner, repo, prNumber, opts).
					Return([]*github.IssueComment{
						{ID: new(int64(10)), Body: new("page 1 comment")},
					}, &github.Response{NextPage: 2}, nil).Once()
				m.On(
					"ListComments",
					context.Background(),
					owner,
					repo,
					prNumber,
					&github.IssueListCommentsOptions{ListOptions: github.ListOptions{PerPage: 100, Page: 2}},
				).
					Return(
						[]*github.IssueComment{
							{ID: new(int64(11)), Body: new("page 2 comment")},
						},
						&github.Response{NextPage: 0},
						nil,
					).
					Once()
				m.On("CreateComment", context.Background(), owner, repo, prNumber, mock.Anything).
					Return(&github.IssueComment{}, &github.Response{}, nil).Once()
				return m
			},
		},
		{
			name: "listComments error - propagated",
			setup: func(t *testing.T) *mocks.MockIssuesInterface {
				t.Helper()
				m := mocks.NewMockIssuesInterface(t)
				m.On("ListComments", context.Background(), owner, repo, prNumber, opts).
					Return(nil, nil, assert.AnError).Once()
				return m
			},
			expectedErr: assert.AnError,
		},
		{
			name: "createComment error - propagated",
			setup: func(t *testing.T) *mocks.MockIssuesInterface {
				t.Helper()
				m := mocks.NewMockIssuesInterface(t)
				m.On("ListComments", context.Background(), owner, repo, prNumber, opts).
					Return([]*github.IssueComment{}, &github.Response{NextPage: 0}, nil).Once()
				m.On("CreateComment", context.Background(), owner, repo, prNumber, mock.Anything).
					Return(nil, nil, assert.AnError).Once()
				return m
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := gh.NewTestClient(tc.setup(t), nil, nil)
			err := client.PostOrUpdateComment(context.Background(), pr, "body")
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
