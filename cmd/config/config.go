package config

import (
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "config(配置管理)",
	Long:  `管理 Geelato CLI 工具的配置信息`,
}

func init() {
	ConfigCmd.AddCommand(configListCmd, configGetCmd, configSetCmd, configRemoveCmd, configEditCmd)
}

func NewConfigCmd() *cobra.Command {
	return ConfigCmd
}
