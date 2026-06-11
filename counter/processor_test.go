package counter_test

import (
	"strings"
	"testing"

	"github.com/vimek-go/pr-analizer-action/counter"
	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var goLang = models.Language{
	Name:                  "Go",
	Extensions:            []string{".go"},
	LineComment:           "//",
	MultiLineCommentStart: "/*",
	MultiLineCommentEnd:   "*/",
	TestPattern:           "*_test.go",
}

func TestAnalyze_LineTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected []models.LineType
	}{
		{
			name:     "blank line",
			input:    "   \n",
			expected: []models.LineType{models.LineBlank},
		},
		{
			name:     "code line",
			input:    "x := 1\n",
			expected: []models.LineType{models.LineCode},
		},
		{
			name:     "line comment",
			input:    "// comment\n",
			expected: []models.LineType{models.LineComment},
		},
		{
			name:     "indented line comment",
			input:    "\t// comment\n",
			expected: []models.LineType{models.LineComment},
		},
		{
			name:     "single-line block comment",
			input:    "/* comment */\n",
			expected: []models.LineType{models.LineComment},
		},
		{
			name:  "multi-line block comment",
			input: "/*\nstill comment\n*/\n",
			expected: []models.LineType{
				models.LineComment,
				models.LineComment,
				models.LineComment,
			},
		},
		{
			name:  "code before block comment open is code",
			input: "x := compute() /* inline */\n",
			expected: []models.LineType{
				models.LineCode,
			},
		},
		{
			name:  "code before opening starts multi-line",
			input: "x := 1 /* start\nstill comment\n*/\n",
			expected: []models.LineType{
				models.LineCode,
				models.LineComment,
				models.LineComment,
			},
		},
		{
			name:  "mixed lines",
			input: "code\n// comment\n\ncode2\n",
			expected: []models.LineType{
				models.LineCode,
				models.LineComment,
				models.LineBlank,
				models.LineCode,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := counter.Analyze(strings.NewReader(tt.input), goLang)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestCount_NonTestFile(t *testing.T) {
	t.Parallel()
	input := "x := 1\n// comment\n\n"
	stats, err := counter.Count(strings.NewReader(input), "main.go", goLang)
	require.NoError(t, err)
	assert.Equal(t, 1, stats.Code)
	assert.Equal(t, 1, stats.Comments)
	assert.Equal(t, 1, stats.Blanks)
	assert.Equal(t, 0, stats.Test)
}

func TestCount_TestFile_BaseName(t *testing.T) {
	t.Parallel()
	input := "x := 1\n// comment\n\n"
	stats, err := counter.Count(strings.NewReader(input), "foo_test.go", goLang)
	require.NoError(t, err)
	assert.Equal(t, 0, stats.Code)
	assert.Equal(t, 0, stats.Comments)
	assert.Equal(t, 3, stats.Test)
}

func TestCount_TestFile_FullPath(t *testing.T) {
	t.Parallel()
	input := "x := 1\n"
	stats, err := counter.Count(strings.NewReader(input), "src/pkg/foo_test.go", goLang)
	require.NoError(t, err)
	assert.Equal(t, 1, stats.Test)
	assert.Equal(t, 0, stats.Code)
}
