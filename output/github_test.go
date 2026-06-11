package output_test

import (
	"strings"
	"testing"

	"github.com/vimek-go/pr-analizer-action/models"
	"github.com/vimek-go/pr-analizer-action/output"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMarkdownDiff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		report        models.DiffReport
		expected      []string
		expectedOrder []string
	}{
		{
			name:   "header",
			report: models.DiffReport{},
			expected: []string{
				"### PR Analizer",
				"TOTAL",
			},
		},
		{
			name:   "empty report",
			report: models.DiffReport{},
			expected: []string{
				"--",
			},
		},
		{
			name: "single Llanguage",
			report: models.DiffReport{
				Total: models.DiffStats{CodeAdded: 10, CodeRemoved: 2},
				ByLanguage: map[string]models.DiffStats{
					"Go": {CodeAdded: 10, CodeRemoved: 2},
				},
			},
			expected: []string{
				"Go",
				"+10",
				"-2",
			},
		},
		{
			name: "languages sorted alphabetically",
			report: models.DiffReport{
				ByLanguage: map[string]models.DiffStats{
					"Python": {CodeAdded: 5},
					"Go":     {CodeAdded: 10},
				},
			},
			expectedOrder: []string{
				"Go",
				"Python",
			},
		},
		{
			name: "zero added or removed",
			report: models.DiffReport{
				ByLanguage: map[string]models.DiffStats{
					"Go": {CodeAdded: 5},
				},
			},
			expected: []string{
				"+5 / 0",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := output.GenerateMarkdownDiff(tc.report)

			for _, expected := range tc.expected {
				assert.Contains(t, actual, expected)
			}

			if len(tc.expectedOrder) > 0 {
				// We need to remove first 2 header lines
				actualLanguages := strings.Split(actual, "\n")[3:]
				// and there are 2 lines at the end for total
				assert.Len(t, actualLanguages, len(tc.expectedOrder)+2)
				for i := range tc.expectedOrder {
					assert.Contains(t, actualLanguages[i], tc.expectedOrder[i])
				}
			}
		})
	}
}
