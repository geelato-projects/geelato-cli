package mcp

import (
	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
)

var mcpInfoCmd = &cobra.Command{
	Use:   "info [name]",
	Short: "查看能力详情",
	Long: `查看指定 MCP 能力的详细信息。

示例：
  geelato mcp info database
  geelato mcp info auth`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInfo(args)
	},
}

func runInfo(args []string) error {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	logger.Infof("查看能力详情: %s", name)
	logger.Info("")
	logger.Info("能力详情功能正在开发中...")

	return nil
}
