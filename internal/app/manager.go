package app

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/geelato/cli/internal/config"
)

type Manager struct {
	appsDir string
}

func NewManager() *Manager {
	cfg := config.Get()
	var appsDir string
	if cfg != nil && cfg.Cache.Dir != "" {
		appsDir = filepath.Join(cfg.Cache.Dir, "apps")
	} else {
		appsDir = ".geelato/apps"
	}
	return &Manager{
		appsDir: appsDir,
	}
}

func (m *Manager) Exists(path string) bool {
	geelatoFile := filepath.Join(path, "geelato.json")
	_, err := os.Stat(geelatoFile)
	return err == nil
}

func (m *Manager) Load(path string) (*Application, error) {
	geelatoFile := filepath.Join(path, "geelato.json")
	if _, err := os.Stat(geelatoFile); os.IsNotExist(err) {
		return nil, nil
	}

	content, err := os.ReadFile(geelatoFile)
	if err != nil {
		return nil, err
	}

	var app Application
	if err := json.Unmarshal(content, &app); err != nil {
		return nil, err
	}

	app.Path = path
	return &app, nil
}

func (m *Manager) Save(app *Application) error {
	geelatoFile := filepath.Join(app.Path, "geelato.json")

	data, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(geelatoFile, data, 0644)
}

func (m *Manager) FindApps(searchDir string) ([]ApplicationInfo, error) {
	var apps []ApplicationInfo

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		appDir := filepath.Join(searchDir, entry.Name())
		geelatoFile := filepath.Join(appDir, "geelato.json")

		if _, err := os.Stat(geelatoFile); err == nil {
			content, _ := os.ReadFile(geelatoFile)
			info := ApplicationInfo{
				Path: appDir,
			}
			parseAppInfo(string(content), &info)
			if stat, _ := os.Stat(geelatoFile); stat != nil {
				info.LastModTime = stat.ModTime()
			}
			apps = append(apps, info)
		}
	}

	return apps, nil
}

func parseAppInfo(content string, info *ApplicationInfo) {
	// Simplified JSON parsing for basic info
	// A full implementation would use json.Unmarshal
}

// LoadAppConfig 加载 geelato.json 配置文件
func LoadAppConfig(appPath string) (map[string]interface{}, error) {
	geelatoFile := filepath.Join(appPath, "geelato.json")
	content, err := os.ReadFile(geelatoFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read geelato.json: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse geelato.json: %w", err)
	}

	return config, nil
}

// GetRepoFromConfig 从配置中获取 repo 地址
func GetRepoFromConfig(config map[string]interface{}) string {
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

// ParseRepoURL 解析 repo URL，返回 tenant, appCode, apiURL
func ParseRepoURL(repoURL string) (tenant, appCode, apiURL string, err error) {
	repoURL = strings.TrimSpace(repoURL)

	if !strings.HasPrefix(repoURL, "http://") && !strings.HasPrefix(repoURL, "https://") {
		repoURL = "http://" + repoURL
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse URL: %w", err)
	}

	path := strings.TrimPrefix(u.Path, "/")
	path = strings.TrimSuffix(path, "/")

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("URL path should contain tenant and app code (e.g., /tenant/app-code)")
	}

	tenant = parts[0]
	appCode = parts[1]

	if tenant == "" || appCode == "" {
		return "", "", "", fmt.Errorf("tenant and app code cannot be empty")
	}

	apiURL = u.Scheme + "://" + u.Host

	return tenant, appCode, apiURL, nil
}
