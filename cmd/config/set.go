package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "设置配置值",
	Long:  `设置指定配置项的值`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		viper.Set(key, value)

		cfgFile := viper.ConfigFileUsed()
		if cfgFile == "" {
			cfgFile = "$HOME/.geelato/geelato.yaml"
		}

		if err := viper.WriteConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Printf("警告: 配置文件不存在，变更仅在当前会话有效\n")
				fmt.Printf("如需保存，请使用 'geelato config edit' 创建配置文件\n")
			} else {
				return fmt.Errorf("保存配置失败: %w", err)
			}
		} else {
			fmt.Printf("配置已更新: %s = %s\n", key, value)
		}

		return nil
	},
}
