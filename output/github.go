package output

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vimek-go/pr-analizer-action/models"
)

func GenerateMarkdownDiff(report models.DiffReport) string {
	var sb strings.Builder
	sb.WriteString("### PR Analizer\n")

	// Table Header
	sb.WriteString("| Language | Code | Comments | Blanks | Test |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- | :--- |\n")

	// Sort languages
	var languages []string
	for lang := range report.ByLanguage {
		languages = append(languages, lang)
	}
	sort.Strings(languages)

	// Rows for each language
	for _, lang := range languages {
		diff := report.ByLanguage[lang]
		fmt.Fprintf(&sb, "| %s | %s | %s | %s | %s |\n",
			lang,
			formatCell(diff.CodeAdded, diff.CodeRemoved),
			formatCell(diff.CommentsAdded, diff.CommentsRemoved),
			formatCell(diff.BlanksAdded, diff.BlanksRemoved),
			formatCell(diff.TestAdded, diff.TestRemoved))
	}

	// Total Row
	diff := report.Total
	fmt.Fprintf(&sb, "| **TOTAL** | **%s** | **%s** | **%s** | **%s** |\n",
		formatCell(diff.CodeAdded, diff.CodeRemoved),
		formatCell(diff.CommentsAdded, diff.CommentsRemoved),
		formatCell(diff.BlanksAdded, diff.BlanksRemoved),
		formatCell(diff.TestAdded, diff.TestRemoved))

	return sb.String()
}

func formatCell(added, removed int) string {
	if added == 0 && removed == 0 {
		return "--"
	}

	addedStr := "0"
	if added > 0 {
		addedStr = fmt.Sprintf("+%d", added)
	}
	removedStr := "0"
	if removed > 0 {
		removedStr = fmt.Sprintf("-%d", removed)
	}

	return fmt.Sprintf("%s / %s", addedStr, removedStr)
}
