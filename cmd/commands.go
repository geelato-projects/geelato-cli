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
	validateCmd *cobra.Command
	pushCmd     *cobra.Command
	pullCmd     *cobra.Command
	diffCmd     *cobra.Command
	pageCmd     *cobra.Command
)

func init() {
	configCmd = config.NewConfigCmd()
	initCmd = initCmdFn()
	modelCmd = NewModelCmd()
	apiCmd = NewApiCmd()
	workflowCmd = workflow.NewWorkflowCmd()
	mcpCmd = NewMcpCmd()
	validateCmd = NewValidateCmd()
	pushCmd = NewPushCmd()
	pullCmd = NewPullCmd()
	diffCmd = NewDiffCmd()
	pageCmd = NewPageCmd()
}

func NewMcpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "mcp(MCP管理)",
		Long:  `管理 MCP 平台能力同步`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}
