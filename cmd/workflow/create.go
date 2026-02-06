package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
)

var (
	workflowName     string
	workflowDesc    string
	workflowFormat   string
	workflowTemplate string
	interactive     bool
)

var workflowCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "创建新工作流",
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
	workflowCreateCmd.Flags().StringVar(&workflowTemplate, "template", "", "使用模板创建")
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

	wf := createDefaultWorkflow()

	filename := filepath.Join(dir, workflowName+"."+workflowFormat)

	content := wf.ToJSON()

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("创建工作流文件失败: %w", err)
	}

	logger.Success("工作流创建成功")
	logger.Infof("工作流文件: %s", filename)

	return nil
}

func createDefaultWorkflow() *Workflow {
	return &Workflow{
		Name:        workflowName,
		Description: workflowDesc,
		Version:     "1.0",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		StartEvents: []StartEvent{
			{
				ID:   "start_1",
				Name: "开始",
			},
		},
		EndEvents: []EndEvent{
			{
				ID:   "end_1",
				Name: "结束",
			},
		},
		Tasks: []Task{
			{
				ID:   "task_1",
				Name: "处理任务",
				Type: "serviceTask",
			},
		},
		SequenceFlows: []SequenceFlow{
			{
				ID:        "flow_1",
				SourceRef: "start_1",
				TargetRef: "task_1",
			},
			{
				ID:        "flow_2",
				SourceRef: "task_1",
				TargetRef: "end_1",
			},
		},
	}
}

type Workflow struct {
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Version       string       `json:"version"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
	StartEvents   []StartEvent `json:"startEvents"`
	EndEvents     []EndEvent   `json:"endEvents"`
	Tasks         []Task       `json:"tasks"`
	SequenceFlows []SequenceFlow `json:"sequenceFlows"`
}

type StartEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EndEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type SequenceFlow struct {
	ID        string `json:"id"`
	SourceRef string `json:"sourceRef"`
	TargetRef string `json:"targetRef"`
}

func (w *Workflow) ToJSON() string {
	data := fmt.Sprintf(`{
  "meta": {
    "name": "%s",
    "description": "%s",
    "version": "%s",
    "createdAt": "%s",
    "updatedAt": "%s"
  },
  "startEvents": [
    {
      "id": "%s",
      "name": "%s"
    }
  ],
  "endEvents": [
    {
      "id": "%s",
      "name": "%s"
    }
  ],
  "tasks": [
    {
      "id": "%s",
      "name": "%s",
      "type": "%s"
    }
  ],
  "sequenceFlows": [
    {
      "id": "%s",
      "sourceRef": "%s",
      "targetRef": "%s"
    },
    {
      "id": "%s",
      "sourceRef": "%s",
      "targetRef": "%s"
    }
  ]
}`,
		w.Name, w.Description, w.Version, w.CreatedAt.Format(time.RFC3339), w.UpdatedAt.Format(time.RFC3339),
		w.StartEvents[0].ID, w.StartEvents[0].Name,
		w.EndEvents[0].ID, w.EndEvents[0].Name,
		w.Tasks[0].ID, w.Tasks[0].Name, w.Tasks[0].Type,
		w.SequenceFlows[0].ID, w.SequenceFlows[0].SourceRef, w.SequenceFlows[0].TargetRef,
		w.SequenceFlows[1].ID, w.SequenceFlows[1].SourceRef, w.SequenceFlows[1].TargetRef,
	)
	return data
}
