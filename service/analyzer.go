package service

import (
	"context"
	"log"
	"math"
	"path/filepath"
	"slices"
	"strings"

	"github.com/vimek-go/pr-analizer-action/counter"
	"github.com/vimek-go/pr-analizer-action/enum"
	"github.com/vimek-go/pr-analizer-action/models"
	"github.com/vimek-go/pr-analizer-action/output"
	"github.com/vimek-go/pr-analizer-action/rules"

	"github.com/pkg/errors"
)

type Analyzer struct {
	client           Client
	config           *models.Config
	ignoreNotDefined bool
	verboseLogging   bool
}

type AnalyzerOption func(*Analyzer)

func NewAnalyzer(client Client, cfg *models.Config, options ...AnalyzerOption) *Analyzer {
	a := &Analyzer{
		client: client,
		config: cfg,
	}
	for _, option := range options {
		option(a)
	}
	return a
}

func WithIgnoreNotDefined(ignore bool) AnalyzerOption {
	return func(a *Analyzer) {
		a.ignoreNotDefined = ignore
	}
}

func WithVerboseLogging(verbose bool) AnalyzerOption {
	return func(a *Analyzer) {
		a.verboseLogging = verbose
	}
}

func (a *Analyzer) Run(ctx context.Context, pr *models.PullRequest) error {
	changedFiles, err := a.client.GetChangedFiles(ctx, pr)
	if err != nil {
		return errors.Wrapf(err, "getting changed files for PR #%d", pr.Number)
	}

	diffReport := models.DiffReport{
		ByLanguage: make(map[string]models.DiffStats),
	}

	for _, file := range changedFiles {
		// check if file is not ignored globally
		// Normalize to forward slashes for consistent matching
		normalizedPath := filepath.ToSlash(file.Filename)
		if a.matchesAnyPattern(normalizedPath, a.config.GlobalIgnore) {
			a.logWithCheckf("file %s - ignored globally", normalizedPath)
			continue
		}

		found, lang := a.GetLanguageForFile(normalizedPath)
		if !found && a.ignoreNotDefined {
			continue
		}

		diff, err := a.processFile(ctx, pr, file, lang)
		if err != nil {
			log.Printf("Warning: processing file %s: %v", file.Filename, err)
			continue
		}
		if diff == nil {
			continue
		}

		diffReport.Add(lang.Name, *diff)
	}

	return a.generateAndPostReport(ctx, pr, diffReport)
}

func (a *Analyzer) processFile(
	ctx context.Context,
	pr *models.PullRequest,
	file models.ChangedFile,
	lang models.Language,
) (*models.DiffStats, error) {
	switch file.Status {
	case "removed":
		stats, err := a.fetchAndCount(ctx, pr.Owner, pr.Repo, file.Filename, pr.BaseBranch, lang)
		if err != nil {
			return nil, err
		}
		diff := stats.ToModelDiff(enum.DiffTypes.Removed())
		return &diff, nil

	case "added":
		stats, err := a.fetchAndCount(ctx, pr.Owner, pr.Repo, file.Filename, pr.HeadBranch, lang)
		if err != nil {
			return nil, err
		}
		diff := stats.ToModelDiff(enum.DiffTypes.Added())
		return &diff, nil

	case "modified", "renamed":
		patch := file.Patch
		if patch == "" {
			if file.Status == "renamed" {
				return nil, nil
			}
			log.Printf("Warning: No patch available for %s. Skipping.", file.Filename)
			return nil, nil
		}

		baseLines, err := a.fetchAndAnalyze(ctx, pr.Owner, pr.Repo, file.Filename, pr.BaseBranch, lang)
		if err != nil {
			return nil, errors.Wrap(err, "analyzing base")
		}

		headLines, err := a.fetchAndAnalyze(ctx, pr.Owner, pr.Repo, file.Filename, pr.HeadBranch, lang)
		if err != nil {
			return nil, errors.Wrap(err, "analyzing head")
		}

		diff, err := counter.CountDiff(patch, baseLines, headLines, file.Filename, lang)
		if err != nil {
			return nil, errors.Wrap(err, "counting diff")
		}
		return &diff, nil

	default:
		log.Printf("Info: Skipping file %s with status %s", file.Filename, file.Status)
		return nil, nil
	}
}

func (a *Analyzer) fetchAndCount(
	ctx context.Context,
	owner, repo, path, ref string,
	lang models.Language,
) (models.Stats, error) {
	reader, err := a.client.GetFileContent(ctx, owner, repo, path, ref)
	if err != nil {
		return models.Stats{}, errors.Wrap(err, "error getting content")
	}
	defer reader.Close()

	return counter.Count(reader, path, lang)
}

