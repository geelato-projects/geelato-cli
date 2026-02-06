package mcp

import (
	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
)

var (
	listJSON     bool
	listCategory string
)

var mcpListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出可用能力",
	Long: `列出本地和云端平台可用的 MCP 能力。

能力分类包括：
  - 认证能力 (auth)
  - 数据库能力 (database)
  - 消息队列能力 (mq)
  - 存储能力 (storage)
  - AI 能力 (ai)

示例：
  geelato mcp list
  geelato mcp list --category database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList()
	},
}

func init() {
	mcpListCmd.Flags().BoolVar(&listJSON, "json", false, "JSON 格式输出")
	mcpListCmd.Flags().StringVar(&listCategory, "category", "", "按分类筛选")
}

func runList() error {
	logger.Info("")
	logger.Info("可用 MCP 能力")
	logger.Info("===============")
	logger.Info("")

	logger.Info("MCP 功能正在开发中...")
	logger.Info("")

	capabilities := []Capability{
		{
			ID:          "database",
			Name:        "数据库能力",
			Version:     "1.0.0",
			Category:    "database",
			Description: "提供数据库查询、更新、事务等能力",
		},
		{
			ID:          "auth",
			Name:        "认证能力",
			Version:     "1.0.0",
			Category:    "auth",
			Description: "提供用户认证、授权、令牌管理等能力",
		},
	}

	for _, cap := range capabilities {
		if listCategory != "" && cap.Category != listCategory {
			continue
		}
		logger.Infof("- %s (%s)", cap.Name, cap.Category)
		logger.Infof("  版本: %s", cap.Version)
		logger.Infof("  描述: %s", cap.Description)
		logger.Info("")
	}

	logger.Infof("共 %d 个能力", len(capabilities))

	return nil
}

type Capability struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Installed   bool   `json:"installed"`
}
