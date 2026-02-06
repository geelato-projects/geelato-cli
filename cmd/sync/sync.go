package sync

import (
	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
)

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "云端同步",
	Long: `云端同步命令，用于将本地变更同步到云端平台或从云端拉取变更。

同步功能包括：
  geelato sync push     推送本地变更到云端
  geelato sync pull     从云端拉取变更到本地
  geelato sync status   查看同步状态
  geelato sync resolve  解决同步冲突

使用 "geelato sync [command] --help" 查看命令帮助。`,
}

func init() {
	SyncCmd.AddCommand(syncPushCmd, syncPullCmd, syncStatusCmd, syncResolveCmd)
}

func NewSyncCmd() *cobra.Command {
	return SyncCmd
}
