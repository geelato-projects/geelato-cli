package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit(编辑配置文件)",
	Long:  `使用默认编辑器打开配置文件`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgFile := viper.ConfigFileUsed()

		if cfgFile == "" {
			homeDir, _ := os.UserHomeDir()
			cfgDir := filepath.Join(homeDir, ".geelato")
			cfgFile = filepath.Join(cfgDir, "geelato.yaml")

			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				fmt.Printf("配置文件不存在，将创建: %s\n", cfgFile)

				if err := os.MkdirAll(cfgDir, 0755); err != nil {
					return fmt.Errorf("创建配置目录失败: %w", err)
				}

				viper.SetConfigFile(cfgFile)
				if err := viper.SafeWriteConfig(); err != nil {
					return fmt.Errorf("创建配置文件失败: %w", err)
				}
			}
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		editCmd := exec.Command(editor, cfgFile)
		editCmd.Stdin = os.Stdin
		editCmd.Stdout = os.Stdout
		editCmd.Stderr = os.Stderr

		return editCmd.Run()
	},
}
