package sync

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
)

var (
	pullVersion   string
	pullAutoMerge bool
)

var syncPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "从云端拉取变更",
	Long: `从云端平台拉取最新的变更到本地。

拉取内容包括：
  - 模型定义的更新
  - API 接口的变更
  - 工作流的调整
  - 配置文件的更新

示例：
  geelato sync pull
  geelato sync pull --version latest`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			pullVersion = args[0]
		}

		return runPull()
	},
}

func init() {
	syncPullCmd.Flags().StringVar(&pullVersion, "version", "latest", "指定版本")
	syncPullCmd.Flags().BoolVar(&pullAutoMerge, "auto-merge", false, "自动合并冲突")
}

func runPull() error {
	logger.Info("从云端拉取变更...")

	manager := NewManager()

	status, err := manager.GetStatus()
	if err != nil {
		return fmt.Errorf("获取状态失败: %w", err)
	}

	logger.Infof("当前版本: %s", status.LocalVersion)
	logger.Infof("最新版本: %s", status.RemoteVersion)

	logger.Info("拉取功能正在开发中...")

	return nil
}
