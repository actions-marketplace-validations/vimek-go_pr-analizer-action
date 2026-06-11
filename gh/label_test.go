package gh_test

import (
	"context"
	"testing"

	"github.com/vimek-go/pr-analizer-action/gh"
	"github.com/vimek-go/pr-analizer-action/gh/mocks"
	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
)

func TestClient_AddLabels(t *testing.T) {
	t.Parallel()
	const prNumber = 123
	pr := &models.PullRequest{Owner: "owner", Repo: "repo", Number: prNumber}

	tests := []struct {
		name        string
		labels      []string
		setup       func(t *testing.T) *mocks.MockIssuesInterface
		expectedErr error
	}{
		{
			name:   "nil slice - returns nil without SDK call",
			labels: nil,
			setup: func(t *testing.T) *mocks.MockIssuesInterface {
				return mocks.NewMockIssuesInterface(t)
			},
		},
		{
			name:   "labels added successfully",
			labels: []string{"small", "medium"},
			setup: func(t *testing.T) *mocks.MockIssuesInterface {
				m := mocks.NewMockIssuesInterface(t)
				m.On("AddLabelsToIssue", context.Background(), "owner", "repo", prNumber, []string{"small", "medium"}).
					Return([]*github.Label{}, &github.Response{}, nil).Once()
				return m
			},
		},
		{
			name:   "github error",
			labels: []string{"large"},
			setup: func(t *testing.T) *mocks.MockIssuesInterface {
				m := mocks.NewMockIssuesInterface(t)
				m.On("AddLabelsToIssue", context.Background(), "owner", "repo", prNumber, []string{"large"}).
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
			err := client.AddLabels(context.Background(), pr, tc.labels)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
