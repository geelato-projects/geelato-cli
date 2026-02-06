package sync

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
)

var (
	statusJSON    bool
	statusVerbose bool
)

var syncStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看同步状态",
	Long: `查看当前应用的同步状态。

显示内容包括：
  - 本地版本和云端版本
  - 待推送的本地变更
  - 待拉取的云端变更
  - 同步冲突列表

示例：
  geelato sync status
  geelato sync status --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStatus()
	},
}

func init() {
	syncStatusCmd.Flags().BoolVarP(&statusJSON, "json", "j", false, "JSON 格式输出")
	syncStatusCmd.Flags().BoolVarP(&statusVerbose, "verbose", "v", false, "详细输出")
}

func runStatus() error {
	manager := NewManager()

	status, err := manager.GetStatus()
	if err != nil {
		return fmt.Errorf("获取同步状态失败: %w", err)
	}

	if statusJSON {
		return printJSON(status)
	}

	return printStatus(status)
}

func printStatus(status *SyncStatusView) {
	logger.Info("")
	logger.Info("同步状态")
	logger.Info("========")
	logger.Info("")

	logger.Infof("本地版本:  %s", status.LocalVersion)
	logger.Infof("云端版本:  %s", status.RemoteVersion)
	logger.Info("")

	if status.AheadBy > 0 {
		logger.Warnf("待推送: %d 个变更", status.AheadBy)
	} else {
		logger.Success("已同步，无待推送变更")
	}

	if status.BehindBy > 0 {
		logger.Warnf("待拉取: %d 个变更", status.BehindBy)
	} else {
		logger.Success("已是最新版本")
	}

	if status.HasConflict {
		logger.Error("存在同步冲突，请使用 'geelato sync resolve' 解决")
	}

	if status.LastSyncTime != "" {
		logger.Infof("最后同步: %s", status.LastSyncTime)
	}
}

func printJSON(status *SyncStatusView) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"localVersion\": \"%s\",\n", status.LocalVersion)
	fmt.Printf("  \"remoteVersion\": \"%s\",\n", status.RemoteVersion)
	fmt.Printf("  \"aheadBy\": %d,\n", status.AheadBy)
	fmt.Printf("  \"behindBy\": %d,\n", status.BehindBy)
	fmt.Printf("  \"hasConflict\": %v,\n", status.HasConflict)
	fmt.Printf("  \"lastSyncTime\": \"%s\"\n", status.LastSyncTime)
	fmt.Printf("}\n")
	return nil
}
