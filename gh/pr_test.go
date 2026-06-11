package gh_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/vimek-go/pr-analizer-action/gh"
	"github.com/vimek-go/pr-analizer-action/gh/mocks"
	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetChangedFiles(t *testing.T) {
	t.Parallel()

	const (
		owner    = "test-owner"
		repo     = "test-repo"
		prNumber = 123
	)

	pr := &models.PullRequest{Owner: owner, Repo: repo, Number: prNumber}
	opts := &github.ListOptions{PerPage: 100}

	tests := []struct {
		name        string
		setup       func(t *testing.T) *mocks.MockPullRequestsInterface
		expected    []models.ChangedFile
		expectedErr error
	}{
		{
			name: "single page - returns all files",
			setup: func(t *testing.T) *mocks.MockPullRequestsInterface {
				t.Helper()
				m := mocks.NewMockPullRequestsInterface(t)
				m.On("ListFiles", context.Background(), owner, repo, prNumber, opts).
					Return([]*github.CommitFile{
						{
							Filename: new("main.go"),
							Status:   new("modified"),
							Patch:    new("@@ -1 +1 @@"),
						},
						{
							Filename: new("README.md"),
							Status:   new("added"),
							Patch:    new(""),
						},
					}, &github.Response{NextPage: 0}, nil).Once()
				return m
			},
			expected: []models.ChangedFile{
				{Filename: "main.go", Status: "modified", Patch: "@@ -1 +1 @@"},
				{Filename: "README.md", Status: "added", Patch: ""},
			},
		},
		{
			name: "pagination - accumulates files from all pages",
			setup: func(t *testing.T) *mocks.MockPullRequestsInterface {
				t.Helper()
				m := mocks.NewMockPullRequestsInterface(t)
				m.On("ListFiles", context.Background(), owner, repo, prNumber, opts).
					Return([]*github.CommitFile{
						{Filename: new("a.go"), Status: new("added")},
					}, &github.Response{NextPage: 2}, nil).Once()
				m.On("ListFiles", context.Background(), owner, repo, prNumber, &github.ListOptions{PerPage: 100, Page: 2}).
					Return([]*github.CommitFile{
						{Filename: new("b.go"), Status: new("modified")},
					}, &github.Response{NextPage: 0}, nil).
					Once()
				return m
			},
			expected: []models.ChangedFile{
				{Filename: "a.go", Status: "added"},
				{Filename: "b.go", Status: "modified"},
			},
		},
		{
			name: "files with empty filename are skipped",
			setup: func(t *testing.T) *mocks.MockPullRequestsInterface {
				t.Helper()
				m := mocks.NewMockPullRequestsInterface(t)
				m.On("ListFiles", context.Background(), owner, repo, prNumber, opts).
					Return([]*github.CommitFile{
						{Filename: new(""), Status: new("added")},
						{Filename: new("valid.go"), Status: new("added")},
					}, &github.Response{NextPage: 0}, nil).Once()
				return m
			},
			expected: []models.ChangedFile{
				{Filename: "valid.go", Status: "added"},
			},
		},
		{
			name: "empty result - returns nil slice",
			setup: func(t *testing.T) *mocks.MockPullRequestsInterface {
				t.Helper()
				m := mocks.NewMockPullRequestsInterface(t)
				m.On("ListFiles", context.Background(), owner, repo, prNumber, opts).
					Return([]*github.CommitFile{}, &github.Response{NextPage: 0}, nil).Once()
				return m
			},
			expected: nil,
		},
		{
			name: "github SDK error",
			setup: func(t *testing.T) *mocks.MockPullRequestsInterface {
				t.Helper()
				m := mocks.NewMockPullRequestsInterface(t)
				m.On("ListFiles", context.Background(), owner, repo, prNumber, opts).
					Return(nil, nil, assert.AnError).Once()
				return m
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := gh.NewTestClient(nil, tc.setup(t), nil)
			actual, err := client.GetChangedFiles(context.Background(), pr)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestClient_GetFileContent(t *testing.T) {
	t.Parallel()

	const (
		owner    = "test-owner"
		repo     = "test-repo"
		prNumber = 123
		path     = "main.go"
		ref      = "main"
	)
	opts := &github.RepositoryContentGetOptions{Ref: ref}

	tests := []struct {
		name        string
		setup       func(t *testing.T) *mocks.MockRepositoriesInterface
		expectedErr error
	}{
		{
			name: "returns reader on success",
			setup: func(t *testing.T) *mocks.MockRepositoriesInterface {
				t.Helper()
				m := mocks.NewMockRepositoriesInterface(t)
				m.On("DownloadContents", context.Background(), owner, repo, path, opts).
					Return(io.NopCloser(strings.NewReader("content")), &github.Response{}, nil).Once()
				return m
			},
		},
		{
			name: "SDK error",
			setup: func(t *testing.T) *mocks.MockRepositoriesInterface {
				t.Helper()
				m := mocks.NewMockRepositoriesInterface(t)
				m.On("DownloadContents", context.Background(), owner, repo, path, opts).
					Return(nil, nil, assert.AnError).Once()
				return m
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := gh.NewTestClient(nil, nil, tc.setup(t))
			actual, err := client.GetFileContent(context.Background(), owner, repo, path, ref)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, actual)
		})
	}
}
