package cmd

import (
	"os"
	"path/filepath"

	"github.com/geelato/cli/internal/app"
	"github.com/geelato/cli/pkg/logger"
	"github.com/spf13/cobra"
)

func NewValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate application(验证应用配置)",
		Long: `Validate the current Geelato application configuration.

This command checks:
- geelato.json configuration
- Required directories (meta/, api/, workflow/)
- Model definitions
- API scripts
- Workflow definitions

Example:
  geelato validate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate()
		},
	}
}

func runValidate() error {
	logger.Info("Validating application...")

	cwd, err := os.Getwd()
	if err != nil {
		logger.Errorf("Failed to get working directory: %v", err)
		return err
	}

	geelatoPath := filepath.Join(cwd, "geelato.json")
	if _, err := os.Stat(geelatoPath); os.IsNotExist(err) {
		logger.Errorf("Not a Geelato application: geelato.json not found in %s", cwd)
		logger.Info("Please run 'geelato init' first to initialize an application.")
	}

	validator := app.NewValidator()
	result, err := validator.Validate(cwd)
	if err != nil {
		logger.Errorf("Validation failed: %v", err)
		return err
	}

	logger.Info("Validation complete!")
	logger.Infof("Found %d models, %d APIs, %d workflows",
		result.Models, result.APIs, result.Workflows)

	if len(result.Errors) > 0 {
		logger.Error("Validation found errors:")
		for _, e := range result.Errors {
			logger.Errorf("  - %s", e)
		}
	} else {
		logger.Success("Application structure is valid!")
	}

	return nil
}
