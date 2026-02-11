package cmd

import (
	"errors"
	"os"
	"runtime/debug"

	"github.com/geelato/cli/internal/config"
	gerrors "github.com/geelato/cli/pkg/errors"
	"github.com/geelato/cli/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	verbose  bool
	jsonLogs bool
	version  = "dev"
	commit   = "unknown"
	date     = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "geelato",
	Short: "Geelato CLI - 低代码平台开发工具",
	Long: `Geelato CLI 是一个用于低代码平台开发的命令行工具。

== 本地创建(针对本地内容的创建) ==
  geelato api         - API管理
  geelato model       - 模型管理
  geelato page        - 页面管理
  geelato workflow    - 工作流管理

== 应用操作(针对整个应用的操作) ==
  geelato init        - 初始化应用
  geelato clone       - 从服务器克隆应用
  geelato pull        - 从云端拉取最新应用
  geelato push        - 推送变更到云端
  geelato diff        - 显示本地与云端的差异
  geelato validate    - 验证应用配置
  geelato config     - 配置管理
  geelato mcp        - MCP平台能力管理

使用 "geelato [command] --help" 查看命令帮助。`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			logger.SetLevel(logger.DebugLevel)
		}

		if jsonLogs {
			logger.SetFormatter(logger.NewJSONFormatter("2006-01-02T15:04:05Z07:00", false))
		}

		cfg, err := config.Load(cfgFile)
		if err != nil {
			logger.Warn("加载配置失败: %v", err)
		}
		config.SetGlobal(cfg)

		logger.Debugf("命令: %s %v", cmd.Name(), args)
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("程序发生 panic: %v\n%s", r, string(debug.Stack()))
			os.Exit(1)
		}
	}()

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "显示详细日志")
	rootCmd.PersistentFlags().BoolVar(&jsonLogs, "json", false, "使用 JSON 格式输出日志")
	rootCmd.Flags().Bool("version", false, "显示版本信息")

	rootCmd.AddCommand(
		apiCmd,
		cloneCmd,
		configCmd,
		diffCmd,
		initCmd,
		mcpCmd,
		modelCmd,
		pageCmd,
		pullCmd,
		pushCmd,
		validateCmd,
		workflowCmd,
	)

	rootCmd.SetVersionTemplate("{{.Name}} version {{.Version}}\n")
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if ver, _ := os.LookupEnv("VERSION"); ver != "" {
		version = ver
	}

	rootCmd.Version = version

	if err := rootCmd.Execute(); err != nil {
		var ge *gerrors.GeelatoError
		if errors.As(err, &ge) {
			logger.Errorf("错误 [%d]: %s", ge.Code, ge.Message)
			if ge.Details != "" {
				logger.Errorf("详情: %s", ge.Details)
			}
			if ge.Err != nil && verbose {
				logger.Debugf("详细错误: %v", ge.Err)
			}
			return err
		}

		logger.Errorf("执行失败: %v", err)
		return err
	}

	return nil
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("geelato")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.geelato")
		viper.AddConfigPath("$HOME/.config/geelato")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Warn("读取配置文件失败: %v", err)
		}
	}
}
