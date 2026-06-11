package models_test

import (
	"testing"

	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validLang(name string) models.Language {
	return models.Language{Name: name, Extensions: []string{".go"}}
}

func TestConfig_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		config    models.Config
		errString string
	}{
		{
			name:      "no languages",
			config:    models.Config{},
			errString: "no languages defined",
		},
		{
			name:      "empty language name",
			config:    models.Config{Languages: []models.Language{{Extensions: []string{".go"}}}},
			errString: "name is required",
		},
		{
			name:      "reserved name",
			config:    models.Config{Languages: []models.Language{{Name: "Not detected", Extensions: []string{".go"}}}},
			errString: "reserved",
		},
		{
			name:      "duplicate names",
			config:    models.Config{Languages: []models.Language{validLang("Python"), validLang("Python")}},
			errString: "duplicate",
		},
		{
			name:      "no extensions or filenames",
			config:    models.Config{Languages: []models.Language{{Name: "Go"}}},
			errString: "extension or filename",
		},
		{
			name: "unpaired multi-line comment start only",
			config: models.Config{Languages: []models.Language{{
				Name:                  "Go",
				Extensions:            []string{".go"},
				MultiLineCommentStart: "/*",
			}}},
			errString: "multi_line_comment_start",
		},
		{
			name: "unpaired multi-line comment end only",
			config: models.Config{Languages: []models.Language{{
				Name:                "Go",
				Extensions:          []string{".go"},
				MultiLineCommentEnd: "*/",
			}}},
			errString: "multi_line_comment_start",
		},
		{
			name: "empty label name",
			config: models.Config{
				Languages:  []models.Language{validLang("Go")},
				LabelRules: []models.LabelRule{{Label: ""}},
			},
			errString: "label name is required",
		},
		{
			name:   "valid minimal config",
			config: models.Config{Languages: []models.Language{validLang("Go")}},
		},
		{
			name: "valid config with all fields",
			config: models.Config{
				Languages: []models.Language{{
					Name:                  "Go",
					Extensions:            []string{".go"},
					LineComment:           "//",
					MultiLineCommentStart: "/*",
					MultiLineCommentEnd:   "*/",
					TestPattern:           "*_test.go",
				}},
				LabelRules: []models.LabelRule{{Label: "size/L", Conditions: []string{"total.code > 500"}}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.config.Validate()
			if tc.errString != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
