package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/internal/app"
	"github.com/geelato/cli/pkg/logger"
)

var appInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "查看当前应用信息",
	Long: `查看当前 Geelato 应用的基本信息和配置。

显示内容包括：
  - 应用名称、描述、版本
  - 创建时间、最后修改时间
  - 配置信息

示例：
  geelato app info`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInfo()
	},
}

func runInfo() error {
	geelatoFile := "geelato.json"
	if _, err := os.Stat(geelatoFile); os.IsNotExist(err) {
		return fmt.Errorf("当前目录不是有效的 Geelato 应用")
	}

	data, err := os.ReadFile(geelatoFile)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config struct {
		Meta struct {
			Version    string `json:"version"`
			AppID     string `json:"appId"`
			Name      string `json:"name"`
			Description string `json:"description"`
			CreatedAt string `json:"createdAt"`
		} `json:"meta"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	logger.Info("")
	logger.Info("应用信息")
	logger.Info("=========")
	logger.Info("")
	logger.Infof("名称:      %s", config.Meta.Name)
	logger.Infof("应用ID:   %s", config.Meta.AppID)
	logger.Infof("描述:     %s", config.Meta.Description)
	logger.Infof("版本:     %s", config.Meta.Version)
	logger.Infof("创建时间: %s", config.Meta.CreatedAt)
	logger.Info("")

	cwd, _ := os.Getwd()
	logger.Infof("路径:     %s", cwd)

	return nil
}
