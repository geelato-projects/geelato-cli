package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/geelato/cli/cmd/initializer"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
	"github.com/spf13/cobra"
)

var (
	workflowName   string
	workflowDesc   string
	workflowFormat string
	interactive    bool
)

var workflowCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "create(创建新工作流)",
	Long: `创建一个新的 BPMN 工作流定义。

工作流定义包括：
  - 基本信息：名称、描述
  - 流程元素：开始事件、结束事件、任务
  - 连接关系：顺序流

示例：
  geelato workflow create approval
  geelato workflow create approval --desc "审批流程"`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if interactive {
			return runInteractiveCreate()
		}

		if len(args) > 0 {
			workflowName = args[0]
		}

		if workflowName == "" {
			return fmt.Errorf("工作流名称不能为空")
		}

		return runCreate()
	},
}

func init() {
	workflowCreateCmd.Flags().StringVar(&workflowDesc, "desc", "", "工作流描述")
	workflowCreateCmd.Flags().StringVar(&workflowFormat, "format", "json", "输出格式 (json)")
	workflowCreateCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "交互式模式")
}

func runInteractiveCreate() error {
	logger.Info("创建新工作流")
	logger.Info("")

	name, err := prompt.Input("工作流名称", "")
	if err != nil {
		return err
	}
	workflowName = name

	desc, err := prompt.Input("工作流描述", "")
	if err != nil {
		return err
	}
	workflowDesc = desc

	return runCreate()
}

func runCreate() error {
	logger.Infof("创建工作流: %s", workflowName)

	dir := "workflow"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建工作流目录失败: %w", err)
	}

	filename := filepath.Join(dir, workflowName+"."+workflowFormat)
	now := time.Now().Format(time.RFC3339)

	if err := initializer.CreateWorkflowFile(filename, workflowName, workflowDesc, now, now); err != nil {
		return fmt.Errorf("创建工作流文件失败: %w", err)
	}

	logger.Success("工作流创建成功")
	logger.Infof("工作流文件: %s", filename)

	return nil
}
