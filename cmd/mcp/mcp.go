package mcp

import (
	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
)

var McpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP 平台能力管理",
	Long: `MCP（Model-Context-Protocol）平台能力管理命令。

MCP 功能用于管理 Geelato 平台的能力同步：
  geelato mcp list        列出可用能力
  geelato mcp sync       同步能力到本地
  geelato mcp search     搜索可用能力
  geelato mcp info       查看能力详情

使用 "geelato mcp [command] --help" 查看命令帮助。`,
}

func init() {
	McpCmd.AddCommand(mcpListCmd, mcpSyncCmd, mcpSearchCmd, mcpInfoCmd)
}

func NewMcpCmd() *cobra.Command {
	return McpCmd
}
