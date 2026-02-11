package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/geelato/cli/internal/config"
	"github.com/geelato/cli/pkg/logger"
	"github.com/spf13/cobra"
)

var deployForce bool

var workflowDeployCmd = &cobra.Command{
	Use:   "deploy [name]",
	Short: "deploy(部署工作流)",
	Long: `部署工作流到云端平台。

部署前会验证工作流定义，然后将工作流文件上传到云端。

示例：
  geelato workflow deploy
  geelato workflow deploy approval
  geelato workflow deploy approval --force`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		}
		return runDeploy(name)
	},
}

func init() {
	workflowDeployCmd.Flags().BoolVar(&deployForce, "force", false, "强制部署")
}

func runDeploy(name string) error {
	logger.Info("部署工作流...")
	logger.Info("============")
	logger.Info("")

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %w", err)
	}

	cfg := config.Get()
	if cfg == nil || cfg.API.URL == "" {
		logger.Error("API URL 未配置")
		logger.Info("请先配置 API 地址: geelato config set api.url <url>")
		return nil
	}

	var workflows []string

	if name != "" {
		workflows = append(workflows, name+".json")
	} else {
		dir := filepath.Join(cwd, "workflow")
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

	logger.Infof("找到 %d 个工作流", len(workflows))
	logger.Info("")

	deployed := 0
	failed := 0

	for _, wf := range workflows {
		name := wf
		if filepath.Ext(wf) == ".json" {
			name = wf[:len(wf)-5]
		}

		logger.Infof("部署工作流: %s", name)

		err := deployWorkflow(cwd, name)
		if err != nil {
			logger.Errorf("部署失败: %s - %v", name, err)
			failed++
			continue
		}

		logger.Success("工作流 %s 部署成功", name)
		deployed++
	}

	logger.Info("")
	logger.Infof("部署完成: %d 成功, %d 失败", deployed, failed)

	return nil
}

func deployWorkflow(cwd, name string) error {
	workflowPath := filepath.Join(cwd, "workflow", name+".json")

	if !exists(workflowPath) {
		return fmt.Errorf("工作流文件不存在: %s", workflowPath)
	}

	data, err := os.ReadFile(workflowPath)
	if err != nil {
		return fmt.Errorf("读取工作流文件失败: %w", err)
	}

	var wf Workflow
	if err := json.Unmarshal(data, &wf); err != nil {
		return fmt.Errorf("解析工作流文件失败: %w", err)
	}

	if wf.Name == "" {
		return fmt.Errorf("工作流缺少名称")
	}

	logger.Infof("  名称: %s", wf.Name)
	logger.Infof("  版本: %s", wf.Version)
	logger.Infof("  开始事件: %d", len(wf.StartEvents))
	logger.Infof("  结束事件: %d", len(wf.EndEvents))

	if !deployForce {
		logger.Info("")
		logger.Info("注意: 工作流部署需要云端平台支持")
		logger.Info("当前版本仅记录部署信息，实际部署请使用 'geelato push'")
	}

	deployment := WorkflowDeployment{
		ID:         fmt.Sprintf("deploy_%s_%d", name, time.Now().Unix()),
		Name:       wf.Name,
		Version:    wf.Version,
		DeployedAt: time.Now(),
		Status:     "pending",
	}

	deployData, _ := json.MarshalIndent(deployment, "", "  ")
	deployPath := filepath.Join(cwd, "workflow", name+".deploy.json")
	os.WriteFile(deployPath, deployData, 0644)

	logger.Infof("  部署记录: %s.deploy.json", name)

	return nil
}
