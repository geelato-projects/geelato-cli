package sync

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
)

var (
	resolveStrategy string
	resolveAll     bool
)

var syncResolveCmd = &cobra.Command{
	Use:   "resolve [file]",
	Short: "解决同步冲突",
	Long: `解决同步过程中检测到的冲突。

冲突解决策略：
  - ours: 保留本地版本
  - theirs: 保留云端版本
  - manual: 手动合并

示例：
  geelato sync resolve
  geelato sync resolve --strategy theirs
  geelato sync resolve models/user.yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runResolve(args)
	},
}

func init() {
	syncResolveCmd.Flags().StringVar(&resolveStrategy, "strategy", "manual", "解决策略")
	syncResolveCmd.Flags().BoolVar(&resolveAll, "all", false, "解决所有冲突")
}

func runResolve(args []string) error {
	manager := NewManager()

	conflicts, err := manager.DetectConflicts()
	if err != nil {
		return fmt.Errorf("检测冲突失败: %w", err)
	}

	if len(conflicts) == 0 {
		logger.Success("没有检测到冲突")
		return nil
	}

	logger.Warnf("检测到 %d 个冲突", len(conflicts))

	displayConflicts(conflicts)

	logger.Info("")
	logger.Info("冲突解决功能正在开发中...")
	logger.Info("请手动编辑冲突文件，然后使用 'geelato sync push' 重新推送")

	return nil
}

func displayConflicts(conflicts []Conflict) {
	logger.Info("")
	logger.Info("冲突列表:")
	logger.Info("---------")

	for i, conflict := range conflicts {
		logger.Infof("%d. %s", i+1, conflict.Path)
		logger.Infof("   本地哈希: %s", conflict.LocalHash[:8])
		logger.Infof("   云端哈希: %s", conflict.RemoteHash[:8])
	}
}
