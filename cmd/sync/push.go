package sync

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
)

var (
	pushMessage string
	pushAll    bool
	pushDryRun bool
)

var syncPushCmd = &cobra.Command{
	Use:   "push",
	Short: "推送变更到云端",
	Long: `推送本地变更到云端平台。

推送内容包括：
  - 新增的模型定义
  - 修改的 API 接口
  - 新增或修改的工作流
  - 配置文件的变更

示例：
  geelato sync push "更新用户模型"
  geelato sync push --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			pushMessage = args[0]
		}

		if pushMessage == "" {
			pushMessage = generateDefaultMessage()
		}

		return runPush()
	},
}

func init() {
	syncPushCmd.Flags().StringVar(&pushMessage, "message", "", "提交消息")
	syncPushCmd.Flags().BoolVar(&pushAll, "all", false, "推送所有变更")
	syncPushCmd.Flags().BoolVar(&pushDryRun, "dry-run", false, "演练模式（不实际推送）")
}

func runPush() error {
	logger.Infof("开始推送变更: %s", pushMessage)

	manager := NewManager()

	changes, err := manager.DetectChanges(".")
	if err != nil {
		return fmt.Errorf("检测变更失败: %w", err)
	}

	if len(changes) == 0 {
		logger.Info("没有需要推送的变更")
		return nil
	}

	logger.Infof("发现 %d 个变更", len(changes))

	displayChanges(changes)

	if pushDryRun {
		logger.Info("演练模式完成，未实际推送")
		return nil
	}

	if !pushAll && len(changes) > 0 {
		confirm, err := prompt.Confirm(fmt.Sprintf("确认推送 %d 个变更？", len(changes)), true)
		if err != nil {
			return err
		}
		if !confirm {
			logger.Info("已取消推送")
			return nil
		}
	}

	result, err := manager.Push(changes, pushMessage)
	if err != nil {
		return fmt.Errorf("推送失败: %w", err)
	}

	logger.Success("推送完成")
	logger.Infof("版本: %s", result.Version)
	logger.Infof("变更数量: %d", len(changes))

	return nil
}

func displayChanges(changes []Change) {
	logger.Info("")
	logger.Info("变更列表:")
	logger.Info("---------")

	for _, change := range changes {
		var icon string
		switch change.Type {
		case string(ChangeTypeAdded):
			icon = "[新增]"
		case string(ChangeTypeModified):
			icon = "[修改]"
		case string(ChangeTypeDeleted):
			icon = "[删除]"
		default:
			icon = "[未知]"
		}
		logger.Infof("  %s %s (%s)", icon, change.Path, change.Description)
	}
}

func generateDefaultMessage() string {
	return fmt.Sprintf("同步更新 - %s", time.Now().Format("2006-01-02 15:04:05"))
}
