# PR Analizer Action

A fast, configurable GitHub Action to analize in your Pull Requests. It analyzes changes (delta) and posts a detailed summary comment directly to your PR.

![PR Analizer Badge](https://img.shields.io/badge/PR-Analizer-blue)

<a href="https://www.buymeacoffee.com/vimekgo" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 36px !important;width: 130px !important;" ></a>

## 🚀 Features

*   **Delta Analysis:** Calculates the exact change in Code, Comments, Blank lines, and Test lines (e.g., `Code: +120`, `Comments: -5`).
*   **Per-Language Breakdown:** Detailed statistics for each language involved in the PR.
*   **Smart Commenting:** Posts a single comment on the PR and updates it on subsequent pushes (no spamming).
*   **Multi-Language Support:** Configure **any language** via YAML.
*   **PR Labeling:** Automatically adds labels to your Pull Request based on configurable analysis conditions.

## 📸 Example Output

| Language | Code | Comments | Blanks | Test |
| :--- | :--- | :--- | :--- | :--- |
| Go | +450 | +120 | +30 | +150 |
| Python | +20 | +5 | +2 | +10 |
| **TOTAL** | **+470** | **+125** | **+32** | **+160** |

*(Actual output is a beautifully formatted Markdown table in your PR comment)*

## 📦 Usage

1. Create a Github token with read access to the repository and write access to pull requests to add labels and comments. The example assumes you have a `GITHUB_TOKEN` secret set in your repository.
2. Create a workflow file (e.g., `.github/workflows/analizer_config.yml`) in your repository:
3. Copy configuration from [/configuration_examples](configuration_examples) and paste it into your repository. Default path is `analizer_config.yaml`.

```yaml
name: PR Analizer

on:
  pull_request:
    types: [opened, synchronize, reopened]

permissions:
  contents: read
  pull-requests: write # Required to post comments AND add labels

jobs:
  analyze-code:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Analize code
        uses: vimek-go/pr-analizer-action@1.0.0
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          # Optional: Path to your custom config
          # config_path: analizer_config.yaml - default value
```

## ⚙️ Configuration

The action includes example configurations for common programming languages; just copy the .yaml configuration file. To customise languages, ignore patterns, test file patterns, or add PR labelling rules, you can either modify one of the examples in [/configuration_examples](configuration_examples) or provide your own custom configuration file.

```yaml
global_ignore:
  # Supports glob patterns (*, **, ?). Matched against the full file path.
  # Use ** to match files within a directory at any depth.
  - "node_modules/**"
  - "vendor/**"
  - ".git/**"

languages:
  - name: "Go"
    extensions:
      - ".go"
    line_comment: "//"
    multi_line_comment_start: "/*"
    multi_line_comment_end: "*/"
    test_pattern: "*_test.go"
    # Optional: Patterns to include/exclude files for this language.
    # Supports glob patterns (*, **, ?).
    include_patterns:
      - "src/**/*.go"
    exclude_patterns:
      - "vendor/**"
    # Optional: Priority for language detection (higher value = higher priority).
    # Useful when multiple languages match the same file.
    priority: 10

  - name: "Makefile"
    file_names:
      - "Makefile"
    line_comment: "#"

label_rules:
  # Example: Label PRs with small amount of code changes
  - label: "size/small"
    conditions:
      - "total.code <= 200" # If total added/modified code lines are less or equal than 200

  # Example: Label PRs that add Go code but have no Go test changes
  - label: "needs-go-tests"
    conditions:
      - "language.Go.code > 0"  # If Go code lines were added/modified
      - "language.Go.test == 0" # And no Go test lines were added/modified

  # Example: Label PRs with significant comment changes
  - label: "documentation-heavy"
    conditions:
      - "total.comments > 100"

  # Example: Label PRs that remove a lot of code
  - label: "refactor/cleanup"
    conditions:
      - "total.code < -200" # Note: Negative values for removals

# Available variables for conditions:
# - total.code: Total change (net) in code lines across all languages.
# - total.comments: Total change (net) in comment lines across all languages.
# - total.blanks: Total change (net) in blank lines across all languages.
# - total.test: Total change (net) in test lines across all languages.
#
# You can also target specifically added or removed lines:
# - total.code_added, total.code_removed
# - total.comments_added, total.comments_removed
# - total.blanks_added, total.blanks_removed
# - total.test_added, total.test_removed
#
# - language.{LanguageName}.code: Change in code lines for a specific language (e.g., 'language.Go.code').
# - language.{LanguageName}.comments: Change in comment lines for a specific language.
# - language.{LanguageName}.blanks: Change in blank lines for a specific language.
# - language.{LanguageName}.test: Change in test lines for a specific language.
#
# - language.{LanguageName}.code_added, language.{LanguageName}.code_removed
# ... (and so on for other metrics)
#
# Supported operators: >, <, >=, <=, =, ==, !=

```

Changing the config path is as simple as adding the following to your workflow:

```yaml
with:
  config_path: configs/my_custom_config.yaml
```

## 📥 Inputs

| Input | Description | Required | Default |
| :--- | :--- | :--- | :--- |
| `github_token` | The GitHub Token to authenticate with the API. | **Yes** | `${{ github.token }}` |
| `config_path` | Path to the YAML configuration file relative to the repo root. | No | `.github/workflows/analizer_config.yaml` |
| `ignore_not_defined_language` | If `true`, files with languages not defined in the configuration will be ignored. | No | `true` |
| `verbose_logging` | If `true`, enables detailed per-file matching logs in the action output. | No | `false` |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
