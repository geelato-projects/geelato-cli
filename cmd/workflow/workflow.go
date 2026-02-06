package workflow

import (
	"github.com/spf13/cobra"
)

var WorkflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "工作流管理",
	Long: `工作流管理命令，用于创建和管理 BPMN 流程。

工作流功能包括：
  geelato workflow create    创建新工作流
  geelato workflow list     列出所有工作流
  geelato workflow validate  验证工作流定义
  geelato workflow deploy    部署工作流

使用 "geelato workflow [command] --help" 查看命令帮助。`,
}

func init() {
	WorkflowCmd.AddCommand(workflowCreateCmd, workflowListCmd, workflowValidateCmd, workflowDeployCmd)
}

func NewWorkflowCmd() *cobra.Command {
	return WorkflowCmd
}
