package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geelato/cli/cmd/initializer"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func NewApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "api(API管理)",
		Long:  `管理 API 脚本，包括创建、测试、运行等操作`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(apiCreateCmd)
	cmd.AddCommand(NewApiTestCmd())
	cmd.AddCommand(NewApiRunCmd())

	return cmd
}

var (
	apiType string
)

var apiCreateCmd = &cobra.Command{
	Use:   "create <api-name>",
	Short: "创建API",
	Long: `创建一个新的 API 脚本

示例:
  geelato api create getUserList
  geelato api create getUserList -t js
  geelato api create saveUser -t python
  geelato api create myHandler -t go`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var apiName string
		if len(args) > 0 {
			apiName = args[0]
		} else {
			input, err := prompt.Input("Enter API name (e.g., getUserList):")
			if err != nil {
				return fmt.Errorf("failed to input API name: %w", err)
			}
			apiName = strings.TrimSpace(input)
		}

		if apiName == "" {
			return fmt.Errorf("API name cannot be empty")
		}

		if apiType == "" {
			apiType = "js"
		}

		if err := createAPI(apiName, apiType); err != nil {
			return fmt.Errorf("failed to create API: %w", err)
		}

		logger.Infof("API '%s' created successfully!", apiName)
		logger.Info("")
		logger.Info("Created file:")
		logger.Info("  api/%s.api.%s", apiName, apiType)

		return nil
	},
}

func init() {
	apiCreateCmd.Flags().StringVarP(&apiType, "type", "t", "", "API type (js, python, go)")
}

func createAPI(apiName, apiType string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	geelatoPath := filepath.Join(cwd, "geelato.json")
	if _, err := os.Stat(geelatoPath); os.IsNotExist(err) {
		return fmt.Errorf("current directory is not a valid Geelato application")
	}

	apiDir := filepath.Join(cwd, "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create API directory: %w", err)
	}

	fileName := fmt.Sprintf("%s.api.%s", apiName, apiType)
	filePath := filepath.Join(apiDir, fileName)

	return initializer.CreateAPIFile(filePath, apiName, apiType)
}

func NewApiTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test <api-file>",
		Short: "test(测试API)",
		Long:  `测试 API 脚本`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("API test functionality not yet implemented")
			return nil
		},
	}
}

func NewApiRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run <api-file>",
		Short: "run(运行API)",
		Long:  `运行 API 脚本`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("API run functionality not yet implemented")
			return nil
		},
	}
}
