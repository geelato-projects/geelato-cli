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

func NewPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push [message]",
		Short: "Push to cloud(推送变更到云端)",
		Long: `Push the current application to cloud platform

Example:
  geelato push "feat: add new model"
  geelato push`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			message := ""
			if len(args) > 0 {
				message = args[0]
			}
			return runPush(message)
		},
	}
}

func runPush(message string) error {
	logger.Info("Preparing to push application to cloud...")

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

	if message == "" {
		message = "Update application via CLI"
	}

	progressBar := progress.NewBar(100, "Pushing")
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

	err = svc.Push(ctx, message)
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

	logger.Success("Application pushed successfully!")
	return nil
}
