package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有配置",
	Long:  `列出当前配置的所有选项及其值`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := viper.ConfigFileUsed()
		if cfg == "" {
			cfg = "未使用配置文件"
		}

		fmt.Println()
		fmt.Printf("配置文件: %s\n", cfg)
		fmt.Println()
		fmt.Println("当前配置:")
		fmt.Println("--------")

		settings := viper.AllSettings()
		for key, value := range settings {
			fmt.Printf("  %s: %v\n", key, value)
		}

		return nil
	},
}
