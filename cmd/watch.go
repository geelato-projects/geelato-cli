package cmd

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/internal/config"
	"github.com/geelato/cli/internal/watcher"
	"github.com/geelato/cli/pkg/logger"
)

func NewWatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Watch for changes",
		Long: `Watch for file changes and auto-sync to cloud

Example:
  geelato watch`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatch()
		},
	}
}

func runWatch() error {
	logger.Info("Starting file watcher...")

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

	watch, err := watcher.NewWatcher(cwd)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Watching for file changes. Press Ctrl+C to stop.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-sigChan
		logger.Info("Shutting down watcher...")
		cancel()
	}()

	err = watch.Start(ctx)
	if err != nil {
		return err
	}

	logger.Info("Watcher stopped.")
	return nil
}
