package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/internal/config"
	"github.com/geelato/cli/internal/sync"
	"github.com/geelato/cli/pkg/logger"
)

func NewDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Show differences",
		Long: `Show differences between local and cloud

Example:
  geelato diff`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiff()
		},
	}
}

func runDiff() error {
	logger.Info("Checking differences between local and cloud...")

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

	svc, err := sync.NewSyncService(cwd, cfg.API.URL, cfg.API.Key)
	if err != nil {
		return err
	}

	diff, err := svc.GetDiff()
	if err != nil {
		return err
	}

	if len(diff.Added) == 0 && len(diff.Modified) == 0 && len(diff.Deleted) == 0 {
		logger.Info("No differences found. Local and cloud are in sync.")
		return nil
	}

	logger.Info("Differences found:")
	logger.Infof("  Added: %d", len(diff.Added))
	logger.Infof("  Modified: %d", len(diff.Modified))
	logger.Infof("  Deleted: %d", len(diff.Deleted))

	return nil
}
