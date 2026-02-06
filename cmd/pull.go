package cmd

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/internal/config"
	"github.com/geelato/cli/internal/sync"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/progress"
)

func NewPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pull from cloud",
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

	cfg := config.Get()
	if cfg == nil {
		logger.Error("Configuration not loaded. Please run 'geelato config set api.url <url>' first.")
		return nil
	}

	if cfg.API.URL == "" {
		logger.Error("API URL not configured. Please run 'geelato config set api.url <url>' first.")
		return nil
	}

	progressBar := progress.NewBar(100, "Pulling")
	if progressBar != nil {
		progressBar.Start()
	}

	svc, err := sync.NewSyncService(cwd, cfg.API.URL, cfg.API.Key)
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
