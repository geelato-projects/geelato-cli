package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/geelato/cli/cmd/initializer"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func NewPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "page",
		Short: "page(页面管理)",
		Long:  `管理应用页面配置，包括创建页面、编辑源码、发布等操作`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(pageCreateCmd)

	return cmd
}

var (
	pageType string
	pageDesc string
)

var pageCreateCmd = &cobra.Command{
	Use:   "create <page-name>",
	Short: "create(创建页面)",
	Long: `创建一个新的页面配置

示例:
  geelato page create userList
  geelato page create userList --type form
  geelato page create userList --desc "用户列表页面"`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var pageName string
		if len(args) > 0 {
			pageName = args[0]
		} else {
			input, err := prompt.Input("Enter page name (e.g., userList):")
			if err != nil {
				return fmt.Errorf("failed to input page name: %w", err)
			}
			pageName = strings.TrimSpace(input)
		}

		if pageName == "" {
			return fmt.Errorf("page name cannot be empty")
		}

		if err := createPage(pageName); err != nil {
			return fmt.Errorf("failed to create page: %w", err)
		}

		logger.Infof("Page '%s' created successfully!", pageName)
		logger.Info("")
		logger.Info("Created files:")
		logger.Info("  page/%s/%s.define.json", pageName, pageName)
		logger.Info("  page/%s/%s.source.json", pageName, pageName)
		logger.Info("  page/%s/%s.release.json", pageName, pageName)
		logger.Info("  page/%s/%s.preview.json", pageName, pageName)

		return nil
	},
}

func init() {
	pageCreateCmd.Flags().StringVarP(&pageType, "type", "t", "page", "Page type (page, form, dashboard)")
	pageCreateCmd.Flags().StringVarP(&pageDesc, "desc", "d", "", "Page description")
}

func createPage(pageName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	geelatoPath := filepath.Join(cwd, "geelato.json")
	if _, err := os.Stat(geelatoPath); os.IsNotExist(err) {
		return fmt.Errorf("current directory is not a valid Geelato application")
	}

	// Read geelato.json to get appId
	appId, err := getAppIdFromGeelatoJSON(geelatoPath)
	if err != nil {
		return fmt.Errorf("failed to read appId from geelato.json: %w", err)
	}

	pageDir := filepath.Join(cwd, "page", pageName)
	if err := os.MkdirAll(pageDir, 0755); err != nil {
		return fmt.Errorf("failed to create page directory: %w", err)
	}

	pageID := "pg_" + strings.ToLower(pageName)
	now := time.Now().Format(time.RFC3339)

	if pageDesc == "" {
		pageDesc = pageName + " page"
	}
	if pageType == "" {
		pageType = "page"
	}

	return initializer.CreatePageFiles(pageDir, pageName, pageID, now, appId, pageDesc, pageType)
}
