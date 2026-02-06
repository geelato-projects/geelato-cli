package mcp

import (
	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
)

var mcpSearchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "搜索可用能力",
	Long: `搜索云端平台可用的 MCP 能力。

示例：
  geelato mcp search database
  geelato mcp search auth`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSearch(args)
	},
}

func runSearch(args []string) error {
	keyword := ""
	if len(args) > 0 {
		keyword = args[0]
	}

	logger.Infof("搜索 MCP 能力: %s", keyword)
	logger.Info("")
	logger.Info("搜索功能正在开发中...")

	return nil
}
