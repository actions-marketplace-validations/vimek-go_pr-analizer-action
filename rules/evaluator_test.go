package rules_test

import (
	"testing"

	"github.com/vimek-go/pr-analizer-action/models"
	"github.com/vimek-go/pr-analizer-action/rules"

	"github.com/stretchr/testify/assert"
)

func Test_Evaluate(t *testing.T) {
	report := models.DiffReport{
		Total: models.DiffStats{
			Code:            600,
			Test:            50,
			CodeAdded:       800,
			CodeRemoved:     200,
			CommentsAdded:   100,
			CommentsRemoved: 50,
		},
		ByLanguage: map[string]models.DiffStats{
			"Go": {
				Code:            200,
				Test:            0,
				CodeAdded:       250,
				CodeRemoved:     50,
				CommentsRemoved: 20,
			},
			"Python": {
				Code:      400,
				Test:      50,
				CodeAdded: 400,
			},
		},
	}

	tests := []struct {
		name     string
		rules    []models.LabelRule
		expected []string
	}{
		{
			name: "total code > 500",
			rules: []models.LabelRule{
				{Label: "size/L", Conditions: []string{"total.code > 500"}},
			},
			expected: []string{"size/L"},
		},
		{
			name: "total code_added > 700",
			rules: []models.LabelRule{
				{Label: "lots-of-additions", Conditions: []string{"total.code_added > 700"}},
			},
			expected: []string{"lots-of-additions"},
		},
		{
			name: "go code_removed > 30",
			rules: []models.LabelRule{
				{Label: "cleanup-go", Conditions: []string{"language.go.code_removed > 30"}},
			},
			expected: []string{"cleanup-go"},
		},
		{
			name: "go code > 100 and no tests",
			rules: []models.LabelRule{
				{Label: "needs-tests", Conditions: []string{"language.go.code > 100", "language.go.test == 0"}},
			},
			expected: []string{"needs-tests"},
		},
		{
			name: "python code < 100 (false)",
			rules: []models.LabelRule{
				{Label: "python-small", Conditions: []string{"language.python.code < 100"}},
			},
			expected: nil,
		},
		{
			name: "missing language treated as err",
			rules: []models.LabelRule{
				{Label: "rust-zero", Conditions: []string{"language.rust.code == 0"}},
			},
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := rules.Evaluate(tc.rules, report)
			assert.NoError(t, err)

			assert.Len(t, actual, len(tc.expected))
			for i, v := range actual {
				assert.Equal(t, tc.expected[i], v)
			}
		})
	}
}
