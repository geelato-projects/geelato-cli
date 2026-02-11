package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "get(获取配置值)",
	Long:  `获取指定配置项的值`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := viper.Get(key)

		if value == nil {
			fmt.Printf("配置项 '%s' 不存在\n", key)
			return nil
		}

		fmt.Printf("%v\n", value)
		return nil
	},
}
