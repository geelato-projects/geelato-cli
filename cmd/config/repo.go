package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo [url]",
	Short: "管理仓库地址",
	Long: `管理 Geelato 应用的仓库地址（repo）。

如果不提供参数，则显示当前仓库地址。
如果提供 URL 参数，则设置仓库地址。

示例:
  geelato config repo                           # 查看当前仓库地址
  geelato config repo http://localhost:8080/tenant/app  # 设置仓库地址`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		geelatoPath := filepath.Join(cwd, "geelato.json")
		if _, err := os.Stat(geelatoPath); os.IsNotExist(err) {
			return fmt.Errorf("not a Geelato application: geelato.json not found in %s", cwd)
		}

		content, err := os.ReadFile(geelatoPath)
		if err != nil {
			return fmt.Errorf("failed to read geelato.json: %w", err)
		}

		var config map[string]interface{}
		if err := json.Unmarshal(content, &config); err != nil {
			return fmt.Errorf("failed to parse geelato.json: %w", err)
		}

		if len(args) == 0 {
			// 查看当前 repo
			repoURL := getRepoURL(config)
			if repoURL == "" {
				fmt.Println("当前未设置仓库地址 (repo)")
				fmt.Println("使用 'geelato config repo <url>' 设置仓库地址")
			} else {
				fmt.Printf("当前仓库地址: %s\n", repoURL)
			}
			return nil
		}

		// 设置 repo
		newRepo := args[0]
		setRepoURL(config, newRepo)

		newContent, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		if err := os.WriteFile(geelatoPath, newContent, 0644); err != nil {
			return fmt.Errorf("failed to write geelato.json: %w", err)
		}

		fmt.Printf("仓库地址已更新: %s\n", newRepo)
		return nil
	},
}

func getRepoURL(config map[string]interface{}) string {
	// 新结构：config.repo.url
	if configObj, ok := config["config"].(map[string]interface{}); ok {
		if repo, ok := configObj["repo"].(map[string]interface{}); ok {
			if url, ok := repo["url"].(string); ok {
				return url
			}
		}
	}
	// 兼容旧结构：repo（顶层）
	if repo, ok := config["repo"].(string); ok {
		return repo
	}
	return ""
}

func setRepoURL(config map[string]interface{}, url string) {
	// 确保 config 对象存在
	if _, ok := config["config"].(map[string]interface{}); !ok {
		config["config"] = map[string]interface{}{}
	}

	configObj := config["config"].(map[string]interface{})

	// 确保 repo 对象存在
	if _, ok := configObj["repo"].(map[string]interface{}); !ok {
		configObj["repo"] = map[string]interface{}{}
	}

	repo := configObj["repo"].(map[string]interface{})
	repo["url"] = url
}

func init() {
	ConfigCmd.AddCommand(repoCmd)
}