func (a *Analyzer) fetchAndAnalyze(
	ctx context.Context,
	owner, repo, path, ref string,
	lang models.Language,
) ([]models.LineType, error) {
	reader, err := a.client.GetFileContent(ctx, owner, repo, path, ref)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting content")
	}
	defer reader.Close()

	return counter.Analyze(reader, lang)
}

func (a *Analyzer) generateAndPostReport(ctx context.Context, pr *models.PullRequest, report models.DiffReport) error {
	commentBody := output.GenerateMarkdownDiff(report)
	if err := a.client.PostOrUpdateComment(ctx, pr, commentBody); err != nil {
		return errors.Wrap(err, "posting/updating comment")
	}

	if len(a.config.LabelRules) == 0 {
		return nil
	}

	labelsToAdd, err := rules.Evaluate(a.config.LabelRules, report)
	if err != nil {
		return errors.Wrap(err, "evaluating label rules")
	}
	if len(labelsToAdd) > 0 {
		if err := a.client.AddLabels(ctx, pr, labelsToAdd); err != nil {
			return errors.Wrap(err, "adding labels")
		}
		log.Printf("Added labels: %v", labelsToAdd)
	}

	log.Println("Analysis complete and comment posted/updated.")
	return nil
}

func (a *Analyzer) GetLanguageForFile(normalizedPath string) (bool, models.Language) {
	ext := strings.ToLower(filepath.Ext(normalizedPath))
	fileName := filepath.Base(normalizedPath)

	bestPriority := math.MinInt
	bestLang := models.LangNotDetected
	found := false

	for _, lang := range a.config.Languages {
		// Check if file is explicitly excluded
		if a.matchesAnyPattern(normalizedPath, lang.ExcludePatterns) {
			a.logWithCheckf("file %s - ignored for language name: [%s]", normalizedPath, lang.Name)
			continue
		}

		// Check extension, filename, or include pattern
		// If include pattern not empty then we need to match only included
		matched := (slices.Contains(lang.Extensions, ext) ||
			slices.Contains(lang.FileNames, fileName)) &&
			(len(lang.IncludePatterns) == 0 || a.matchesAnyPattern(normalizedPath, lang.IncludePatterns))

		if !matched {
			a.logWithCheckf("file %s - not matched for language name: [%s]", normalizedPath, lang.Name)
			continue
		}

		priority := math.MinInt
		if lang.Priority != nil {
			priority = *lang.Priority
		}

		if !found || priority > bestPriority {
			a.logWithCheckf("file %s - matched for language name: [%s]", normalizedPath, lang.Name)
			bestPriority = priority
			bestLang = lang
			found = true
		}
	}

	if !found {
		a.logWithCheckf("file %s - not matched for any language name", normalizedPath)
	}

	return found, bestLang
}

// matchesAnyPattern returns true if the path matches any of the given glob patterns.
// Supports ** for multi-level directory matching and * for single-level matching.
func (a *Analyzer) matchesAnyPattern(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if a.matchGlob(pattern, path) {
			return true
		}
	}
	return false
}

// matchGlob matches a path against a glob pattern supporting:
//   - **  matches zero or more directories
//   - *   matches a single path segment (no slashes)
//   - ?   matches a single non-slash character
func (a *Analyzer) matchGlob(pattern, path string) bool {
	// Split pattern and path into segments by "/"
	patternParts := strings.Split(filepath.ToSlash(pattern), "/")
	pathParts := strings.Split(filepath.ToSlash(path), "/")

	return a.matchSegments(patternParts, pathParts)
}

func (a *Analyzer) matchSegments(pattern, path []string) bool {
	for len(pattern) > 0 {
		seg := pattern[0]

		if seg == "**" {
			// Consume consecutive ** segments
			for len(pattern) > 0 && pattern[0] == "**" {
				pattern = pattern[1:]
			}
			// If ** is the last pattern segment, it matches everything remaining
			if len(pattern) == 0 {
				return true
			}
			// Try matching the rest of the pattern at every remaining position
			for i := 0; i <= len(path); i++ {
				if a.matchSegments(pattern, path[i:]) {
					return true
				}
			}
			return false
		}

		// No more path segments but pattern still has non-** parts
		if len(path) == 0 {
			return false
		}

		// Match the current segment using filepath.Match (handles * and ? within a segment)
		matched, err := filepath.Match(seg, path[0])
		if err != nil || !matched {
			return false
		}

		pattern = pattern[1:]
		path = path[1:]
	}

	return len(path) == 0
}

func (a *Analyzer) logWithCheckf(format string, args ...any) {
	if a.verboseLogging {
		log.Printf("[Analyzer] "+format, args...)
	}
}
