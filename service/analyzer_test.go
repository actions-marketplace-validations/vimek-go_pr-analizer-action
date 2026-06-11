package service_test

import (
	"testing"

	"github.com/vimek-go/pr-analizer-action/models"
	"github.com/vimek-go/pr-analizer-action/service"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzer_GetLanguageForFile(t *testing.T) {
	type expected struct {
		found bool
		lang  models.Language
	}

	t.Parallel()
	priority := 0
	langs := []models.Language{
		{
			Name:            "test priority 0",
			Extensions:      []string{".test"},
			IncludePatterns: []string{"**/test/**", "**/best?/**"},
			ExcludePatterns: []string{"**/test/*/excluded/**"},
		},
		{
			Name:            "priority override",
			Priority:        &priority,
			Extensions:      []string{".test"},
			IncludePatterns: []string{"**/asd/test/**"},
			ExcludePatterns: []string{"**/test/*/excluded/**"},
		},
	}
	testCases := []struct {
		name     string
		path     string
		expected expected
	}{
		{
			name:     "priority override",
			path:     "/abs/asd/test/age/asd/er1/test.test",
			expected: expected{found: true, lang: langs[1]},
		},
		{
			name:     "explicitly excluded one level dir",
			path:     "/abs/test/age/excluded/er1/test.test",
			expected: expected{found: false, lang: models.LangNotDetected},
		},
		{
			name:     "2 level dir from ignore",
			path:     "/abs/test/age/er1/excluded/test.test",
			expected: expected{found: true, lang: langs[0]},
		},
		{
			name:     "fails as there is not additional char",
			path:     "/abs/best/age/er1/excluded/test.test",
			expected: expected{found: false, lang: models.LangNotDetected},
		},
		{
			name:     "included for dir name",
			path:     "/abs/best2/age/er1/excluded/test2.test",
			expected: expected{found: true, lang: langs[0]},
		},
		{
			name:     "excluded too many additional chars",
			path:     "/abs/best-add/age/er1/excluded/test2.test",
			expected: expected{found: false, lang: models.LangNotDetected},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			subject := service.NewAnalyzer(nil, &models.Config{Languages: langs})
			found, lang := subject.GetLanguageForFile(tc.path)
			assert.Equal(t, tc.expected.found, found)
			assert.Equal(t, tc.expected.lang, lang)
		})
	}
}
