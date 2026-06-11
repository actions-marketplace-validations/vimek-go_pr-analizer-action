package counter

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/vimek-go/pr-analizer-action/models"
)

var hunkHeaderRegex = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)

func CountDiff(
	patch string,
	baseLines, headLines []models.LineType,
	filename string,
	lang models.Language,
) (models.DiffStats, error) {
	stats := models.DiffStats{}
	isTest := false

	if lang.TestPattern != "" {
		fullPath := filepath.ToSlash(filename)
		matchFull, err := filepath.Match(lang.TestPattern, fullPath)
		if err != nil {
			return stats, errors.Wrapf(
				err,
				"failed to match full test file pattern: [%s] to path: [%s]",
				lang.TestPattern,
				fullPath,
			)
		}
		matchBase, err := filepath.Match(lang.TestPattern, filepath.Base(fullPath))
		if err != nil {
			return stats, errors.Wrapf(
				err,
				"failed to match base test file pattern: [%s] to path: [%s]",
				lang.TestPattern,
				fullPath,
			)
		}
		isTest = matchFull || matchBase
	}

	lines := strings.Split(patch, "\n")

	// 1-based line numbers
	var currentBaseLine int
	var currentHeadLine int

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			if base, head, ok := parseHunkHeader(line); ok {
				currentBaseLine = base
				currentHeadLine = head
			}
			continue
		}

		if len(line) == 0 {
			continue
		}

		firstChar := line[0]

		switch firstChar {
		case '-':
			// Removed line
			// Look up in baseLines (0-based index)
			if currentBaseLine > 0 && currentBaseLine <= len(baseLines) {
				lineType := baseLines[currentBaseLine-1]
				incrementStats(&stats, lineType, false, isTest)
			}
			currentBaseLine++
		case '+':
			// Added line
			// Look up in headLines
			if currentHeadLine > 0 && currentHeadLine <= len(headLines) {
				lineType := headLines[currentHeadLine-1]
				incrementStats(&stats, lineType, true, isTest)
			}
			currentHeadLine++
		case ' ':
			// Context line
			currentBaseLine++
			currentHeadLine++
		}
		// Ignore other lines (like "\ No newline at end of file")
	}

	// Calculate Net Stats
	stats.Code = stats.CodeAdded - stats.CodeRemoved
	stats.Comments = stats.CommentsAdded - stats.CommentsRemoved
	stats.Blanks = stats.BlanksAdded - stats.BlanksRemoved
	stats.Test = stats.TestAdded - stats.TestRemoved

	return stats, nil
}

func parseHunkHeader(line string) (baseStart, headStart int, ok bool) {
	matches := hunkHeaderRegex.FindStringSubmatch(line)
	if len(matches) == 0 {
		return 0, 0, false
	}
	baseStart, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, false
	}
	headStart, err = strconv.Atoi(matches[3])
	if err != nil {
		return 0, 0, false
	}
	return baseStart, headStart, true
}

func incrementStats(stats *models.DiffStats, lineType models.LineType, isAdded bool, isTest bool) {
	if isTest {
		if isAdded {
			stats.TestAdded++
		} else {
			stats.TestRemoved++
		}
		return
	}

	switch lineType {
	case models.LineCode:
		if isAdded {
			stats.CodeAdded++
		} else {
			stats.CodeRemoved++
		}
	case models.LineComment:
		if isAdded {
			stats.CommentsAdded++
		} else {
			stats.CommentsRemoved++
		}
	case models.LineBlank:
		if isAdded {
			stats.BlanksAdded++
		} else {
			stats.BlanksRemoved++
		}
	}
}
