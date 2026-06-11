package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vimek-go/pr-analizer-action/config"
	"github.com/vimek-go/pr-analizer-action/models"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config_*.yaml")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

func TestLoad(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		configContent  string
		useInvalidPath bool
		wantErr        bool
		errorContains  string
		validate       func(t *testing.T, cfg *models.Config)
	}{
		{
			name: "valid config",
			configContent: `
languages:
  - name: Go
    extensions: [".go"]
    line_comment: "//"
    multi_line_comment_start: "/*"
    multi_line_comment_end: "*/"
    test_pattern: "*_test.go"
`,
			wantErr: false,
			validate: func(t *testing.T, cfg *models.Config) {
				assert.Len(t, cfg.Languages, 1)
				assert.Equal(t, "Go", cfg.Languages[0].Name)
			},
		},
		{
			name:           "file not found",
			useInvalidPath: true,
			wantErr:        true,
			errorContains:  "no such file or directory",
		},
		{
			name:          "invalid YAML",
			configContent: ":::invalid yaml:::",
			wantErr:       true,
			errorContains: "config: no languages defined",
		},
		{
			name: "fails validation - duplicate languages",
			configContent: `
languages:
  - name: Go
    extensions: [".go"]
  - name: Go
    extensions: [".go"]
`,
			wantErr:       true,
			errorContains: "duplicate",
		},
		{
			name: "with label rules",
			configContent: `
languages:
  - name: Go
    extensions: [".go"]
label_rules:
  - label: size/L
    conditions:
      - total.code > 500
`,
			wantErr: false,
			validate: func(t *testing.T, cfg *models.Config) {
				assert.Len(t, cfg.LabelRules, 1)
				assert.Equal(t, "size/L", cfg.LabelRules[0].Label)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var path string
			if tt.useInvalidPath {
				path = filepath.Join(t.TempDir(), "nonexistent.yaml")
			} else {
				path = writeConfig(t, tt.configContent)
			}

			cfg, err := config.Load(path)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}
