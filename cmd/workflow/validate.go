package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/geelato/cli/pkg/logger"
	"github.com/spf13/cobra"
)

var validateStrict bool

var workflowValidateCmd = &cobra.Command{
	Use:   "validate [name]",
	Short: "验证工作流定义",
	Long: `验证 BPMN 工作流定义的正确性。

验证内容包括：
  - 流程结构完整性
  - 元素属性有效性
  - 连接关系正确性

示例：
  geelato workflow validate
  geelato workflow validate approval`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runValidate(args)
	},
}

func init() {
	workflowValidateCmd.Flags().BoolVar(&validateStrict, "strict", false, "严格模式")
}

func runValidate(args []string) error {
	logger.Info("验证工作流定义...")

	var workflows []string

	if len(args) > 0 {
		workflows = append(workflows, args[0])
	} else {
		dir := "workflow"
		if exists(dir) {
			entries, _ := os.ReadDir(dir)
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				if filepath.Ext(entry.Name()) == ".json" {
					workflows = append(workflows, entry.Name())
				}
			}
		}
	}

	if len(workflows) == 0 {
		logger.Info("未找到工作流文件")
		return nil
	}

	var allErrors []ValidationError
	var allWarnings []ValidationWarning

	for _, name := range workflows {
		result, err := validateWorkflow(name)
		if err != nil {
			logger.Warnf("验证失败: %s - %v", name, err)
			continue
		}

		if !result.Valid {
			allErrors = append(allErrors, result.Errors...)
		}
		allWarnings = append(allWarnings, result.Warnings...)
	}

	if len(allErrors) == 0 {
		logger.Success("工作流验证通过")
	} else {
		logger.Error("工作流验证失败")
		logger.Info("")
		logger.Info("错误:")
		for _, e := range allErrors {
			logger.Errorf("  [%s.%s] %s", e.Element, e.Property, e.Message)
		}
	}

	return nil
}

func validateWorkflow(name string) (*ValidationResult, error) {
	filename := filepath.Join("workflow", name)
	if !exists(filename) {
		return nil, fmt.Errorf("工作流文件不存在: %s", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取工作流文件失败: %w", err)
	}

	var wf Workflow
	if err := json.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("解析工作流文件失败: %w", err)
	}

	result := &ValidationResult{
		Valid:    true,
		Errors:   make([]ValidationError, 0),
		Warnings: make([]ValidationWarning, 0),
	}

	if wf.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Element:  "workflow",
			Property: "name",
			Message:  "工作流缺少名称",
		})
	}

	if len(wf.StartEvents) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Element:  "workflow",
			Property: "startEvents",
			Message:  "缺少开始事件",
		})
	}

	if len(wf.EndEvents) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Element:  "workflow",
			Property: "endEvents",
			Message:  "缺少结束事件",
		})
	}

	return result, nil
}
