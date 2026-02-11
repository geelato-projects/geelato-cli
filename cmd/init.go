package cmd

import (
	"fmt"
	"strings"

	"github.com/geelato/cli/cmd/initializer"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func initCmdFn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <app-name>",
		Short: "Initialize a new application(初始化新应用)",
		Long: `Initialize a new Geelato application

Examples:
  geelato init my-app`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var appName string
			if len(args) > 0 {
				appName = args[0]
			} else {
				input, err := prompt.Input("Enter app name (app-name):")
				if err != nil {
					return fmt.Errorf("failed to input app name: %w", err)
				}
				appName = strings.TrimSpace(input)
			}

			if appName == "" {
				return fmt.Errorf("app name cannot be empty")
			}

			appDir, err := initializer.InitializeApp(appName, "")
			if err != nil {
				return fmt.Errorf("failed to initialize app: %w", err)
			}

			logger.Infof("App '%s' initialized successfully!", appName)
			logger.Infof("App directory: %s", appDir)
			logger.Info("")
			logger.Info("Next steps:")
			logger.Info("  cd %s", appName)
			logger.Info("  geelato model create User        # Create your first model")
			logger.Info("  geelato api create getUserList  # Create your first API")
			logger.Info("  geelato app push                # Push to cloud")

			return nil
		},
	}

	return cmd
}
