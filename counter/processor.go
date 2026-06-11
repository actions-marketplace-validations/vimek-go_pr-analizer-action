package counter

import (
	"bufio"
	"io"
	"path/filepath"
	"strings"

	"github.com/vimek-go/pr-analizer-action/models"
)

func Count(r io.Reader, filePath string, lang models.Language) (models.Stats, error) {
	lineTypes, err := Analyze(r, lang)
	if err != nil {
		return models.Stats{}, err
	}

	stats := models.Stats{}
	isTest := false

	if lang.TestPattern != "" {
		fullPath := filepath.ToSlash(filePath)
		matchFull, err := filepath.Match(lang.TestPattern, fullPath)
		if err != nil {
			return stats, err
		}
		matchBase, err := filepath.Match(lang.TestPattern, filepath.Base(fullPath))
		if err != nil {
			return stats, err
		}
		isTest = matchFull || matchBase
	}

	for _, lt := range lineTypes {
		switch lt {
		case models.LineCode:
			stats.Code++
		case models.LineComment:
			stats.Comments++
		case models.LineBlank:
			stats.Blanks++
		}
	}

	if isTest {
		stats.Test = stats.Code + stats.Comments + stats.Blanks
		stats.Code = 0
		stats.Comments = 0
		stats.Blanks = 0
	}

	return stats, nil
}

// Analyze processes the content and returns the classification of each line.
// This preserves line-by-line context (like multi-line comments) which is crucial for diff parsing.
func Analyze(r io.Reader, lang models.Language) ([]models.LineType, error) {
	var results []models.LineType
	scanner := bufio.NewScanner(r)
	const maxCapacity = 1024 * 1024 // 1 MB
	buf := make([]byte, 0, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	inMultilineComment := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" {
			results = append(results, models.LineBlank)
			continue
		}

		if inMultilineComment {
			results = append(results, models.LineComment)
			if strings.Contains(line, lang.MultiLineCommentEnd) {
				inMultilineComment = false
			}
			continue
		}

		if lang.MultiLineCommentStart != "" && strings.Contains(line, lang.MultiLineCommentStart) {
			startIdx := strings.Index(line, lang.MultiLineCommentStart)
			// If code precedes the opening marker (e.g. `x := 1 /* comment */`), count as code.
			if startIdx > 0 && strings.TrimSpace(line[:startIdx]) != "" {
				results = append(results, models.LineCode)
			} else {
				results = append(results, models.LineComment)
			}
			if lang.MultiLineCommentEnd != "" && !strings.Contains(line, lang.MultiLineCommentEnd) {
				inMultilineComment = true
			}
			continue
		}

		if lang.LineComment != "" && strings.HasPrefix(trimmedLine, lang.LineComment) {
			results = append(results, models.LineComment)
			continue
		}

		results = append(results, models.LineCode)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
