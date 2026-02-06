package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
)

var workflowListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有工作流",
	Long: `列出当前应用下的所有工作流定义。

示例：
  geelato workflow list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList()
	},
}

func runList() error {
	logger.Info("工作流列表")
	logger.Info("==========")
	logger.Info("")

	dir := "workflow"
	if !exists(dir) {
		logger.Info("未找到工作流目录")
		logger.Info("使用 'geelato workflow create' 创建工作流")
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取工作流目录失败: %w", err)
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		name := entry.Name()
		logger.Infof("- %s", name[:len(name)-5])
		count++
	}

	logger.Info("")
	logger.Infof("共 %d 个工作流", count)

	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
