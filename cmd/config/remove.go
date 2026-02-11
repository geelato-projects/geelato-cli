package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configRemoveCmd = &cobra.Command{
	Use:   "remove <key>",
	Short: "remove(删除配置项)",
	Long:  `删除指定配置项`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		if !viper.IsSet(key) {
			fmt.Printf("配置项 '%s' 不存在\n", key)
			return nil
		}

		viper.Set(key, nil)

		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("保存配置失败: %w", err)
		}

		fmt.Printf("配置已删除: %s\n", key)
		return nil
	},
}
