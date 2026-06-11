package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/vimek-go/pr-analizer-action/config"
	"github.com/vimek-go/pr-analizer-action/event"
	"github.com/vimek-go/pr-analizer-action/gh"
	"github.com/vimek-go/pr-analizer-action/service"
)

func main() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}

	prDetails, err := event.GetPRDetails()
	if err != nil {
		log.Fatalf("Error getting PR details: %v", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "analizer_config.yaml" // Default config path
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	ignoreNotDefinedFlag, err := strconv.ParseBool(os.Getenv("IGNORE_NOT_DEFINED_LANGUAGE"))
	if err != nil {
		// This triggers for empty string or non parsable values
		// Default to true
		ignoreNotDefinedFlag = true
	}

	verboseLogging, err := strconv.ParseBool(os.Getenv("VERBOSE_LOGGING"))
	if err != nil {
		verboseLogging = false
	}

	ctx := context.Background()
	ghClient := gh.NewClient(ctx, githubToken)
	analyzer := service.NewAnalyzer(
		ghClient,
		cfg,
		service.WithIgnoreNotDefined(ignoreNotDefinedFlag),
		service.WithVerboseLogging(verboseLogging),
	)

	if err := analyzer.Run(context.Background(), prDetails); err != nil {
		log.Fatalf("Error running analysis: %v", err)
	}
}
