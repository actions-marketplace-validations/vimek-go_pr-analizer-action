package models

import (
	"path/filepath"

	"github.com/pkg/errors"
)

const NotDetected = "Not detected"

var LangNotDetected = Language{Name: NotDetected}

type Config struct {
	GlobalIgnore []string    `yaml:"global_ignore"`
	Languages    []Language  `yaml:"languages"`
	LabelRules   []LabelRule `yaml:"label_rules"`
}

func (c *Config) Validate() error {
	if len(c.Languages) == 0 {
		return errors.New("config: no languages defined")
	}

	seen := make(map[string]struct{}, len(c.Languages))
	for i, lang := range c.Languages {
		if lang.Name == "" {
			return errors.Errorf("language[%d]: name is required", i)
		}
		if lang.Name == NotDetected {
			return errors.Errorf("language name '%s' is reserved", lang.Name)
		}
		if _, dup := seen[lang.Name]; dup {
			return errors.Errorf("duplicate language name: '%s'", lang.Name)
		}
		seen[lang.Name] = struct{}{}

		if len(lang.Extensions) == 0 && len(lang.FileNames) == 0 {
			return errors.Errorf("language '%s': must define at least one extension or filename", lang.Name)
		}

		if (lang.MultiLineCommentStart == "") != (lang.MultiLineCommentEnd == "") {
			return errors.Errorf(
				"language '%s': multi_line_comment_start and multi_line_comment_end must both be set or both be empty",
				lang.Name,
			)
		}

		if lang.TestPattern != "" {
			if _, err := filepath.Match(lang.TestPattern, ""); err != nil {
				return errors.Errorf("language '%s': invalid test_pattern '%s': %v", lang.Name, lang.TestPattern, err)
			}
		}
	}

	for _, rule := range c.LabelRules {
		if rule.Label == "" {
			return errors.New("label_rules: label name is required")
		}
	}

	return nil
}

type LabelRule struct {
	Label      string   `yaml:"label"`
	Conditions []string `yaml:"conditions"`
}

type Language struct {
	Name                  string   `yaml:"name"`
	Extensions            []string `yaml:"extensions"`
	FileNames             []string `yaml:"file_names"`
	LineComment           string   `yaml:"line_comment"`
	MultiLineCommentStart string   `yaml:"multi_line_comment_start"`
	MultiLineCommentEnd   string   `yaml:"multi_line_comment_end"`
	TestPattern           string   `yaml:"test_pattern"`
	IncludePatterns       []string `yaml:"include_patterns"`
	ExcludePatterns       []string `yaml:"exclude_patterns"`
	Priority              *int     `yaml:"priority"`
}
