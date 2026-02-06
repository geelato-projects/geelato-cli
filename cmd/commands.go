package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/geelato/cli/cmd/config"
	"github.com/geelato/cli/cmd/workflow"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
	"github.com/spf13/cobra"
)

var (
	configCmd      *cobra.Command
	initCmd        *cobra.Command
	modelCmd       *cobra.Command
	apiCmd         *cobra.Command
	workflowCmd    *cobra.Command
	mcpCmd         *cobra.Command
	gitCmd         *cobra.Command
	syncCmd        *cobra.Command
	cloneCmd       *cobra.Command
	validateCmd    *cobra.Command
	pushCmd        *cobra.Command
	pullCmd        *cobra.Command
	diffCmd        *cobra.Command
	watchCmd       *cobra.Command
	cloneTargetDir string
	cloneAPIURL    string
	cloneForce     bool
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
	cloneCmd = NewCloneCmd()
	validateCmd = NewValidateCmd()
	pushCmd = NewPushCmd()
	pullCmd = NewPullCmd()
	diffCmd = NewDiffCmd()
	watchCmd = NewWatchCmd()
	cloneCmd.Flags().StringVar(&cloneTargetDir, "dir", "", "Target directory for clone")
	cloneCmd.Flags().StringVar(&cloneAPIURL, "api-url", "", "API server URL")
	cloneCmd.Flags().BoolVarP(&cloneForce, "force", "f", false, "Overwrite existing directory")
}

func NewAppCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "app",
		Short: "Application management",
		Long:  `Application management commands (deprecated, use direct commands instead)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Warn("The 'app' command is deprecated. Use direct commands like 'validate', 'push', 'pull' instead.")
			return cmd.Help()
		},
	}
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

func NewCloneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone <repository-url>",
		Short: "克隆远程应用",
		Long: `从远程仓库克隆应用到本地

URL 格式: http://cli.geelato.cn/{tenant}/{app-code}

示例:
  geelato clone http://cli.geelato.cn/tenant/myapp
  geelato clone http://cli.geelato.cn/my-tenant/order-system`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}
			return runClone(args[0])
		},
	}
}

func runClone(repoURL string) error {
	logger.Infof("Cloning from: %s", repoURL)

	tenant, appCode, err := parseCloneURL(repoURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	logger.Infof("Tenant: %s", tenant)
	logger.Infof("App Code: %s", appCode)

	targetDir := cloneTargetDir
	if targetDir == "" {
		targetDir = appCode
	}

	if _, err := os.Stat(targetDir); err == nil && !cloneForce {
		confirm, err := prompt.Confirm(fmt.Sprintf("Directory '%s' already exists. Overwrite?", targetDir), false)
		if err != nil {
			return err
		}
		if !confirm {
			logger.Info("Clone cancelled")
			return nil
		}
		os.RemoveAll(targetDir)
	}

	apiURL := cloneAPIURL
	if apiURL == "" {
		apiURL = extractAPIURL(repoURL)
	}

	logger.Infof("API URL: %s", apiURL)

	data := map[string]string{
		"tenant": tenant,
		"code":   appCode,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	downloadURL := apiURL + "/api/cli/app/download"

	logger.Infof("Downloading from: %s", downloadURL)

	resp, err := http.Post(downloadURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	tmpFile := filepath.Join(os.TempDir(), "geelato-clone-"+appCode+".zip")
	tmpFile, err = filepath.Abs(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to get temp file path: %w", err)
	}

	outFile, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save package: %w", err)
	}

	logger.Infof("Package downloaded to: %s", tmpFile)

	if err := extractAndInitialize(tmpFile, targetDir); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to initialize application: %w", err)
	}

	os.Remove(tmpFile)

	logger.Info("Application cloned successfully!")
	logger.Infof("Target directory: %s", targetDir)
	logger.Info("")
	logger.Info("Next steps:")
	logger.Info("  cd %s", targetDir)
	logger.Info("  geelato model list        # List available models")
	logger.Info("  geelato api list         # List available APIs")

	return nil
}

func parseCloneURL(repoURL string) (tenant, appCode string, err error) {
	repoURL = strings.TrimSpace(repoURL)

	if !strings.HasPrefix(repoURL, "http://") && !strings.HasPrefix(repoURL, "https://") {
		repoURL = "http://" + repoURL
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse URL: %w", err)
	}

	path := strings.TrimPrefix(u.Path, "/")
	path = strings.TrimSuffix(path, "/")

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("URL path should contain tenant and app code (e.g., /tenant/app-code)")
	}

	tenant = parts[0]
	appCode = parts[1]

	if tenant == "" || appCode == "" {
		return "", "", fmt.Errorf("tenant and app code cannot be empty")
	}

	return tenant, appCode, nil
}

func extractAPIURL(repoURL string) string {
	u, err := url.Parse(repoURL)
	if err != nil {
		return ""
	}

	baseURL := u.Scheme + "://" + u.Host
	if u.Port() != "" {
		baseURL += ":" + u.Port()
	}

	return baseURL
}

func extractAndInitialize(zipPath, targetDir string) error {
	logger.Info("Extracting package...")

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer reader.Close()

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	for _, file := range reader.File {
		path := filepath.Join(targetDir, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open zip file: %w", err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return fmt.Errorf("failed to extract file: %w", err)
		}
	}

	logger.Infof("Package extracted to: %s", targetDir)
	return nil
}
