package cmd

import (
	"github.com/geelato/cli/cmd/config"
	"github.com/geelato/cli/cmd/workflow"
	"github.com/spf13/cobra"
)

var (
	configCmd   *cobra.Command
	initCmd     *cobra.Command
	modelCmd    *cobra.Command
	apiCmd      *cobra.Command
	workflowCmd *cobra.Command
	mcpCmd      *cobra.Command
	gitCmd      *cobra.Command
	syncCmd     *cobra.Command
	validateCmd *cobra.Command
	pushCmd     *cobra.Command
	pullCmd     *cobra.Command
	diffCmd     *cobra.Command
	watchCmd    *cobra.Command
	pageCmd     *cobra.Command
)

func init() {
	configCmd = config.NewConfigCmd()
	initCmd = initCmdFn()
	modelCmd = NewModelCmd()
	apiCmd = NewApiCmd()
	workflowCmd = workflow.NewWorkflowCmd()
	mcpCmd = NewMcpCmd()
	gitCmd = NewGitCmd()
	syncCmd = NewSyncCmd()
	validateCmd = NewValidateCmd()
	pushCmd = NewPushCmd()
	pullCmd = NewPullCmd()
	diffCmd = NewDiffCmd()
	watchCmd = NewWatchCmd()
	pageCmd = NewPageCmd()
	rootCmd.AddCommand(cloneCmd)
}

func NewMcpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "MCP管理",
		Long:  `管理 MCP 平台能力同步`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

func NewGitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "git",
		Short: "Git集成",
		Long:  `Git 集成命令，用于版本控制`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

func NewSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "云端同步",
		Long:  `云端同步命令，用于推送和拉取变更`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}
