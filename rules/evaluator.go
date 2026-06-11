package rules

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/pkg/errors"
)

var (
	conditionRegex = regexp.MustCompile(`^([a-zA-Z0-9_.]+)\s*([<>=!]+)\s*(\d+)$`)

	errLanguageNotFound = errors.New("language not found")
)

func Evaluate(rules []models.LabelRule, report models.DiffReport) ([]string, error) {
	var appliedLabels []string

	for _, rule := range rules {
		match := true
		for _, cond := range rule.Conditions {
			result, err := evaluateCondition(cond, report)
			if err != nil {
				if errors.Is(err, errLanguageNotFound) {
					// Skip this rule if a language is missing in the report
					match = false
					break
				}
				return nil, errors.Wrapf(err, "error evaluating condition '%s' for label '%s'", cond, rule.Label)
			}
			if !result {
				match = false
				break
			}
		}
		if match {
			appliedLabels = append(appliedLabels, rule.Label)
		}
	}

	return appliedLabels, nil
}

func evaluateCondition(condition string, report models.DiffReport) (bool, error) {
	matches := conditionRegex.FindStringSubmatch(condition)
	if len(matches) != 4 {
		return false, errors.Errorf("invalid condition format: %q", condition)
	}

	lhsPath := matches[1]
	operator := matches[2]
	rhsVal, err := strconv.Atoi(matches[3])
	if err != nil {
		return false, errors.Errorf("invalid number: %s", matches[3])
	}

	lhsVal, err := getLHSValue(lhsPath, report)
	if err != nil {
		// If a language is missing in the report, treating it as error
		// We need to check for this error not to return false positives
		return false, err
	}

	switch operator {
	case ">":
		return lhsVal > rhsVal, nil
	case "<":
		return lhsVal < rhsVal, nil
	case ">=":
		return lhsVal >= rhsVal, nil
	case "<=":
		return lhsVal <= rhsVal, nil
	case "=":
		return lhsVal == rhsVal, nil // Support both = and ==
	case "==":
		return lhsVal == rhsVal, nil
	case "!=":
		return lhsVal != rhsVal, nil
	default:
		return false, errors.Errorf("unknown operator: %s", operator)
	}
}

func getLHSValue(path string, report models.DiffReport) (int, error) {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return 0, errors.Errorf("invalid path: %s", path)
	}

	switch parts[0] {
	case "total":
		return getTotalValue(path, parts, report)
	case "language":
		return getLanguageValue(path, parts, report)
	default:
		return 0, errors.Errorf("unknown category: %s", parts[0])
	}
}

func getTotalValue(path string, parts []string, report models.DiffReport) (int, error) {
	if len(parts) != 2 {
		return 0, errors.Errorf("invalid total path: %s", path)
	}
	return getFieldValue(report.Total, parts[1])
}

func getLanguageValue(path string, parts []string, report models.DiffReport) (int, error) {
	if len(parts) != 3 {
		return 0, errors.Errorf("invalid language path: %s", path)
	}
	langName, field := parts[1], parts[2]

	stats, ok := report.ByLanguage[langName]
	if !ok {
		stats, ok = findLanguageCaseInsensitive(report.ByLanguage, langName)
		if !ok {
			return 0, errLanguageNotFound
		}
	}
	return getFieldValue(stats, field)
}

func findLanguageCaseInsensitive(byLanguage map[string]models.DiffStats, langName string) (models.DiffStats, bool) {
	lower := strings.ToLower(langName)
	for k, v := range byLanguage {
		if strings.ToLower(k) == lower {
			return v, true
		}
	}
	return models.DiffStats{}, false
}

func getFieldValue(stats models.DiffStats, field string) (int, error) {
	switch strings.ToLower(field) {
	case "code":
		return stats.Code, nil
	case "comments":
		return stats.Comments, nil
	case "blanks":
		return stats.Blanks, nil
	case "test", "tests":
		return stats.Test, nil

	// Added fields
	case "code_added":
		return stats.CodeAdded, nil
	case "comments_added":
		return stats.CommentsAdded, nil
	case "blanks_added":
		return stats.BlanksAdded, nil
	case "test_added", "tests_added":
		return stats.TestAdded, nil

	// Removed fields
	case "code_removed":
		return stats.CodeRemoved, nil
	case "comments_removed":
		return stats.CommentsRemoved, nil
	case "blanks_removed":
		return stats.BlanksRemoved, nil
	case "test_removed", "tests_removed":
		return stats.TestRemoved, nil

	default:
		return 0, errors.Errorf("unknown field: %s", field)
	}
}
