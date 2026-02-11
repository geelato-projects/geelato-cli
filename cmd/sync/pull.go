package sync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geelato/cli/internal/app"
	"github.com/geelato/cli/pkg/logger"
	"github.com/spf13/cobra"
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

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	geelatoPath := filepath.Join(cwd, "geelato.json")
	if _, err := os.Stat(geelatoPath); os.IsNotExist(err) {
		return fmt.Errorf("not a Geelato application: geelato.json not found in %s", cwd)
	}

	// 从 geelato.json 读取 repo 配置
	appConfig, err := app.LoadAppConfig(cwd)
	if err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}

	repoURL := app.GetRepoFromConfig(appConfig)
	if repoURL == "" {
		return fmt.Errorf("repo URL not configured. Please run 'geelato config repo <url>' first")
	}

	_, _, apiURL, err := app.ParseRepoURL(repoURL)
	if err != nil {
		return fmt.Errorf("failed to parse repo URL: %w", err)
	}

	logger.Infof("使用 API URL: %s", apiURL)

	manager := NewManagerWithAPI(apiURL)

	status, err := manager.GetStatus()
	if err != nil {
		return fmt.Errorf("获取状态失败: %w", err)
	}

	logger.Infof("当前版本: %s", status.LocalVersion)
	logger.Infof("最新版本: %s", status.RemoteVersion)

	logger.Info("拉取功能正在开发中...")

	return nil
}
