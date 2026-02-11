package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/geelato/cli/internal/app"
	"github.com/geelato/cli/internal/sync"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/progress"
	"github.com/spf13/cobra"
)

func NewPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pull from cloud(从云端拉取最新应用)",
		Long: `Pull the latest application from cloud platform

Example:
  geelato pull`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPull()
		},
	}
}

func runPull() error {
	logger.Info("Preparing to pull application from cloud...")

	cwd, err := os.Getwd()
	if err != nil {
		logger.Errorf("Failed to get working directory: %v", err)
		return err
	}

	geelatoPath := filepath.Join(cwd, "geelato.json")
	if _, err := os.Stat(geelatoPath); os.IsNotExist(err) {
		logger.Errorf("Not a Geelato application: geelato.json not found in %s", cwd)
		logger.Info("Please run 'geelato init' first to initialize an application.")
		return nil
	}

	// 从 geelato.json 读取 repo 配置
	appConfig, err := app.LoadAppConfig(cwd)
	if err != nil {
		logger.Errorf("Failed to load app config: %v", err)
		return err
	}

	repoURL := app.GetRepoFromConfig(appConfig)
	if repoURL == "" {
		logger.Error("Repo URL not configured. Please run 'geelato config repo <url>' first.")
		return nil
	}

	_, _, apiURL, err := app.ParseRepoURL(repoURL)
	if err != nil {
		logger.Errorf("Failed to parse repo URL: %v", err)
		return err
	}

	progressBar := progress.NewBar(100, "Pulling")
	if progressBar != nil {
		progressBar.Start()
	}

	svc, err := sync.NewSyncService(cwd, apiURL, "")
	if err != nil {
		if progressBar != nil {
			progressBar.Stop()
		}
		return err
	}
	if progressBar != nil {
		progressBar.Update(30)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	err = svc.Pull(ctx)
	if err != nil {
		if progressBar != nil {
			progressBar.Stop()
		}
		return err
	}
	if progressBar != nil {
		progressBar.Update(100)
		progressBar.Stop()
	}

	logger.Success("Application pulled successfully!")
	return nil
}
